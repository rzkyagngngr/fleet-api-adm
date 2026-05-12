package vesselschedule

import (
	"fmt"
	"net/http"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VesselScheduleHandler struct {
	service VesselScheduleService
}

func NewVesselScheduleHandler(service VesselScheduleService) *VesselScheduleHandler {
	return &VesselScheduleHandler{service: service}
}

// Search godoc
// @Summary Search vessel schedules
// @Description Retrieve vessel schedule data with server-side pagination, filtering, sorting, and download range support
// @Tags plan-vessel-schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body vesselschedule.SearchVesselScheduleRequest true "Vessel schedule search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /plan/vessel-schedule/search [post]
func (h *VesselScheduleHandler) Search(c *gin.Context) {
	var input SearchVesselScheduleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	branchCodeVal, exists := c.Get(middleware.BranchCodeKey)
	if !exists || branchCodeVal == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "branch code not found in token")
		return
	}
	terminalCodeVal, exists := c.Get(middleware.TerminalCodeKey)
	if !exists || terminalCodeVal == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "terminal code not found in token")
		return
	}

	branchCode, err := parseContextInt(branchCodeVal)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
		return
	}
	terminalCode, err := parseContextInt(terminalCodeVal)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
		return
	}

	if input.Filters == nil {
		input.Filters = map[string]string{}
	}
	input.Filters["branch_code"] = strconv.Itoa(branchCode)
	input.Filters["terminal_code"] = strconv.Itoa(terminalCode)

	rows, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search vessel schedules")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "vessel schedules retrieved successfully", rows, meta)
}

// GetByID godoc
// @Summary Get vessel schedule detail
// @Description Retrieve vessel schedule detail by id
// @Tags plan-vessel-schedule
// @Produce json
// @Security BearerAuth
// @Param id path int true "Vessel Schedule ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
func (h *VesselScheduleHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid vessel schedule id")
		return
	}

	row, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "vessel schedule not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "vessel schedule detail retrieved successfully", row)
}

// GetByScheduleCode godoc
// @Summary Get vessel schedule detail by schedule code
// @Description Retrieve vessel schedule detail by schedule_code
// @Tags plan-vessel-schedule
// @Produce json
// @Security BearerAuth
// @Param schedule_code query string true "Vessel Schedule Code"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /plan/vessel-schedule [get]
func (h *VesselScheduleHandler) GetByScheduleCode(c *gin.Context) {
	scheduleCode := c.Query("schedule_code")
	if scheduleCode == "" {
		scheduleCode = c.Param("schedule_code")
	}
	if scheduleCode == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "schedule_code is required")
		return
	}

	row, err := h.service.FindByScheduleCode(c.Request.Context(), scheduleCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "vessel schedule not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "vessel schedule detail retrieved successfully", row)
}

// Create godoc
// @Summary Create vessel schedule
// @Description Create a new vessel schedule record
// @Tags plan-vessel-schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body vesselschedule.VesselScheduleRequest true "Vessel schedule payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /plan/vessel-schedule [post]
func (h *VesselScheduleHandler) Create(c *gin.Context) {
	var req VesselScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	schedule, err := h.buildScheduleFromRequest(c, req)
	if err != nil {
		writeBuildError(c, err)
		return
	}

	if err := h.service.Create(c.Request.Context(), schedule); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create vessel schedule")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "vessel schedule created successfully", schedule)
}

// Update godoc
// @Summary Update vessel schedule
// @Description Update vessel schedule by schedule_code
// @Tags plan-vessel-schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param schedule_code query string true "Vessel Schedule Code"
// @Param payload body vesselschedule.VesselScheduleRequest true "Vessel schedule payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /plan/vessel-schedule [put]
func (h *VesselScheduleHandler) Update(c *gin.Context) {
	scheduleCode := c.Query("schedule_code")
	if scheduleCode == "" {
		scheduleCode = c.Param("schedule_code")
	}
	if scheduleCode == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "schedule_code is required")
		return
	}

	var req VesselScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	schedule, err := h.buildScheduleFromRequest(c, req)
	if err != nil {
		writeBuildError(c, err)
		return
	}

	if err := h.service.Update(c.Request.Context(), scheduleCode, schedule); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update vessel schedule")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "vessel schedule updated successfully", schedule)
}

// UpdateStatus godoc
// @Summary Update vessel schedule status
// @Description Update status by schedule_code from request body
// @Tags plan-vessel-schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body vesselschedule.UpdateVesselScheduleStatusRequest true "Vessel schedule status payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /plan/vessel-schedule/status [put]
func (h *VesselScheduleHandler) UpdateStatus(c *gin.Context) {
	var req UpdateVesselScheduleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	scheduleCode := strings.TrimSpace(req.ScheduleCode)
	if scheduleCode == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "schedule_code is required")
		return
	}
	if req.Status == nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "status is required")
		return
	}

	userName := middleware.GetUserEmail(c)
	if userName == "" {
		userName = "SYSTEM"
	}

	if err := h.service.UpdateStatus(c.Request.Context(), scheduleCode, *req.Status, userName); err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.ErrorResponse(c, http.StatusNotFound, "vessel schedule not found")
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update vessel schedule status")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "vessel schedule status updated successfully", gin.H{
		"schedule_code": scheduleCode,
		"status":        *req.Status,
	})
}

// Delete godoc
// @Summary Delete vessel schedule
// @Description Delete vessel schedule by id
// @Tags plan-vessel-schedule
// @Produce json
// @Security BearerAuth
// @Param id path int true "Vessel Schedule ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
func (h *VesselScheduleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid vessel schedule id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete vessel schedule")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "vessel schedule deleted successfully", nil)
}

func (h *VesselScheduleHandler) buildScheduleFromRequest(c *gin.Context, req VesselScheduleRequest) (*VesselSchedule, error) {
	userName := middleware.GetUserEmail(c)
	if userName == "" {
		userName = "SYSTEM"
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return nil, fmt.Errorf("user id not found in token")
	}
	branchCodeVal, exists := c.Get(middleware.BranchCodeKey)
	if !exists || branchCodeVal == nil {
		return nil, fmt.Errorf("branch code not found in token")
	}
	terminalCodeVal, exists := c.Get(middleware.TerminalCodeKey)
	if !exists || terminalCodeVal == nil {
		return nil, fmt.Errorf("terminal code not found in token")
	}

	branchCode, err := parseContextInt(branchCodeVal)
	if err != nil {
		return nil, fmt.Errorf("invalid branch code in token")
	}
	terminalCode, err := parseContextInt(terminalCodeVal)
	if err != nil {
		return nil, fmt.Errorf("invalid terminal code in token")
	}
	authLocation, err := h.service.GetAuthLocation(c.Request.Context(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve auth location")
	}
	branchName := authLocation.BranchName
	terminalName := authLocation.TerminalName

	return &VesselSchedule{
		BranchCode:             &branchCode,
		TerminalCode:           &terminalCode,
		BranchName:             &branchName,
		TerminalName:           &terminalName,
		VesselName:             req.VesselName,
		VesselCode:             req.VesselCode,
		VesselType:             req.VesselType,
		VesselHatchNumber:      req.VesselHatchNumber,
		VoyageNumber:           req.VoyageNumber,
		PKKNumber:              req.PKKNumber,
		PPKNumber:              req.PPKNumber,
		VoyageType:             req.VoyageType,
		GRT:                    req.GRT,
		LOA:                    req.LOA,
		AgencyName:             req.AgencyName,
		PortAgent:              req.PortAgent,
		EmergencyContact:       req.EmergencyContact,
		OriginPortCode:         req.OriginPortCode,
		OriginPortName:         req.OriginPortName,
		DestinationPortCode:    req.DestinationPortCode,
		DestinationPortName:    req.DestinationPortName,
		DischargePortCode:      req.DischargePortCode,
		DischargePortName:      req.DischargePortName,
		AssignedBerthName:      req.AssignedBerthName,
		DockID:                 req.DockID,
		DockCode:               req.DockCode,
		DockName:               req.DockName,
		BerthCode:              req.BerthCode,
		BerthName:              req.BerthName,
		BerthLatitude:          req.BerthLatitude,
		BerthLongitude:         req.BerthLongitude,
		CodeInaportnet:         req.CodeInaportnet,
		LocationNameInaportnet: req.LocationNameInaportnet,
		StartBerthPosition:     req.StartBerthPosition,
		EndBerthPosition:       req.EndBerthPosition,
		ETA:                    req.ETA,
		ETB:                    req.ETB,
		ETC:                    req.ETC,
		ETD:                    req.ETD,
		Status:                 req.Status,
		CreationBy:             &userName,
		LastUpdatedBy:          &userName,
	}, nil
}

func writeBuildError(c *gin.Context, err error) {
	switch err.Error() {
	case "branch code not found in token", "terminal code not found in token", "user id not found in token", "invalid branch code in token", "invalid terminal code in token":
		helper.ErrorResponse(c, http.StatusUnauthorized, err.Error())
	default:
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}

func parseContextInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case *int64:
		if v == nil {
			return 0, fmt.Errorf("nil int64 pointer")
		}
		return int(*v), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return i, nil
	case *string:
		if v == nil {
			return 0, fmt.Errorf("nil string pointer")
		}
		i, err := strconv.Atoi(*v)
		if err != nil {
			return 0, err
		}
		return i, nil
	default:
		return 0, fmt.Errorf("unsupported type %T", value)
	}
}

package vessel

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type VesselHandler struct {
	service VesselService
}

func NewVesselHandler(service VesselService) *VesselHandler {
	return &VesselHandler{service: service}
}

// Search godoc
// @Summary Search vessels
// @Description Search and filter vessels
// @Tags master-vessel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body vessel.SearchVesselsRequest true "Search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/vessel/search [post]
func (h *VesselHandler) Search(c *gin.Context) {
	var input SearchVesselsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rows, meta, err := h.service.SearchVessels(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
}

// Create godoc
// @Summary Create vessel
// @Description Create a new vessel
// @Tags master-vessel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body vessel.VesselRequest true "Vessel payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/vessel [post]
func (h *VesselHandler) Create(c *gin.Context) {
	var input VesselRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)

	if err := h.service.CreateVessel(c.Request.Context(), &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "vessel created successfully", nil)
}

// Update godoc
// @Summary Update vessel
// @Description Update vessel by id
// @Tags master-vessel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Vessel ID"
// @Param payload body vessel.VesselRequest true "Vessel payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/vessel [put]
// @Router /master/vessel/{id} [put]
func (h *VesselHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid vessel id")
		return
	}

	var input VesselRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)

	if err := h.service.UpdateVessel(c.Request.Context(), id, &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "vessel updated successfully", nil)
}

// GetByID godoc
// @Summary Get vessel by id
// @Description Get vessel by id
// @Tags master-vessel
// @Produce json
// @Security BearerAuth
// @Param id query int true "Vessel ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/vessel [get]
// @Router /master/vessel/{id} [get]
func (h *VesselHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid vessel id")
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

// Delete godoc
// @Summary Delete vessel
// @Description Delete vessel by id
// @Tags master-vessel
// @Produce json
// @Security BearerAuth
// @Param id query int true "Vessel ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/vessel [delete]
// @Router /master/vessel/{id} [delete]
func (h *VesselHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid vessel id")
		return
	}

	if err := h.service.DeleteVessel(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "vessel deleted successfully", nil)
}

// GetStats godoc
// @Summary Get vessel statistics
// @Description Get aggregated vessel statistics
// @Tags master-vessel
// @Produce json
// @Security BearerAuth
// @Param branch_code query int false "Branch Code"
// @Param terminal_code query int false "Terminal Code"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/vessel/stats [get]
func (h *VesselHandler) GetStats(c *gin.Context) {
	branchCode, _ := strconv.Atoi(c.DefaultQuery("branch_code", "0"))
	terminalCode, _ := strconv.Atoi(c.DefaultQuery("terminal_code", "0"))

	// Fallback to middleware context if query params are missing
	if branchCode == 0 {
		if val, ok := c.Get(middleware.BranchCodeKey); ok {
			branchCode, _ = strconv.Atoi(val.(string))
		}
	}
	if terminalCode == 0 {
		if val, ok := c.Get(middleware.TerminalCodeKey); ok {
			terminalCode, _ = strconv.Atoi(val.(string))
		}
	}

	res, err := h.service.GetVesselStats(c.Request.Context(), branchCode, terminalCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

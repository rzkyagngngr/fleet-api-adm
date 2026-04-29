package equipment

import (
	"fmt"
	"net/http"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EquipmentHandler struct {
	service EquipmentService
}

func NewEquipmentHandler(service EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{service: service}
}

// ListEquipmentGroupOptions godoc
// @Summary List equipment group options
// @Description Retrieve lightweight equipment group options from reference data
// @Tags master-equipments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body equipment.EquipmentGroupOptionRequest false "Equipment group option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/equipment/group-options/search [post]
func (h *EquipmentHandler) ListEquipmentGroupOptions(c *gin.Context) {
	var req EquipmentGroupOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	options, err := h.service.ListEquipmentGroupOptions(c.Request.Context(), req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve equipment group options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "equipment group options retrieved successfully", options)
}

// ListCustomerOptions godoc
// @Summary List customer options for equipment
// @Description Retrieve lightweight active customer options for equipment owner select inputs
// @Tags master-equipments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body equipment.CustomerOptionRequest false "Customer option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/equipment/customer-options/search [post]
func (h *EquipmentHandler) ListCustomerOptions(c *gin.Context) {
	var req CustomerOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
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

	options, err := h.service.ListCustomerOptions(
		c.Request.Context(),
		branchCode,
		terminalCode,
		req.Q,
		req.Limit,
	)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve customer options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "customer options retrieved successfully", options)
}

// SearchEquipments godoc
// @Summary Search equipments
// @Description Retrieve equipments with server-side pagination, filtering, sorting, and download range support
// @Tags master-equipments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body equipment.SearchEquipmentRequest true "Equipment search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/equipment/search [post]
func (h *EquipmentHandler) SearchEquipments(c *gin.Context) {
	var input SearchEquipmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
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

	if input.Filters == nil {
		input.Filters = map[string]string{}
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

	input.Filters["branch_code"] = strconv.Itoa(branchCode)
	input.Filters["terminal_code"] = strconv.Itoa(terminalCode)

	equipments, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search equipments")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "equipments retrieved successfully", equipments, meta)
}

// CreateEquipment godoc
// @Summary Create equipment
// @Description Create a new equipment record
// @Tags master-equipments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body equipment.EquipmentReq true "Equipment payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/equipment [post]
func (h *EquipmentHandler) CreateEquipment(c *gin.Context) {
	var req EquipmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userName := middleware.GetUserEmail(c)
	if userName == "" {
		userName = "SYSTEM"
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		helper.ErrorResponse(c, http.StatusUnauthorized, "user id not found in token")
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
	authLocation, err := h.service.GetAuthLocation(c.Request.Context(), userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to resolve auth location")
		return
	}
	branchName := authLocation.BranchName
	terminalName := authLocation.TerminalName

	equipment := Equipment{
		BranchCode:          &branchCode,
		BranchName:          &branchName,
		TerminalCode:        &terminalCode,
		TerminalName:        &terminalName,
		EquipmentName:       req.EquipmentName,
		EquipmentGroup:      req.EquipmentGroup,
		EquipmentType:       req.EquipmentType,
		Capacity:            req.Capacity,
		MinimalLoadCapacity: req.MinimalLoadCapacity,
		MaxLoadCapacity:     req.MaxLoadCapacity,
		OwnershipStatus:     req.OwnershipStatus,
		OwnerName:           req.OwnerName,
		OwnerCode:           req.OwnerCode,
		StartDate:           req.StartDate,
		EndDate:             req.EndDate,
		EquipmentCondition:  req.EquipmentCondition,
		Status:              req.Status,
		CreationBy:          &userName,
		LastUpdatedBy:       &userName,
	}

	if err := h.service.Create(c.Request.Context(), &equipment); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create equipment")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "equipment created successfully", equipment)
}

// UpdateEquipment godoc
// @Summary Update equipment
// @Description Update an existing equipment by id
// @Tags master-equipments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Equipment ID"
// @Param payload body equipment.EquipmentReq true "Equipment payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/equipment/{id} [put]
func (h *EquipmentHandler) UpdateEquipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid equipment id")
		return
	}

	var req EquipmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userName := middleware.GetUserEmail(c)
	if userName == "" {
		userName = "SYSTEM"
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		helper.ErrorResponse(c, http.StatusUnauthorized, "user id not found in token")
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
	authLocation, err := h.service.GetAuthLocation(c.Request.Context(), userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to resolve auth location")
		return
	}
	branchName := authLocation.BranchName
	terminalName := authLocation.TerminalName

	equipment := Equipment{
		BranchCode:          &branchCode,
		BranchName:          &branchName,
		TerminalCode:        &terminalCode,
		TerminalName:        &terminalName,
		EquipmentName:       req.EquipmentName,
		EquipmentGroup:      req.EquipmentGroup,
		EquipmentType:       req.EquipmentType,
		Capacity:            req.Capacity,
		MinimalLoadCapacity: req.MinimalLoadCapacity,
		MaxLoadCapacity:     req.MaxLoadCapacity,
		OwnershipStatus:     req.OwnershipStatus,
		OwnerName:           req.OwnerName,
		OwnerCode:           req.OwnerCode,
		StartDate:           req.StartDate,
		EndDate:             req.EndDate,
		EquipmentCondition:  req.EquipmentCondition,
		Status:              req.Status,
		LastUpdatedBy:       &userName,
	}

	existingEquipment, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to find equipment")
		return
	}
	equipment.EquipmentCode = existingEquipment.EquipmentCode

	if err := h.service.Update(c.Request.Context(), id, &equipment); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update equipment")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "equipment updated successfully", equipment)
}

// DeleteEquipment godoc
// @Summary Delete equipment
// @Description Delete equipment by id
// @Tags master-equipments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Equipment ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/equipment/{id} [delete]
func (h *EquipmentHandler) DeleteEquipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid equipment id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete equipment")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "equipment deleted successfully", nil)
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

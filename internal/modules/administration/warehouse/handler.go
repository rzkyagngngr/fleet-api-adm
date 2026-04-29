package warehouse

import (
	"fmt"
	"net/http"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	service WarehouseService
}

func NewWarehouseHandler(service WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

// SearchWarehouse godoc
// @Summary Search warehouse
// @Description Retrieve warehouse data with server-side pagination, filtering, sorting, and download range support
// @Tags master-warehouse
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body warehouse.SearchWarehouseRequest true "Warehouse search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/warehouse/search [post]
func (h *WarehouseHandler) SearchWarehouse(c *gin.Context) {
	var input SearchWarehouseRequest
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

	warehouses, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search warehouse")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "warehouse retrieved successfully", warehouses, meta)
}

// GetWarehouseDetail godoc
// @Summary Get warehouse detail
// @Description Retrieve warehouse detail by id
// @Tags master-warehouse
// @Produce json
// @Security BearerAuth
// @Param id path int true "Warehouse ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /master/warehouse/{id} [get]
func (h *WarehouseHandler) GetWarehouseDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid warehouse id")
		return
	}

	warehouse, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "warehouse not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "warehouse detail retrieved successfully", warehouse)
}

// CreateWarehouse godoc
// @Summary Create warehouse
// @Description Create a new warehouse record with details
// @Tags master-warehouse
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body warehouse.WarehouseReq true "Warehouse payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/warehouse [post]
func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var req WarehouseReq
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

	warehouse := Warehouse{
		BranchCode:        &branchCode,
		BranchName:        &branchName,
		TerminalCode:      &terminalCode,
		TerminalName:      &terminalName,
		WarehouseCode:     req.WarehouseCode,
		WarehouseName:     req.WarehouseName,
		WarehouseType:     req.WarehouseType,
		WarehouseCapacity: req.WarehouseCapacity,
		Status:            req.Status,
		CreationBy:        &userName,
		LastUpdatedBy:     &userName,
		Details:           mapWarehouseDetails(req.Details),
	}

	if err := h.service.Create(c.Request.Context(), &warehouse); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create warehouse")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "warehouse created successfully", warehouse)
}

// UpdateWarehouse godoc
// @Summary Update warehouse
// @Description Update an existing warehouse by id with details
// @Tags master-warehouse
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Warehouse ID"
// @Param payload body warehouse.WarehouseReq true "Warehouse payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/warehouse/{id} [put]
func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid warehouse id")
		return
	}

	var req WarehouseReq
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

	warehouse := Warehouse{
		BranchCode:        &branchCode,
		BranchName:        &branchName,
		TerminalCode:      &terminalCode,
		TerminalName:      &terminalName,
		WarehouseCode:     req.WarehouseCode,
		WarehouseName:     req.WarehouseName,
		WarehouseType:     req.WarehouseType,
		WarehouseCapacity: req.WarehouseCapacity,
		Status:            req.Status,
		LastUpdatedBy:     &userName,
		Details:           mapWarehouseDetails(req.Details),
	}

	if err := h.service.Update(c.Request.Context(), id, &warehouse); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update warehouse")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "warehouse updated successfully", warehouse)
}

// DeleteWarehouse godoc
// @Summary Delete warehouse
// @Description Delete warehouse by id
// @Tags master-warehouse
// @Produce json
// @Security BearerAuth
// @Param id path int true "Warehouse ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/warehouse/{id} [delete]
func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid warehouse id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete warehouse")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "warehouse deleted successfully", nil)
}

func mapWarehouseDetails(details []WarehouseDetailReq) []WarehouseDetail {
	if len(details) == 0 {
		return nil
	}

	result := make([]WarehouseDetail, 0, len(details))
	for _, detail := range details {
		result = append(result, WarehouseDetail{
			WarehouseCodeD:        detail.WarehouseCodeD,
			WarehouseNameD:        detail.WarehouseNameD,
			WerehouseDType:        detail.WerehouseDType,
			WarehouseCapacityDM3:  detail.WarehouseCapacityDM3,
			WarehouseCapacityDTon: detail.WarehouseCapacityDTon,
			XPosition:             detail.XPosition,
			YPosition:             detail.YPosition,
			WSize:                 detail.WSize,
			HSize:                 detail.HSize,
			Status:                detail.Status,
		})
	}

	return result
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

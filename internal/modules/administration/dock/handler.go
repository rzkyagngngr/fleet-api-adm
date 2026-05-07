package dock

import (
	"fmt"
	"net/http"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DockHandler struct {
	service DockService
}

func NewDockHandler(service DockService) *DockHandler {
	return &DockHandler{service: service}
}

// SearchDock godoc
// @Summary Search dock
// @Description Retrieve dock data with server-side pagination, filtering, sorting, and download range support
// @Tags master-dock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dock.SearchDockRequest true "Dock search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/dock/search [post]
func (h *DockHandler) SearchDock(c *gin.Context) {
	var input SearchDockRequest
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

	docks, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search dock")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "dock retrieved successfully", docks, meta)
}

// GetDockDetail godoc
// @Summary Get dock detail
// @Description Retrieve dock detail by id
// @Tags master-dock
// @Produce json
// @Security BearerAuth
// @Param id query int true "Dock ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /master/dock [get]
func (h *DockHandler) GetDockDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid dock id")
		return
	}

	dock, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "dock not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "dock detail retrieved successfully", dock)
}

// CreateDock godoc
// @Summary Create dock
// @Description Create a new dock record with details
// @Tags master-dock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dock.DockReq true "Dock payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/dock [post]
func (h *DockHandler) CreateDock(c *gin.Context) {
	var req DockReq
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

	dock := Dock{
		BranchCode:      &branchCode,
		BranchName:      &branchName,
		TerminalCode:    &terminalCode,
		TerminalName:    &terminalName,
		DockCode:        req.DockCode,
		DockName:        req.DockName,
		DockType:        req.DockType,
		DockLengthM:     req.DockLengthM,
		DockWidthM:      req.DockWidthM,
		DockCapacityTon: req.DockCapacityTon,
		CodeInaportnet:  req.CodeInaportnet,
		LocationNameIna: req.LocationNameIna,
		Status:          req.Status,
		CreationBy:      &userName,
		LastUpdatedBy:   &userName,
		Details:         mapDockDetails(req.Details),
	}

	if err := h.service.Create(c.Request.Context(), &dock); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create dock")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "dock created successfully", dock)
}

// UpdateDock godoc
// @Summary Update dock
// @Description Update an existing dock by id with details
// @Tags master-dock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Dock ID"
// @Param payload body dock.DockReq true "Dock payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/dock [put]
func (h *DockHandler) UpdateDock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid dock id")
		return
	}

	var req DockReq
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

	dock := Dock{
		BranchCode:      &branchCode,
		BranchName:      &branchName,
		TerminalCode:    &terminalCode,
		TerminalName:    &terminalName,
		DockCode:        req.DockCode,
		DockName:        req.DockName,
		DockType:        req.DockType,
		DockLengthM:     req.DockLengthM,
		DockWidthM:      req.DockWidthM,
		DockCapacityTon: req.DockCapacityTon,
		CodeInaportnet:  req.CodeInaportnet,
		LocationNameIna: req.LocationNameIna,
		Status:          req.Status,
		LastUpdatedBy:   &userName,
		Details:         mapDockDetails(req.Details),
	}

	if err := h.service.Update(c.Request.Context(), id, &dock); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update dock")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "dock updated successfully", dock)
}

// DeleteDock godoc
// @Summary Delete dock
// @Description Delete dock by id
// @Tags master-dock
// @Produce json
// @Security BearerAuth
// @Param id query int true "Dock ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/dock [delete]
func (h *DockHandler) DeleteDock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid dock id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete dock")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "dock deleted successfully", nil)
}

func mapDockDetails(details []DockDetailReq) []DockDetail {
	if len(details) == 0 {
		return nil
	}

	result := make([]DockDetail, 0, len(details))
	for _, detail := range details {
		result = append(result, DockDetail{
			BerthCode:  detail.BerthCode,
			BerthName:  detail.BerthName,
			MaxLoa:     detail.MaxLoa,
			XPosition:  detail.XPosition,
			YPosition:  detail.YPosition,
			WidthSize:  detail.WidthSize,
			HeightSize: detail.HeightSize,
			Status:     detail.Status,
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

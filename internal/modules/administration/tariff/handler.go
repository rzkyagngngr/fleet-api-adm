package tariff

import (
	"fmt"
	"net/http"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TariffHandler struct {
	service TariffServiceInterface
}

func NewTariffHandler(service TariffServiceInterface) *TariffHandler {
	return &TariffHandler{service: service}
}

// Search godoc
// @Summary Search tariff
// @Description Retrieve tariff data with server-side pagination, filtering, sorting, and download range support
// @Tags master-tariff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body tariff.SearchTariffRequest true "Tariff search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/tariff/search [post]
func (h *TariffHandler) Search(c *gin.Context) {
	var input SearchTariffRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	rows, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search tariff")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "tariff retrieved successfully", rows, meta)
}

// SearchStatusZero godoc
// @Summary Search tariff status zero
// @Description Retrieve tariff data with status 0 using server-side pagination, filtering, sorting, and download range support
// @Tags master-tariff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body tariff.SearchTariffRequest true "Tariff search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/tariff/status-zero/search [post]
func (h *TariffHandler) SearchStatusZero(c *gin.Context) {
	var input SearchTariffRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	rows, meta, err := h.service.SearchStatusZero(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search tariff status zero")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "tariff status zero retrieved successfully", rows, meta)
}

// GetByID godoc
// @Summary Get tariff detail
// @Description Retrieve tariff detail by id
// @Tags master-tariff
// @Produce json
// @Security BearerAuth
// @Param id query int true "Tariff ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /master/tariff [get]
func (h *TariffHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid tariff id")
		return
	}

	row, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "tariff not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "tariff detail retrieved successfully", row)
}

// Create godoc
// @Summary Create tariff
// @Description Create a new tariff record with details
// @Tags master-tariff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body tariff.TariffReq true "Tariff payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/tariff [post]
func (h *TariffHandler) Create(c *gin.Context) {
	var req TariffReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
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

	row := Tariff{
		BranchCode:      &branchCode,
		BranchName:      &branchName,
		TerminalCode:    &terminalCode,
		TerminalName:    &terminalName,
		NameTariff:      req.NameTariff,
		Description:     req.Description,
		Status:          req.Status,
		AgreementNumber: req.AgreementNumber,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		CreationBy:      &userName,
		LastUpdatedBy:   &userName,
		Details:         mapTariffDetails(req.Details, branchCode, branchName, terminalCode, terminalName),
	}

	if err := h.service.Create(c.Request.Context(), &row); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create tariff")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "tariff created successfully", row)
}

// Update godoc
// @Summary Update tariff
// @Description Update an existing tariff by id with details
// @Tags master-tariff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Tariff ID"
// @Param payload body tariff.TariffReq true "Tariff payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/tariff [put]
func (h *TariffHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid tariff id")
		return
	}

	var req TariffReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
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

	row := Tariff{
		BranchCode:      &branchCode,
		BranchName:      &branchName,
		TerminalCode:    &terminalCode,
		TerminalName:    &terminalName,
		NameTariff:      req.NameTariff,
		Description:     req.Description,
		Status:          req.Status,
		AgreementNumber: req.AgreementNumber,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		LastUpdatedBy:   &userName,
		Details:         mapTariffDetails(req.Details, branchCode, branchName, terminalCode, terminalName),
	}

	if err := h.service.Update(c.Request.Context(), id, &row); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update tariff")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "tariff updated successfully", row)
}

// UpdateStatus godoc
// @Summary Update tariff status
// @Description Update tariff status only by id
// @Tags master-tariff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Tariff ID"
// @Param payload body tariff.UpdateTariffStatusRequest true "Tariff status payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/tariff/{id}/status [put]
func (h *TariffHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid tariff id")
		return
	}

	var req UpdateTariffStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	userName := middleware.GetUserEmail(c)
	if userName == "" {
		userName = "SYSTEM"
	}

	if err := h.service.UpdateStatus(c.Request.Context(), id, req.Status, &userName); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update tariff status")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "tariff status updated successfully", gin.H{
		"id":     id,
		"status": req.Status,
	})
}

// Delete godoc
// @Summary Delete tariff
// @Description Delete tariff by id
// @Tags master-tariff
// @Produce json
// @Security BearerAuth
// @Param id query int true "Tariff ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/tariff [delete]
func (h *TariffHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid tariff id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete tariff")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "tariff deleted successfully", nil)
}

func mapTariffDetails(details []TariffServiceReq, branchCode int, branchName string, terminalCode int, terminalName string) []TariffService {
	if len(details) == 0 {
		return nil
	}

	result := make([]TariffService, 0, len(details))
	for _, detail := range details {
		result = append(result, TariffService{
			BranchCode:     branchCode,
			BranchName:     &branchName,
			TerminalCode:   terminalCode,
			TerminalName:   &terminalName,
			SequenceNo:     detail.SequenceNo,
			ServiceType:    detail.ServiceType,
			ServiceName:    detail.ServiceName,
			CustomerName:   detail.CustomerName,
			CustomerCode:   detail.CustomerCode,
			CargoCode:      detail.CargoCode,
			CargoName:      detail.CargoName,
			CargoPackaging: detail.CargoPackaging,
			CargoUnit:      detail.CargoUnit,
			EquipmentCode:  detail.EquipmentCode,
			EquipmentName:  detail.EquipmentName,
			EquipmentGroup: detail.EquipmentGroup,
			EquipmentUnit:  detail.EquipmentUnit,
			BaseTariff:     detail.BaseTariff,
			CurrencyCode:   detail.CurrencyCode,
			Discount:       detail.Discount,
			Attrib1:        detail.Attrib1,
			Attrib2:        detail.Attrib2,
			Attrib3:        detail.Attrib3,
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

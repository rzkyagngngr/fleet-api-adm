package handler

import (
	"gin-boilerplate/internal/middleware"
	"gin-boilerplate/internal/model/dto"
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/service"
	"gin-boilerplate/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DermagaHandler interface {
	Create(c *gin.Context)
	FindAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	FindByID(c *gin.Context)
}

type dermagaHandler struct {
	dermagaService service.DermagaService
}

func NewDermagaHandler(dermagaService service.DermagaService) DermagaHandler {
	return &dermagaHandler{dermagaService: dermagaService}
}

func (h *dermagaHandler) Create(c *gin.Context) {
	var input dto.DermagaRequest

	// Extract values from context (populated by AuthMiddleware)
	employeeID, _ := c.Get(middleware.EmployeeIDKey)
	branchCode, _ := c.Get(middleware.BranchCodeKey)

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Note: You might need to adjust this depending on how Dermaga entity maps to the new structure
	dermaga := &entity.Dermaga{
		KdCabang:    uint(*branchCode.(*int64)), // Casting based on the new uint64/int64 schema
		NmCabang:    input.NmDermaga,           // Fallback or adjust as needed
		KdTerminal:  0,                         // Adjust as needed
		NmTerminal:  "",                        // Adjust as needed
		NmDermaga:   input.NmDermaga,
		KdDermaga:   input.KdDermaga,
		PosisiAwal:  input.PosisiAwal,
		PosisiAkhir: input.PosisiAkhir,
		Keterangan:  input.Keterangan,
		Status:      input.Status,
		CreatedBy:   employeeID.(string),
		UpdatedBy:   employeeID.(string),
	}

	if err := h.dermagaService.Create(c.Request.Context(), dermaga); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "dermaga created successfully", dermaga)
}

// Create godoc
// @Summary Create a new dermaga
// @Tags dermaga
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.DermagaRequest true "Dermaga payload"
// @Success 201 {object} utils.Response
// @Router /api/v1/dermaga [post]
func (h *dermagaHandler) _Create(c *gin.Context) {} // Placeholder for documentation if needed

// FindAll godoc
// @Summary List all dermagas
// @Tags dermaga
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} utils.PaginatedResponse
// @Router /api/v1/dermaga [get]
func (h *dermagaHandler) FindAll(c *gin.Context) {
	// Extract values from context (populated by AuthMiddleware)
	branchCode, _ := c.Get(middleware.BranchCodeKey)

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	dermagas, total, err := h.dermagaService.FindAll(c.Request.Context(), uint(*branchCode.(*int64)), 0, page, size)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PaginatedSuccessResponse(c, http.StatusOK, "dermaga retrieved successfully", dermagas, total, page, size)
}

// Update godoc
// @Summary Update a dermaga
// @Tags dermaga
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Dermaga ID"
// @Param body body dto.DermagaRequest true "Dermaga payload"
// @Success 200 {object} utils.Response
// @Router /api/v1/dermaga/{id} [put]
func (h *dermagaHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid dermaga ID")
		return
	}

	var input dto.DermagaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Extract values from context (populated by AuthMiddleware)
	employeeID, _ := c.Get(middleware.EmployeeIDKey)
	branchCode, _ := c.Get(middleware.BranchCodeKey)

	dermaga := &entity.Dermaga{
		KdCabang:    uint(*branchCode.(*int64)),
		NmDermaga:   input.NmDermaga,
		KdDermaga:   input.KdDermaga,
		PosisiAwal:  input.PosisiAwal,
		PosisiAkhir: input.PosisiAkhir,
		Keterangan:  input.Keterangan,
		Status:      input.Status,
		UpdatedBy:   employeeID.(string),
	}

	if err := h.dermagaService.Update(c.Request.Context(), uint(id), uint(*branchCode.(*int64)), 0, dermaga); err != nil {
		if err == service.ErrUnauthorized {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		if err == service.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "dermaga updated successfully", nil)
}

// Delete godoc
// @Summary Delete a dermaga
// @Tags dermaga
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Dermaga ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/dermaga/{id} [delete]
func (h *dermagaHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid dermaga ID")
		return
	}

	// Extract values from context (populated by AuthMiddleware)
	branchCode, _ := c.Get(middleware.BranchCodeKey)

	if err := h.dermagaService.Delete(c.Request.Context(), uint(id), uint(*branchCode.(*int64)), 0); err != nil {
		if err == service.ErrUnauthorized {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		if err == service.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "dermaga deleted successfully", nil)
}

// FindByID godoc
// @Summary Get a dermaga by ID
// @Tags dermaga
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Dermaga ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/dermaga/{id} [get]
func (h *dermagaHandler) FindByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid dermaga ID")
		return
	}

	// Correctly implement based on the new service signature if needed
	dermaga, err := h.dermagaService.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "dermaga retrieved successfully", dermaga)
}

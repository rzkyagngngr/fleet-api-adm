package cargo

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type CargoHandler struct{ service CargoService }

func NewCargoHandler(service CargoService) *CargoHandler {
	return &CargoHandler{service: service}
}

// Search godoc
// @Summary Search cargo
// @Tags barang
// @Accept json
// @Produce json
// @Param payload body SearchCargoRequest true "Search payload"
// @Success 200 {object} helper.MetaResponse
// @Router /master/barang/search [post]
func (h *CargoHandler) Search(c *gin.Context) {
	var input SearchCargoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rows, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
}

// Create godoc
func (h *CargoHandler) Create(c *gin.Context) {
	var input CargoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)

	if err := h.service.CreateCargo(c.Request.Context(), &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "cargo created successfully", nil)
}

// Update godoc
func (h *CargoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid cargo ID")
		return
	}

	var input CargoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)

	if err := h.service.UpdateCargo(c.Request.Context(), id, &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "cargo updated successfully", nil)
}

// Delete godoc
func (h *CargoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid cargo ID")
		return
	}

	if err := h.service.DeleteCargo(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "cargo deleted successfully", nil)
}

// GetByID godoc
func (h *CargoHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid cargo ID")
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

// GetStats godoc
func (h *CargoHandler) GetStats(c *gin.Context) {
	res, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

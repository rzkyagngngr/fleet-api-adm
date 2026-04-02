package dermaga

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type DermagaHandler interface {
	Create(c *gin.Context)
	FindAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	FindByID(c *gin.Context)
}

type dermagaHandler struct{ dermagaService DermagaService }

func NewDermagaHandler(svc DermagaService) DermagaHandler {
	return &dermagaHandler{dermagaService: svc}
}

// Create godoc
// @Summary Create dermaga
// @Description Create a new dermaga record
// @Tags dermaga
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dermaga.DermagaRequest true "Dermaga payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /dermaga [post]
func (h *dermagaHandler) Create(c *gin.Context) {
	var input DermagaRequest
	employeeID, _ := c.Get(middleware.EmployeeIDKey)
	branchCode, _ := c.Get(middleware.BranchCodeKey)

	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	d := &Dermaga{KdCabang: uint(*branchCode.(*int64)), NmDermaga: input.NmDermaga, KdDermaga: input.KdDermaga, PosisiAwal: input.PosisiAwal, PosisiAkhir: input.PosisiAkhir, Keterangan: input.Keterangan, Status: input.Status, CreatedBy: employeeID.(string), UpdatedBy: employeeID.(string)}

	if err := h.dermagaService.Create(c.Request.Context(), d); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "dermaga created successfully", d)
}

// FindAll godoc
// @Summary Get all dermaga
// @Description Retrieve paginated dermaga list
// @Tags dermaga
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} helper.PaginatedResponse
// @Failure 500 {object} helper.Response
// @Router /dermaga [get]
func (h *dermagaHandler) FindAll(c *gin.Context) {
	branchCode, _ := c.Get(middleware.BranchCodeKey)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	rows, total, err := h.dermagaService.FindAll(c.Request.Context(), uint(*branchCode.(*int64)), 0, page, size)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.PaginatedSuccessResponse(c, http.StatusOK, "dermaga retrieved successfully", rows, total, page, size)
}

// Update godoc
// @Summary Update dermaga
// @Description Update dermaga by id
// @Tags dermaga
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Dermaga ID"
// @Param payload body dermaga.DermagaRequest true "Dermaga payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /dermaga/{id} [put]
func (h *dermagaHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid dermaga ID")
		return
	}

	var input DermagaRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)
	branchCode, _ := c.Get(middleware.BranchCodeKey)

	d := &Dermaga{KdCabang: uint(*branchCode.(*int64)), NmDermaga: input.NmDermaga, KdDermaga: input.KdDermaga, PosisiAwal: input.PosisiAwal, PosisiAkhir: input.PosisiAkhir, Keterangan: input.Keterangan, Status: input.Status, UpdatedBy: employeeID.(string)}

	if err := h.dermagaService.Update(c.Request.Context(), uint(id), uint(*branchCode.(*int64)), 0, d); err != nil {
		if err == ErrUnauthorized {
			helper.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		if err == ErrNotFound {
			helper.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "dermaga updated successfully", nil)
}

// Delete godoc
// @Summary Delete dermaga
// @Description Delete dermaga by id
// @Tags dermaga
// @Produce json
// @Security BearerAuth
// @Param id path int true "Dermaga ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /dermaga/{id} [delete]
func (h *dermagaHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid dermaga ID")
		return
	}
	branchCode, _ := c.Get(middleware.BranchCodeKey)

	if err := h.dermagaService.Delete(c.Request.Context(), uint(id), uint(*branchCode.(*int64)), 0); err != nil {
		if err == ErrUnauthorized {
			helper.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		if err == ErrNotFound {
			helper.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "dermaga deleted successfully", nil)
}

// FindByID godoc
// @Summary Get dermaga by id
// @Description Retrieve dermaga detail by id
// @Tags dermaga
// @Produce json
// @Security BearerAuth
// @Param id path int true "Dermaga ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /dermaga/{id} [get]
func (h *dermagaHandler) FindByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid dermaga ID")
		return
	}
	d, err := h.dermagaService.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "dermaga retrieved successfully", d)
}

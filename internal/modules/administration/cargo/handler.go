package cargo

import (
	"fmt"
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
// @Description Search and filter cargo masters
// @Tags master-barang
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body cargo.SearchCargoRequest true "Search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang/search [post]
func (h *CargoHandler) Search(c *gin.Context) {
	var input SearchCargoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Inject context-based filters
	branchCode, _ := c.Get(middleware.BranchCodeKey)
	terminalCode, _ := c.Get(middleware.TerminalCodeKey)

	if input.Filters == nil {
		input.Filters = make(map[string]string)
	}

	if branchCode != nil {
		input.Filters["branch_code"] = branchCode.(string)
	}
	if terminalCode != nil {
		input.Filters["terminal_code"] = terminalCode.(string)
	}

	rows, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
}

// Create godoc
// @Summary Create cargo
// @Description Create a new cargo master
// @Tags master-barang
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body cargo.CargoRequest true "Cargo payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang [post]
func (h *CargoHandler) Create(c *gin.Context) {
	var input CargoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)
	branchCode, _ := c.Get(middleware.BranchCodeKey)
	terminalCode, _ := c.Get(middleware.TerminalCodeKey)

	if branchCode != nil {
		code, _ := strconv.Atoi(branchCode.(string))
		input.BranchCode = code
	}
	if terminalCode != nil {
		code, _ := strconv.Atoi(terminalCode.(string))
		input.TerminalCode = code
	}

	if err := h.service.CreateCargo(c.Request.Context(), &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	msg := fmt.Sprintf("cargo created successfully (Branch: %v, Terminal: %v)", input.BranchCode, input.TerminalCode)
	helper.SuccessResponse(c, http.StatusCreated, msg, nil)
}

// Update godoc
// @Summary Update cargo
// @Description Update cargo by id
// @Tags master-barang
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cargo ID"
// @Param payload body cargo.CargoRequest true "Cargo payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang/{id} [put]
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
	branchCode, _ := c.Get(middleware.BranchCodeKey)
	terminalCode, _ := c.Get(middleware.TerminalCodeKey)

	if branchCode != nil {
		code, _ := strconv.Atoi(branchCode.(string))
		input.BranchCode = code
	}
	if terminalCode != nil {
		code, _ := strconv.Atoi(terminalCode.(string))
		input.TerminalCode = code
	}

	if err := h.service.UpdateCargo(c.Request.Context(), id, &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	msg := fmt.Sprintf("cargo updated successfully (Branch: %v, Terminal: %v)", input.BranchCode, input.TerminalCode)
	helper.SuccessResponse(c, http.StatusOK, msg, nil)
}

// Delete godoc
// @Summary Delete cargo
// @Description Delete cargo by id
// @Tags master-barang
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cargo ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang/{id} [delete]
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
// @Summary Get cargo by id
// @Description Get cargo detail by id
// @Tags master-barang
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cargo ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang/{id} [get]
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
// @Summary Get cargo statistics
// @Description Get aggregated cargo statistics for active branch/terminal
// @Tags master-barang
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang/stats [get]
func (h *CargoHandler) GetStats(c *gin.Context) {
	branchCodeVal, _ := c.Get(middleware.BranchCodeKey)
	terminalCodeVal, _ := c.Get(middleware.TerminalCodeKey)

	branchCode := 0
	terminalCode := 0
	if branchCodeVal != nil {
		branchCode, _ = strconv.Atoi(branchCodeVal.(string))
	}
	if terminalCodeVal != nil {
		terminalCode, _ = strconv.Atoi(terminalCodeVal.(string))
	}

	res, err := h.service.GetStats(c.Request.Context(), branchCode, terminalCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

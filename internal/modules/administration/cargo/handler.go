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
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "cargo service is not initialized")
		return
	}

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
		branchCodeText, err := parseContextString(branchCode)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		input.Filters["branch_code"] = branchCodeText
	}
	if terminalCode != nil {
		terminalCodeText, err := parseContextString(terminalCode)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		input.Filters["terminal_code"] = terminalCodeText
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
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "cargo service is not initialized")
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
		code, _ := parseContextInt(branchCode)
		input.BranchCode = code
	}
	if terminalCode != nil {
		code, _ := parseContextInt(terminalCode)
		input.TerminalCode = code
	}

	adminEmp, _ := parseContextString(employeeID)
	if err := h.service.CreateCargo(c.Request.Context(), &input, adminEmp); err != nil {
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
// @Param id query int true "Cargo ID"
// @Param payload body cargo.CargoRequest true "Cargo payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang [put]
func (h *CargoHandler) Update(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "cargo service is not initialized")
		return
	}

	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
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
		code, _ := parseContextInt(branchCode)
		input.BranchCode = code
	}
	if terminalCode != nil {
		code, _ := parseContextInt(terminalCode)
		input.TerminalCode = code
	}

	adminEmp, _ := parseContextString(employeeID)
	if err := h.service.UpdateCargo(c.Request.Context(), id, &input, adminEmp); err != nil {
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
// @Param id query int true "Cargo ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang [delete]
func (h *CargoHandler) Delete(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "cargo service is not initialized")
		return
	}

	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
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
// @Param id query int true "Cargo ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/barang [get]
func (h *CargoHandler) GetByID(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "cargo service is not initialized")
		return
	}

	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
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
// @Router /master/barang/stats [get]
func (h *CargoHandler) GetStats(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "cargo service is not initialized")
		return
	}

	// Priority 1: Query parameters
	branchCode, _ := strconv.Atoi(c.DefaultQuery("branch_code", "0"))
	terminalCode, _ := strconv.Atoi(c.DefaultQuery("terminal_code", "0"))

	// Priority 2: Fallback to middleware context
	if branchCode == 0 {
		if val, ok := c.Get(middleware.BranchCodeKey); ok {
			branchCode, _ = parseContextInt(val)
		}
	}
	if terminalCode == 0 {
		if val, ok := c.Get(middleware.TerminalCodeKey); ok {
			terminalCode, _ = parseContextInt(val)
		}
	}

	res, err := h.service.GetStats(c.Request.Context(), branchCode, terminalCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

func parseContextString(value interface{}) (string, error) {
	switch v := value.(type) {
	case nil:
		return "", nil
	case string:
		return v, nil
	case *string:
		if v == nil {
			return "", nil
		}
		return *v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case *int64:
		if v == nil {
			return "", nil
		}
		return strconv.FormatInt(*v, 10), nil
	default:
		return "", fmt.Errorf("unsupported type %T", value)
	}
}

func parseContextInt(value interface{}) (int, error) {
	text, err := parseContextString(value)
	if err != nil || text == "" {
		return 0, err
	}
	return strconv.Atoi(text)
}

package lookup

import (
	"fmt"
	"net/http"
	"strconv"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"omniport-api/internal/modules/administration/equipment"

	"github.com/gin-gonic/gin"
)

type LookupHandler struct {
	service LookupService
}

func NewLookupHandler(service LookupService) *LookupHandler {
	return &LookupHandler{service: service}
}

// ListEquipmentGroupOptions godoc
// @Summary List equipment group options
// @Description Retrieve lightweight equipment group options from reference data
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body equipment.EquipmentGroupOptionRequest false "Equipment group option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/equipment-groups/search [post]
func (h *LookupHandler) ListEquipmentGroupOptions(c *gin.Context) {
	var req equipment.EquipmentGroupOptionRequest
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
// @Summary List customer options
// @Description Retrieve lightweight active customer options for lookup needs
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body equipment.CustomerOptionRequest false "Customer option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/customers/search [post]
func (h *LookupHandler) ListCustomerOptions(c *gin.Context) {
	var req equipment.CustomerOptionRequest
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

	options, err := h.service.ListCustomerOptions(c.Request.Context(), branchCode, terminalCode, req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve customer options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "customer options retrieved successfully", options)
}

// ListEquipmentOptions godoc
// @Summary List equipment options
// @Description Retrieve equipment code and equipment name from equipment master filtered by branch and terminal
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body lookup.SearchOptionRequest false "Equipment option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/equipments/search [post]
func (h *LookupHandler) ListEquipmentOptions(c *gin.Context) {
	var req SearchOptionRequest
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

	options, err := h.service.ListEquipmentOptions(c.Request.Context(), branchCode, terminalCode, req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve equipment options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "equipment options retrieved successfully", options)
}

// ListCargoPackageOptions godoc
// @Summary List cargo package options
// @Description Retrieve cargo package options from reference data
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body lookup.SearchOptionRequest false "Cargo package option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/cargo-packages/search [post]
func (h *LookupHandler) ListCargoPackageOptions(c *gin.Context) {
	var req SearchOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	options, err := h.service.ListCargoPackageOptions(c.Request.Context(), req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve cargo package options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "cargo package options retrieved successfully", options)
}

// ListCargoUnitOptions godoc
// @Summary List cargo unit options
// @Description Retrieve cargo unit options from reference data
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body lookup.SearchOptionRequest false "Cargo unit option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/cargo-units/search [post]
func (h *LookupHandler) ListCargoUnitOptions(c *gin.Context) {
	var req SearchOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	options, err := h.service.ListCargoUnitOptions(c.Request.Context(), req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve cargo unit options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "cargo unit options retrieved successfully", options)
}

// ListBillingServiceOptions godoc
// @Summary List billing service options
// @Description Retrieve billing service options from reference data filtered by branch and terminal
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body lookup.SearchOptionRequest false "Billing service option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/billing-services/search [post]
func (h *LookupHandler) ListBillingServiceOptions(c *gin.Context) {
	var req SearchOptionRequest
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

	options, err := h.service.ListBillingServiceOptions(c.Request.Context(), branchCode, terminalCode, req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve billing service options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "billing service options retrieved successfully", options)
}

// ListCargoOptions godoc
// @Summary List cargo options
// @Description Retrieve cargo code and cargo name from cargo master
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body lookup.SearchOptionRequest false "Cargo option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/cargos/search [post]
func (h *LookupHandler) ListCargoOptions(c *gin.Context) {
	var req SearchOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	options, err := h.service.ListCargoOptions(c.Request.Context(), req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve cargo options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "cargo options retrieved successfully", options)
}

// ListDockOptions godoc
// @Summary List dock options
// @Description Retrieve active dock header and detail options filtered by branch and terminal
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body lookup.SearchOptionRequest false "Dock option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/docks/search [post]
func (h *LookupHandler) ListDockOptions(c *gin.Context) {
	var req SearchOptionRequest
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

	options, err := h.service.ListDockOptions(c.Request.Context(), branchCode, terminalCode, req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve dock options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "dock options retrieved successfully", options)
}

// ListVesselOptions godoc
// @Summary List vessel options
// @Description Retrieve active vessel options from posm_vessel filtered by branch and terminal
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body lookup.SearchOptionRequest false "Vessel option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/vessels/search [post]
func (h *LookupHandler) ListVesselOptions(c *gin.Context) {
	var req SearchOptionRequest
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

	options, err := h.service.ListVesselOptions(c.Request.Context(), branchCode, terminalCode, req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve vessel options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "vessel options retrieved successfully", options)
}

// ListPortOptions godoc
// @Summary List port options
// @Description Retrieve active port options from adm.posm_port
// @Tags master-lookup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body lookup.SearchOptionRequest false "Port option payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/lookup/ports/search [post]
func (h *LookupHandler) ListPortOptions(c *gin.Context) {
	var req SearchOptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	options, err := h.service.ListPortOptions(c.Request.Context(), req.Q, req.Limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve port options")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "port options retrieved successfully", options)
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

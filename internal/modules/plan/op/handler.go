package op

import (
	"net/http"
	"strconv"
	"strings"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type OpsPlanHandler struct {
	service OpsPlanService
}

func NewOpsPlanHandler(service OpsPlanService) *OpsPlanHandler {
	return &OpsPlanHandler{service: service}
}

// ReadyOpsPlan godoc
// @Summary      Ready Operation Plan Requests
// @Description  Paginated list of approved requests that are ready to become operation plans
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.SearchReadyOpsPlanInput true "Search payload"
// @Success      200 {object} helper.MetaResponse
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/readyOp [post]
func (h *OpsPlanHandler) ReadyOpsPlan(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input SearchReadyOpsPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Filters == nil {
		input.Filters = make(map[string]string)
	}
	if bc, ok := c.Get(middleware.BranchCodeKey); ok {
		branchCode, err := parseContextString(bc)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		input.Filters["branch_code"] = branchCode
	}
	if tc, ok := c.Get(middleware.TerminalCodeKey); ok {
		terminalCode, err := parseContextString(tc)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		input.Filters["terminal_code"] = terminalCode
	}

	rows, meta, err := h.service.SearchReady(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
}

// Create godoc
// @Summary      Create Loading Unloading Plan
// @Description  Insert loading/unloading plan header and detail rows, then mark matching approved post_request as planned
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.CreateLoadingUnloadingPlanInput true "Create payload"
// @Success      201 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op [post]
func (h *OpsPlanHandler) Create(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input CreateLoadingUnloadingPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	branchCode, terminalCode := 0, 0
	if v, ok := c.Get(middleware.BranchCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		branchCode = value
	}
	if v, ok := c.Get(middleware.TerminalCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		terminalCode = value
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		helper.ErrorResponse(c, http.StatusUnauthorized, "user id not found in token")
		return
	}
	authLocation, err := h.service.GetAuthLocation(c.Request.Context(), userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to resolve auth location")
		return
	}
	branchName := authLocation.BranchName
	terminalName := authLocation.TerminalName

	createdBy := middleware.GetUserEmail(c)
	if v, ok := c.Get(middleware.EmployeeIDKey); ok {
		if employeeID, ok := v.(string); ok && employeeID != "" {
			createdBy = employeeID
		}
	}

	res, err := h.service.Create(
		c.Request.Context(),
		&input,
		branchCode,
		terminalCode,
		branchName,
		terminalName,
		createdBy,
	)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "loading unloading plan created successfully", res)
}

// CreateDetermination godoc
// @Summary      Create Loading Unloading Determination
// @Description  Insert loading/unloading determination header, detail rows, and equipment rows
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.CreateLoadingUnloadingDeterminationInput true "Create determination payload"
// @Success      201 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/createDetermination [post]
func (h *OpsPlanHandler) CreateDetermination(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input CreateLoadingUnloadingDeterminationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	branchCode, terminalCode := 0, 0
	if v, ok := c.Get(middleware.BranchCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		branchCode = value
	}
	if v, ok := c.Get(middleware.TerminalCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		terminalCode = value
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		helper.ErrorResponse(c, http.StatusUnauthorized, "user id not found in token")
		return
	}
	authLocation, err := h.service.GetAuthLocation(c.Request.Context(), userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to resolve auth location")
		return
	}
	branchName := authLocation.BranchName
	terminalName := authLocation.TerminalName

	createdBy := middleware.GetUserEmail(c)
	if v, ok := c.Get(middleware.EmployeeIDKey); ok {
		if employeeID, ok := v.(string); ok && employeeID != "" {
			createdBy = employeeID
		}
	}

	res, err := h.service.CreateDetermination(
		c.Request.Context(),
		&input,
		branchCode,
		terminalCode,
		branchName,
		terminalName,
		createdBy,
	)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "loading unloading determination created successfully", res)
}

// Update godoc
// @Summary      Update Loading Unloading Plan
// @Description  Update selected header fields and optionally replace detail rows by plan_code
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.UpdateLoadingUnloadingPlanInput true "Update payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/update [post]
func (h *OpsPlanHandler) Update(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input UpdateLoadingUnloadingPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	branchCode, terminalCode := 0, 0
	if v, ok := c.Get(middleware.BranchCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		branchCode = value
	}
	if v, ok := c.Get(middleware.TerminalCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		terminalCode = value
	}

	updatedBy := middleware.GetUserEmail(c)
	if v, ok := c.Get(middleware.EmployeeIDKey); ok {
		if employeeID, ok := v.(string); ok && employeeID != "" {
			updatedBy = employeeID
		}
	}

	res, err := h.service.Update(c.Request.Context(), &input, branchCode, terminalCode, updatedBy)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "loading unloading plan updated successfully", res)
}

// UpdateDeterminedPlan godoc
// @Summary      Update Determined Loading Unloading Plan
// @Description  Update status 1/2 plan and rebuild related determination details without regenerating determination_code or work_order_code
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.UpdateLoadingUnloadingPlanInput true "Update payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/updateDetermination [post]
func (h *OpsPlanHandler) UpdateDeterminedPlan(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input UpdateLoadingUnloadingPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	branchCode, terminalCode := 0, 0
	if v, ok := c.Get(middleware.BranchCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		branchCode = value
	}
	if v, ok := c.Get(middleware.TerminalCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		terminalCode = value
	}

	updatedBy := middleware.GetUserEmail(c)
	if v, ok := c.Get(middleware.EmployeeIDKey); ok {
		if employeeID, ok := v.(string); ok && employeeID != "" {
			updatedBy = employeeID
		}
	}

	res, err := h.service.UpdateDeterminedPlan(c.Request.Context(), &input, branchCode, terminalCode, updatedBy)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "loading unloading determined plan updated successfully", res)
}

// GetDataRequest godoc
// @Summary      Get Operation Plan Request Data
// @Description  Grouped cargo details from approved requests by pkk_number and activity_code
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.GetDataRequestInput true "Request payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/getDataRequest [post]
func (h *OpsPlanHandler) GetDataRequest(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input GetDataRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	requestNumber := input.RequestNumber()
	if requestNumber == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "pkk_number is required")
		return
	}
	if input.ActivityCode == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "activity_code is required")
		return
	}
	if input.ActivityCode != "BONGKAR" && input.ActivityCode != "MUAT" {
		helper.ErrorResponse(c, http.StatusBadRequest, "activity_code must be BONGKAR or MUAT")
		return
	}

	rows, err := h.service.GetDataRequest(c.Request.Context(), requestNumber, input.ActivityCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", rows)
}

// GetDataOp godoc
// @Summary      Get Operation Plan Data
// @Description  Loading/unloading plan header list with first berth name from detail rows
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.GetDataOpInput false "Optional filter payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/getDataOp [post]
func (h *OpsPlanHandler) GetDataOp(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input GetDataOpInput
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&input); err != nil {
			helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	branchCode, terminalCode := 0, 0
	if v, ok := c.Get(middleware.BranchCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		branchCode = value
	}
	if v, ok := c.Get(middleware.TerminalCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		terminalCode = value
	}

	rows, err := h.service.GetDataOp(c.Request.Context(), branchCode, terminalCode, input)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", rows)
}

// GetDetailOp godoc
// @Summary      Get Operation Plan Detail
// @Description  Loading/unloading plan header and detail rows by plan_code
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.GetDetailOpInput true "Detail payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/getDetailOp [post]
func (h *OpsPlanHandler) GetDetailOp(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input GetDetailOpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	branchCode, terminalCode := 0, 0
	if v, ok := c.Get(middleware.BranchCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		branchCode = value
	}
	if v, ok := c.Get(middleware.TerminalCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		terminalCode = value
	}

	res, err := h.service.GetDetailOp(c.Request.Context(), branchCode, terminalCode, input.PlanIdentifier())
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

// GetDetailDetermination godoc
// @Summary      Get Loading Unloading Determination Detail
// @Description  Loading/unloading determination header and detail rows by determination_code
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.GetDetailDeterminationInput true "Detail payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/getDetailDetermination [post]
func (h *OpsPlanHandler) GetDetailDetermination(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input GetDetailDeterminationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	branchCode, terminalCode := 0, 0
	if v, ok := c.Get(middleware.BranchCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid branch code in token")
			return
		}
		branchCode = value
	}
	if v, ok := c.Get(middleware.TerminalCodeKey); ok {
		value, err := parseContextInt(v)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid terminal code in token")
			return
		}
		terminalCode = value
	}

	res, err := h.service.GetDetailDetermination(c.Request.Context(), branchCode, terminalCode, input)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

// GetDataVesselSchedule godoc
// @Summary      Get Vessel Schedule Data
// @Description  Vessel schedule from post.vessel_schedule by pkk_number first, then vessel_code when pkk_number is not found
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.GetDataVesselLookupInput true "Request payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/getDataVesselSchedule [post]
func (h *OpsPlanHandler) GetDataVesselSchedule(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input GetDataVesselLookupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ppkNumber := strings.TrimSpace(input.RequestNumber())
	vesselCode := strings.TrimSpace(input.VesselCode)
	if ppkNumber == "" && vesselCode == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "pkk_number or vessel_code is required")
		return
	}

	rows, err := h.service.GetDataVesselSchedule(c.Request.Context(), ppkNumber, vesselCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", rows)
}

// GetDataVesel godoc
// @Summary      Get Vessel Master Data
// @Description  Vessel master from adm.posm_vessel and adm.adm_vessel_d by vessel_code
// @Tags         plan-op
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body op.GetDataVeselInput true "Request payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/op/getDataVesel [post]
func (h *OpsPlanHandler) GetDataVesel(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "ops plan service is not initialized")
		return
	}

	var input GetDataVeselInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	vesselCode := strings.TrimSpace(input.VesselCode)
	if vesselCode == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "vessel_code is required")
		return
	}

	rows, err := h.service.GetDataVesel(c.Request.Context(), vesselCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", rows)
}

func parseContextString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case *string:
		if v == nil {
			return "", strconv.ErrSyntax
		}
		return *v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case *int64:
		if v == nil {
			return "", strconv.ErrSyntax
		}
		return strconv.FormatInt(*v, 10), nil
	default:
		return "", strconv.ErrSyntax
	}
}

func parseContextInt(value interface{}) (int, error) {
	text, err := parseContextString(value)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(text)
}

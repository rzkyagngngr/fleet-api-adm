package postrequest

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type PostRequestHandler struct{ service PostRequestService }

func NewPostRequestHandler(service PostRequestService) *PostRequestHandler {
	return &PostRequestHandler{service: service}
}

// Search godoc
// @Summary      Search Permohonan Jasa Barang
// @Description  Paginated list of cargo service requests filtered by branch/terminal from JWT context
// @Tags         plan-permohonan-barang
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body postrequest.SearchPostRequestInput true "Search payload"
// @Success      200 {object} helper.MetaResponse
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/request/barang/search [post]
func (h *PostRequestHandler) Search(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "post request service is not initialized")
		return
	}

	var input SearchPostRequestInput
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

	rows, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
}

// Create godoc
// @Summary      Create Permohonan Jasa Barang
// @Description  Create a new cargo service request with manifest detail lines
// @Tags         plan-permohonan-barang
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body postrequest.CreatePostRequestInput true "Create payload"
// @Success      201 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/request/barang [post]
func (h *PostRequestHandler) Create(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "post request service is not initialized")
		return
	}

	var input CreatePostRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	identity := helper.ExtractIdentity(c)

	res, err := h.service.Create(c.Request.Context(), &input, identity)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "permohonan jasa barang created successfully", res)
}

// GetByID godoc
// @Summary      Get Permohonan Jasa Barang Detail
// @Description  Get a single cargo service request with all its manifest lines
// @Tags         plan-permohonan-barang
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Request ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/request/barang/{id} [get]
// @Router       /plan/request/barang [get]
func (h *PostRequestHandler) GetByID(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "post request service is not initialized")
		return
	}

	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid request ID")
		return
	}
	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

// Update godoc
// @Summary      Update Permohonan Jasa Barang
// @Description  Partially update a cargo service request. Providing 'details' will replace all manifest lines.
// @Tags         plan-permohonan-barang
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Request ID"
// @Param        payload body postrequest.UpdatePostRequestInput true "Update payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/request/barang/{id} [put]
// @Router       /plan/request/barang [put]
func (h *PostRequestHandler) Update(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "post request service is not initialized")
		return
	}

	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid request ID")
		return
	}

	var input UpdatePostRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	identity := helper.ExtractIdentity(c)
	res, err := h.service.Update(c.Request.Context(), id, &input, identity)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "permohonan jasa barang updated successfully", res)
}

// Delete godoc
// @Summary      Delete Permohonan Jasa Barang
// @Description  Delete a cargo service request and all its manifest lines
// @Tags         plan-permohonan-barang
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Request ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/request/barang/{id} [delete]
// @Router       /plan/request/barang [delete]
func (h *PostRequestHandler) Delete(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "post request service is not initialized")
		return
	}

	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid request ID")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "permohonan jasa barang deleted successfully", nil)
}

// GetStats godoc
// @Summary      Permohonan Jasa Barang Stats
// @Description  Return aggregated counts (total, pending, approved, rejected) for the active terminal
// @Tags         plan-permohonan-barang
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/request/barang/stats [get]
func (h *PostRequestHandler) GetStats(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "post request service is not initialized")
		return
	}

	identity := helper.ExtractIdentity(c)

	res, err := h.service.GetStats(c.Request.Context(), identity)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

// UpdateStatus godoc
// @Summary      Update Permohonan Jasa Barang Status
// @Description  Approve or Reject a cargo service request with optional remarks
// @Tags         plan-permohonan-barang
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Request ID"
// @Param        payload body postrequest.UpdateStatusInput true "Status update payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/request/barang/{id}/status [put]
// @Router       /plan/request/barang/status [put]
func (h *PostRequestHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid request ID")
		return
	}

	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	identity := helper.ExtractIdentity(c)
	if err := h.service.UpdateStatus(c.Request.Context(), id, input.Status, input.Remarks, identity); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "status updated successfully", nil)
}

// ─────────────────────────────────────────────────────────────
// HELPER
// ─────────────────────────────────────────────────────────────

// extractContext is now deprecated in favor of helper.ExtractIdentity

// SearchVesselSchedule godoc
// @Summary      Search Vessel Schedules
// @Description  Paginated list of vessel schedules filtered by branch/terminal from JWT context
// @Tags         plan-vessel-schedule
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body postrequest.SearchPostRequestInput true "Search payload"
// @Success      200 {object} helper.MetaResponse
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
func (h *PostRequestHandler) SearchVesselSchedule(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "post request service is not initialized")
		return
	}

	var input SearchPostRequestInput // Reusing SearchPostRequestInput for pagination structure
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

	rows, meta, err := h.service.SearchVesselSchedule(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
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

// GetVesselScheduleByID godoc
// @Summary      Get Vessel Schedule Detail
// @Description  Get a single vessel schedule entry
// @Tags         plan-vessel-schedule
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Schedule ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Failure      500 {object} helper.Response
func (h *PostRequestHandler) GetVesselScheduleByID(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "post request service is not initialized")
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid ID")
		return
	}
	res, err := h.service.GetVesselScheduleByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

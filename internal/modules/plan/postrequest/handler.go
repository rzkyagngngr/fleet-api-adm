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
	var input SearchPostRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Filters == nil {
		input.Filters = make(map[string]string)
	}
	if bc, ok := c.Get(middleware.BranchCodeKey); ok {
		input.Filters["branch_code"] = bc.(string)
	}
	if tc, ok := c.Get(middleware.TerminalCodeKey); ok {
		input.Filters["terminal_code"] = tc.(string)
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
	var input CreatePostRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	branchCode, terminalCode, branchName, terminalName, employeeID := extractContext(c)

	res, err := h.service.Create(c.Request.Context(), &input, branchCode, terminalCode, branchName, terminalName, employeeID)
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
func (h *PostRequestHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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
func (h *PostRequestHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid request ID")
		return
	}

	var input UpdatePostRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)
	res, err := h.service.Update(c.Request.Context(), id, &input, employeeID.(string))
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
func (h *PostRequestHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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
	branchCode, terminalCode := 0, 0
	if bc, ok := c.Get(middleware.BranchCodeKey); ok {
		branchCode, _ = strconv.Atoi(bc.(string))
	}
	if tc, ok := c.Get(middleware.TerminalCodeKey); ok {
		terminalCode, _ = strconv.Atoi(tc.(string))
	}

	res, err := h.service.GetStats(c.Request.Context(), branchCode, terminalCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

func (h *PostRequestHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid request ID")
		return
	}

	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	_, _, _, _, employeeID := extractContext(c)
	if err := h.service.UpdateStatus(c.Request.Context(), id, input.Status, input.Remarks, employeeID); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "status updated successfully", nil)
}

// ─────────────────────────────────────────────────────────────
// HELPER
// ─────────────────────────────────────────────────────────────

func extractContext(c *gin.Context) (branchCode, terminalCode int, branchName, terminalName, employeeID string) {

	if v, ok := c.Get(middleware.BranchCodeKey); ok {
		branchCode, _ = strconv.Atoi(v.(string))
	}
	if v, ok := c.Get(middleware.TerminalCodeKey); ok {
		terminalCode, _ = strconv.Atoi(v.(string))
	}
	// branch_name and terminal_name are optional — can be passed as query params
	branchName = c.DefaultQuery("branch_name", "")
	terminalName = c.DefaultQuery("terminal_name", "")
	if v, ok := c.Get(middleware.EmployeeIDKey); ok {
		employeeID, _ = v.(string)
	}
	return
}

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
// @Router       /plan/vessel-schedule/search [post]
func (h *PostRequestHandler) SearchVesselSchedule(c *gin.Context) {
	var input SearchPostRequestInput // Reusing SearchPostRequestInput for pagination structure
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Filters == nil {
		input.Filters = make(map[string]string)
	}
	if bc, ok := c.Get(middleware.BranchCodeKey); ok {
		input.Filters["branch_code"] = bc.(string)
	}
	if tc, ok := c.Get(middleware.TerminalCodeKey); ok {
		input.Filters["terminal_code"] = tc.(string)
	}

	rows, meta, err := h.service.SearchVesselSchedule(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
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
// @Router       /plan/vessel-schedule/{id} [get]
func (h *PostRequestHandler) GetVesselScheduleByID(c *gin.Context) {
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


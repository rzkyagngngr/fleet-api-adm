package vesselrpkmanual

import (
	"net/http"
	"omniport-api/internal/helper"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VesselRpkHandler struct {
	service VesselRpkService
}

func NewVesselRpkHandler(service VesselRpkService) *VesselRpkHandler {
	return &VesselRpkHandler{service: service}
}

// Create godoc
// @Summary      Create Vessel RPK Manual
// @Description  Create a new manual vessel RPK
// @Tags         plan-vessel-rpk-manual
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body vesselrpkmanual.CreateVesselRpkInput true "Create payload"
// @Success      201 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/vessel-rpk-manual [post]
func (h *VesselRpkHandler) Create(c *gin.Context) {
	var input CreateVesselRpkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	branchCode := c.GetInt64("branch_code")
	if branchCode == 0 && input.BranchCode > 0 {
		branchCode = input.BranchCode
	}
	terminalCode := c.GetInt64("terminal_code")
	if terminalCode == 0 && input.TerminalCode > 0 {
		terminalCode = input.TerminalCode
	}
	userID := c.GetString("employee_id")

	res, err := h.service.Create(c.Request.Context(), input, branchCode, terminalCode, userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "Vessel RPK created successfully", res)
}

// GetByID godoc
// @Summary      Get Vessel RPK Manual Detail
// @Description  Get a manual vessel RPK by id
// @Tags         plan-vessel-rpk-manual
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "Vessel RPK Manual ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Router       /plan/vessel-rpk-manual [get]
func (h *VesselRpkHandler) GetByID(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		idStr = c.Param("id")
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Vessel RPK not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Success", res)
}

// Search godoc
// @Summary      Search Vessel RPK Manual
// @Description  Paginated list of manual vessel RPK records
// @Tags         plan-vessel-rpk-manual
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body vesselrpkmanual.SearchVesselRpkRequest true "Search payload"
// @Success      200 {object} helper.MetaResponse
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/vessel-rpk-manual/search [post]
func (h *VesselRpkHandler) Search(c *gin.Context) {
	if h == nil || h.service == nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "vessel rpk service is not initialized")
		return
	}

	var req SearchVesselRpkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 200 {
		req.Limit = 200
	}

	branchCode := c.GetInt64("branch_code")
	terminalCode := c.GetInt64("terminal_code")

	list, meta, err := h.service.List(c.Request.Context(), branchCode, terminalCode, req.Page, req.Limit, req.Search, req.Filters)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "Success", list, meta)
}

// Update godoc
// @Summary      Update Vessel RPK Manual
// @Description  Update a manual vessel RPK by id
// @Tags         plan-vessel-rpk-manual
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "Vessel RPK Manual ID"
// @Param        payload body vesselrpkmanual.CreateVesselRpkInput true "Update payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/vessel-rpk-manual [put]
func (h *VesselRpkHandler) Update(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		idStr = c.Param("id")
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var input CreateVesselRpkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("employee_id")

	if err := h.service.Update(c.Request.Context(), id, input, userID); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Vessel RPK updated successfully", nil)
}

// Delete godoc
// @Summary      Delete Vessel RPK Manual
// @Description  Delete a manual vessel RPK by id
// @Tags         plan-vessel-rpk-manual
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "Vessel RPK Manual ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /plan/vessel-rpk-manual [delete]
func (h *VesselRpkHandler) Delete(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		idStr = c.Param("id")
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Vessel RPK deleted successfully", nil)
}

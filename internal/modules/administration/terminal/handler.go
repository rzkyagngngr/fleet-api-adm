package terminal

import (
	"context"
	"net/http"
	"strconv"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserProvider interface {
	GetProfile(ctx context.Context, userID uint64) (any, error)
}

type TerminalHandler struct {
	service      TerminalService
	userProvider UserProvider
}

func NewTerminalHandler(service TerminalService, userProvider UserProvider) *TerminalHandler {
	return &TerminalHandler{
		service:      service,
		userProvider: userProvider,
	}
}

func (h *TerminalHandler) getCompanyInfo(c *gin.Context) (string, string, string, error) {
	userID := middleware.GetUserID(c)
	res, err := h.userProvider.GetProfile(c.Request.Context(), userID)
	if err != nil {
		return "", "", "", err
	}

	var compCode, compName string
	if m, ok := res.(interface{ GetCompanyData() (string, string) }); ok {
		compCode, compName = m.GetCompanyData()
	} else {
		return "", "", "", err
	}

	empID, _ := c.Get(middleware.EmployeeIDKey)
	return compCode, compName, empID.(string), nil
}

// Search godoc
// @Summary Search terminals
// @Description Search and filter terminals
// @Tags master-terminal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body terminal.SearchTerminalRequest true "Search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/terminals/search [post]
func (h *TerminalHandler) Search(c *gin.Context) {
	var input SearchTerminalRequest
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
// @Summary Create terminal
// @Description Create a new terminal
// @Tags master-terminal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body terminal.TerminalRequest true "Terminal payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/terminals [post]
func (h *TerminalHandler) Create(c *gin.Context) {
	var input TerminalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	compCode, compName, empID, err := h.getCompanyInfo(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "failed to get user context")
		return
	}

	if err := h.service.Create(c.Request.Context(), &input, compCode, compName, empID); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "terminal created successfully", nil)
}

// Update godoc
// @Summary Update terminal
// @Description Update terminal by id
// @Tags master-terminal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Terminal ID"
// @Param payload body terminal.TerminalRequest true "Terminal payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/terminals [put]
func (h *TerminalHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	var input TerminalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	_, _, empID, err := h.getCompanyInfo(c)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "failed to get user context")
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &input, empID); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "terminal updated successfully", nil)
}

// Delete godoc
// @Summary Delete terminal
// @Description Delete terminal by id
// @Tags master-terminal
// @Produce json
// @Security BearerAuth
// @Param id query int true "Terminal ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/terminals [delete]
func (h *TerminalHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "terminal deleted successfully", nil)
}

// GetStats godoc
// @Summary Get terminal statistics
// @Description Get aggregated terminal statistics
// @Tags master-terminal
// @Produce json
// @Security BearerAuth
// @Param company_code query string false "Filter by company code"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/terminals/stats [get]
func (h *TerminalHandler) GetStats(c *gin.Context) {
	compCode, _, _, err := h.getCompanyInfo(c)
	if err != nil {
		// Fallback to query param for superusers or special cases
		compCode = c.Query("company_code")
	}

	res, err := h.service.GetStats(c.Request.Context(), compCode)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

// GetByID godoc
// @Summary Get terminal by id
// @Description Get terminal detail by id
// @Tags master-terminal
// @Produce json
// @Security BearerAuth
// @Param id query int true "Terminal ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/terminals [get]
func (h *TerminalHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "terminal not found")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

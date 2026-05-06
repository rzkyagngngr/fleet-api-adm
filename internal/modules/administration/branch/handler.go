package branch

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

type BranchHandler struct {
	service      BranchService
	userProvider UserProvider
}

func NewBranchHandler(service BranchService, userProvider UserProvider) *BranchHandler {
	return &BranchHandler{
		service:      service,
		userProvider: userProvider,
	}
}

func (h *BranchHandler) getCompanyInfo(c *gin.Context) (string, string, string, error) {
	userID := middleware.GetUserID(c)
	res, err := h.userProvider.GetProfile(c.Request.Context(), userID)
	if err != nil {
		return "", "", "", err
	}

	// Use reflection or type assertion to get fields without importing user package
	// Since we know the underlying type will be UserResponse
	var compCode, compName string
	if m, ok := res.(interface{ GetCompanyData() (string, string) }); ok {
		compCode, compName = m.GetCompanyData()
	} else {
		// Fallback or handle error
		return "", "", "", err
	}

	empID, _ := c.Get(middleware.EmployeeIDKey)
	return compCode, compName, empID.(string), nil
}

// Search godoc
// @Summary Search branches
// @Description Search and filter branches
// @Tags master-branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body branch.SearchBranchRequest true "Search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/branches/search [post]
func (h *BranchHandler) Search(c *gin.Context) {
	var input SearchBranchRequest
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
// @Summary Create branch
// @Description Create a new branch
// @Tags master-branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body branch.BranchRequest true "Branch payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/branches [post]
func (h *BranchHandler) Create(c *gin.Context) {
	var input BranchRequest
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
	helper.SuccessResponse(c, http.StatusCreated, "branch created successfully", nil)
}

// Update godoc
// @Summary Update branch
// @Description Update branch by id
// @Tags master-branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Branch ID"
// @Param payload body branch.BranchRequest true "Branch payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/branches [put]
func (h *BranchHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	var input BranchRequest
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
	helper.SuccessResponse(c, http.StatusOK, "branch updated successfully", nil)
}

// Delete godoc
// @Summary Delete branch
// @Description Delete branch by id
// @Tags master-branch
// @Produce json
// @Security BearerAuth
// @Param id query int true "Branch ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/branches [delete]
func (h *BranchHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "branch deleted successfully", nil)
}

// GetStats godoc
// @Summary Get branch statistics
// @Description Get aggregated branch statistics
// @Tags master-branch
// @Produce json
// @Security BearerAuth
// @Param company_code query string false "Filter by company code"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/branches/stats [get]
func (h *BranchHandler) GetStats(c *gin.Context) {
	compCode, _, _, err := h.getCompanyInfo(c)
	if err != nil {
		// If fails to get company info (e.g. superuser without company),
		// fallback to query param or just show empty/all
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
// @Summary Get branch by id
// @Description Get branch detail by id
// @Tags master-branch
// @Produce json
// @Security BearerAuth
// @Param id query int true "Branch ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/branches [get]
func (h *BranchHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "branch not found")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

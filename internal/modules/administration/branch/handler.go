package branch

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"omniport-api/internal/modules/administration/user"

	"github.com/gin-gonic/gin"
)

type BranchHandler struct {
	service     BranchService
	userService user.UserService
}

func NewBranchHandler(service BranchService, userService user.UserService) *BranchHandler {
	return &BranchHandler{
		service:     service,
		userService: userService,
	}
}

func (h *BranchHandler) getCompanyInfo(c *gin.Context) (string, string, string, error) {
	userID := middleware.GetUserID(c)
	profile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		return "", "", "", err
	}
	empID, _ := c.Get(middleware.EmployeeIDKey)
	return profile.CompanyCode, profile.CompanyName, empID.(string), nil
}

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

func (h *BranchHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
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

func (h *BranchHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "branch deleted successfully", nil)
}

func (h *BranchHandler) GetStats(c *gin.Context) {
	res, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

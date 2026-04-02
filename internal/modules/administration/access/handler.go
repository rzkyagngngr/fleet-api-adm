package access

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"

	"github.com/gin-gonic/gin"
)

type AccessHandler struct{ accessService AccessService }

func NewAccessHandler(accessService AccessService) *AccessHandler {
	return &AccessHandler{accessService: accessService}
}

// GetRoleAccess godoc
// @Summary Get role access
// @Description Retrieve access matrix for a role
// @Tags master-role-access
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/roles/{id}/access [get]
func (h *AccessHandler) GetRoleAccess(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}
	list, err := h.accessService.GetRoleAccess(c.Request.Context(), roleID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to get role access")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "role access retrieved successfully", list)
}

// UpdateRoleAccess godoc
// @Summary Update role access
// @Description Update access matrix for a role
// @Tags master-role-access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param payload body []access.Access true "Role access payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/roles/{id}/access [post]
func (h *AccessHandler) UpdateRoleAccess(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}
	var accessList []Access
	if err := c.ShouldBindJSON(&accessList); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.accessService.UpdateRoleAccess(c.Request.Context(), roleID, accessList); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update role access")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "role access updated successfully", nil)
}

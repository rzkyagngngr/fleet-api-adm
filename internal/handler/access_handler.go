package handler

import (
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/service"
	"gin-boilerplate/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccessHandler struct {
	accessService service.AccessService
}

func NewAccessHandler(accessService service.AccessService) *AccessHandler {
	return &AccessHandler{accessService: accessService}
}

func (h *AccessHandler) GetRoleAccess(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}

	accessList, err := h.accessService.GetRoleAccess(c.Request.Context(), roleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to get role access")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "role access retrieved successfully", accessList)
}

func (h *AccessHandler) UpdateRoleAccess(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}

	var accessList []entity.Access
	if err := c.ShouldBindJSON(&accessList); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.accessService.UpdateRoleAccess(c.Request.Context(), roleID, accessList); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update role access")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "role access updated successfully", nil)
}

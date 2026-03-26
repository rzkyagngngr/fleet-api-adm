package handler

import (
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/service"
	"gin-boilerplate/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleService.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to get roles")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "roles retrieved successfully", roles)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role entity.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.roleService.Create(&role); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create role")
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "role created successfully", role)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}

	var role entity.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.roleService.Update(id, &role); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update role")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "role updated successfully", role)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}

	if err := h.roleService.Delete(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete role")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "role deleted successfully", nil)
}

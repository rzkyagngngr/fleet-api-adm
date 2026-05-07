package role

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct{ roleService RoleService }

func NewRoleHandler(roleService RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// GetAllRoles godoc
// @Summary Get all roles
// @Description Retrieve all roles
// @Tags master-roles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/roles [get]
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleService.FindAll(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to get roles")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "roles retrieved successfully", roles)
}

// GetRoleDetail godoc
// @Summary Get role detail
// @Description Retrieve role detail by id
// @Tags master-roles
// @Produce json
// @Security BearerAuth
// @Param id query int true "Role ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /master/roles [get]
func (h *RoleHandler) GetRoleDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}

	role, err := h.roleService.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "role not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "role detail retrieved successfully", role)
}

// CreateRole godoc
// @Summary Create role
// @Description Create a new role
// @Tags master-roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body role.Role true "Role payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var r Role
	if err := c.ShouldBindJSON(&r); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.roleService.Create(c.Request.Context(), &r); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create role")
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "role created successfully", r)
}

// UpdateRole godoc
// @Summary Update role
// @Description Update role by id
// @Tags master-roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Role ID"
// @Param payload body role.Role true "Role payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/roles [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}
	var r Role
	if err := c.ShouldBindJSON(&r); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.roleService.Update(c.Request.Context(), id, &r); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update role")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "role updated successfully", r)
}

// DeleteRole godoc
// @Summary Delete role
// @Description Delete role by id
// @Tags master-roles
// @Produce json
// @Security BearerAuth
// @Param id query int true "Role ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/roles [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid role id")
		return
	}
	if err := h.roleService.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete role")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "role deleted successfully", nil)
}

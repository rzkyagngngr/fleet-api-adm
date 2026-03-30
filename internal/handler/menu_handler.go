package handler

import (
	"gin-boilerplate/internal/model/entity"
	"gin-boilerplate/internal/service"
	"gin-boilerplate/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuService service.MenuService
}

func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

func (h *MenuHandler) GetAllMenus(c *gin.Context) {
	menus, err := h.menuService.FindAll(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to get menus")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "menus retrieved successfully", menus)
}

func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var menu entity.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.menuService.Create(c.Request.Context(), &menu); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create menu")
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "menu created successfully", menu)
}

func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid menu id")
		return
	}

	var menu entity.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.menuService.Update(c.Request.Context(), id, &menu); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update menu")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "menu updated successfully", menu)
}

func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid menu id")
		return
	}

	if err := h.menuService.Delete(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete menu")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "menu deleted successfully", nil)
}

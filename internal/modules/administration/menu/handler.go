package menu

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct{ menuService MenuService }

func NewMenuHandler(menuService MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

// GetAllMenus godoc
// @Summary Get all menus
// @Description Retrieve all menu records
// @Tags master-menus
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/menus [get]
func (h *MenuHandler) GetAllMenus(c *gin.Context) {
	menus, err := h.menuService.FindAll(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to get menus")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "menus retrieved successfully", menus)
}

// GetMenuDetail godoc
// @Summary Get menu detail
// @Description Retrieve menu detail by id
// @Tags master-menus
// @Produce json
// @Security BearerAuth
// @Param id query int true "Menu ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /master/menus [get]
func (h *MenuHandler) GetMenuDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid menu id")
		return
	}

	menu, err := h.menuService.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "menu not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "menu detail retrieved successfully", menu)
}

// SearchMenus godoc
// @Summary Search menus
// @Description Retrieve menus with server-side pagination, filtering, sorting, and download range support
// @Tags master-menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body menu.SearchMenusRequest true "Menu search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/menus/search [post]
func (h *MenuHandler) SearchMenus(c *gin.Context) {
	var input SearchMenusRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	menus, meta, err := h.menuService.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search menus")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "menus retrieved successfully", menus, meta)
}

// CreateMenu godoc
// @Summary Create menu
// @Description Create a new menu record
// @Tags master-menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body menu.Menu true "Menu payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var m Menu
	if err := c.ShouldBindJSON(&m); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.menuService.Create(c.Request.Context(), &m); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create menu")
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "menu created successfully", m)
}

// UpdateMenu godoc
// @Summary Update menu
// @Description Update an existing menu by id
// @Tags master-menus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Menu ID"
// @Param payload body menu.Menu true "Menu payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/menus [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid menu id")
		return
	}
	var m Menu
	if err := c.ShouldBindJSON(&m); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.menuService.Update(c.Request.Context(), id, &m); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update menu")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "menu updated successfully", m)
}

// DeleteMenu godoc
// @Summary Delete menu
// @Description Delete menu by id
// @Tags master-menus
// @Produce json
// @Security BearerAuth
// @Param id query int true "Menu ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/menus [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid menu id")
		return
	}
	if err := h.menuService.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete menu")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "menu deleted successfully", nil)
}

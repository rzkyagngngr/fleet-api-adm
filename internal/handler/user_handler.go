package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gin-boilerplate/internal/middleware"
	"gin-boilerplate/internal/service"
	"gin-boilerplate/pkg/utils"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to get profile")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "profile retrieved successfully", profile)
}

package user

import (
	"net/http"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{ userService UserService }

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile godoc
// @Summary Get profile
// @Description Retrieve profile for authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		helper.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	profile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to get profile")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "profile retrieved successfully", profile)
}

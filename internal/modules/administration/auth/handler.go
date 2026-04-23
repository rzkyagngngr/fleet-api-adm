package auth

import (
	"net/http"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ authService AuthService }

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body auth.UserRegisterRequest true "Register payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}
	result, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "user registered successfully", result)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return token with menus
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body auth.LoginRequest true "Login payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}
	result, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "login successful", result)
}

// ChangeTerminal godoc
// @Summary Change terminal
// @Description Change branch and terminal for current authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body auth.ChangeTerminalRequest true "Change terminal payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /auth/change-terminal [post]
func (h *AuthHandler) ChangeTerminal(c *gin.Context) {
	var req ChangeTerminalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.authService.ChangeTerminal(c.Request.Context(), userID.(uint64), &req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "terminal changed successfully", result)
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Refresh current authenticated user token
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), userID.(uint64))
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "token refreshed successfully", result)
}

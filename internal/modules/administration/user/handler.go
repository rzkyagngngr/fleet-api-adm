package user

import (
	"net/http"
	"strconv"

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

// FindAll godoc
// @Summary Get all users
// @Description Retrieve paginated user list
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} helper.PaginatedResponse
// @Failure 500 {object} helper.Response
// @Router /master/users [get]
func (h *UserHandler) FindAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	rows, total, err := h.userService.FindAll(c.Request.Context(), page, size)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.PaginatedSuccessResponse(c, http.StatusOK, "users retrieved successfully", rows, total, page, size)
}

// FindByID godoc
// @Summary Get user by id
// @Description Retrieve user detail by id
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/users/{id} [get]
func (h *UserHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid user ID")
		return
	}
	u, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "user retrieved successfully", u)
}

// Create godoc
// @Summary Create user
// @Description Create a new master user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body user.UserRequest true "User payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var input UserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}
	
	if input.Password == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "password is required")
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)

	if err := h.userService.CreateUser(c.Request.Context(), &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "user created successfully", nil)
}

// Update godoc
// @Summary Update user
// @Description Update user by id
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param payload body user.UserRequest true "User payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	var input UserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)

	if err := h.userService.UpdateUser(c.Request.Context(), id, &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "user updated successfully", nil)
}

// Delete godoc
// @Summary Delete user
// @Description Delete user by id
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 401 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "user deleted successfully", nil)
}

// Search godoc
// @Summary Search users
// @Description Search and filter users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body user.SearchUsersRequest true "Search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/users/search [post]
func (h *UserHandler) Search(c *gin.Context) {
	var input SearchUsersRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rows, meta, err := h.userService.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
}

// GetStats godoc
// @Summary Get user statistics
// @Description Get aggregated user statistics
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/users/stats [get]
func (h *UserHandler) GetStats(c *gin.Context) {
	res, err := h.userService.GetUserStats(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "success", res)
}


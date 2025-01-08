package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samzong/share-ai-platform/internal/middleware"
	"github.com/samzong/share-ai-platform/internal/services"
	"github.com/samzong/share-ai-platform/internal/models"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "Registration details"
// @Success 200 {object} services.UserResponse
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login credentials"
// @Success 200 {object} services.UserResponse
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate user's token
// @Tags auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "message: Successfully logged out"
// @Failure 500 {object} map[string]interface{} "error message"
// @Router /auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.userService.Logout(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} services.UserResponse
// @Failure 500 {object} map[string]interface{} "error message"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags users
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param nickname formData string false "User nickname"
// @Param avatar formData file false "User avatar"
// @Success 200 {object} services.UserResponse
// @Failure 400,500 {object} map[string]interface{} "error message"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req services.UpdateProfileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user's username and email
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "Update user request"
// @Success 200 {object} services.UserResponse
// @Failure 400,500 {object} map[string]interface{} "error message"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" example:"johndoe"`
		Email    string `json:"email" example:"john@example.com"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Param("id")
	user, err := h.userService.UpdateUser(userID, req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserRole godoc
// @Summary Update user role
// @Description Update user's role (admin only)
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRoleRequest true "Update role request"
// @Success 204 "No Content"
// @Failure 400,403,500 {object} map[string]interface{} "error message"
// @Router /users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	var req struct {
		Role models.Role `json:"role" example:"admin" enums:"user,admin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证角色是否有效
	if !models.IsValidRole(req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	adminID := middleware.GetUserID(c)
	userID := c.Param("id")

	if err := h.userService.UpdateUserRole(adminID, userID, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUsers godoc
// @Summary List users
// @Description Get paginated list of users (admin only)
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Param page query int true "Page number" minimum(1)
// @Param page_size query int true "Page size" minimum(1) maximum(100)
// @Success 200 {object} ListUsersResponse
// @Failure 400,403,500 {object} map[string]interface{} "error message"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req services.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userService.ListUsers(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Request/Response models for Swagger documentation
type UpdateUserRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
}

type UpdateUserRoleRequest struct {
	Role models.Role `json:"role" example:"admin" enums:"user,admin"`
}

type ListUsersResponse struct {
	Total int64                  `json:"total" example:"100"`
	Users []services.UserResponse `json:"users"`
} 
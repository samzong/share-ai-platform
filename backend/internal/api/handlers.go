package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samzong/share-ai-platform/internal/middleware"
	"github.com/samzong/share-ai-platform/internal/services"
)

type Handler struct {
	userService   *services.UserService
	imageService  *services.ImageService
	deployService *services.DeployService
}

func NewHandler() *Handler {
	return &Handler{
		userService:   services.NewUserService(),
		imageService:  services.NewImageService(),
		deployService: services.NewDeployService(),
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/api/users/register", h.Register)
	r.POST("/api/users/login", h.Login)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// User routes
		api.GET("/users/me", h.GetCurrentUser)
		api.PUT("/users/me", h.UpdateUser)

		// Image routes
		api.GET("/images", h.ListImages)
		api.GET("/images/:id", h.GetImage)
		api.POST("/images/:id/collect", h.CollectImage)
		api.DELETE("/images/:id/collect", h.UncollectImage)

		// Deploy routes
		api.GET("/deploy/:id", h.GetDeployInfo)
		api.POST("/deploy/:id", h.Deploy)
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetCurrentUser returns the current user's information
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates the current user's information
func (h *Handler) UpdateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userService.UpdateUser(userID, req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListImages returns a list of images
func (h *Handler) ListImages(c *gin.Context) {
	var req services.ImageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	images, total, err := h.imageService.ListImages(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  images,
		"total": total,
	})
}

// GetImage returns a single image by ID
func (h *Handler) GetImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	image, err := h.imageService.GetImageByID(c.Request.Context(), imageID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

// CollectImage adds an image to user's collection
func (h *Handler) CollectImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.imageService.CollectImage(userID, imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UncollectImage removes an image from user's collection
func (h *Handler) UncollectImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.imageService.UncollectImage(userID, imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetDeployInfo returns deployment information for an image
func (h *Handler) GetDeployInfo(c *gin.Context) {
	imageID := c.Param("id")
	providerID := c.Query("provider_id")

	info, err := h.deployService.GetDeployInfo(imageID, providerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// Deploy handles image deployment
func (h *Handler) Deploy(c *gin.Context) {
	var req services.DeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ImageID = c.Param("id")

	info, err := h.deployService.Deploy(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
} 
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samzong/share-ai-platform/internal/middleware"
	"github.com/samzong/share-ai-platform/internal/services"
)

type ImageHandler struct {
	imageService *services.ImageService
}

func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		imageService: services.NewImageService(),
	}
}

// ListImages 返回镜像列表
func (h *ImageHandler) ListImages(c *gin.Context) {
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

// GetImage 返回单个镜像信息
func (h *ImageHandler) GetImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	image, err := h.imageService.GetImageByID(c.Request.Context(), imageID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

// CreateImage 创建新镜像
func (h *ImageHandler) CreateImage(c *gin.Context) {
	var req services.CreateImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	image, err := h.imageService.CreateImage(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, image)
}

// UpdateImage 更新镜像信息
func (h *ImageHandler) UpdateImage(c *gin.Context) {
	var req services.UpdateImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imageID := c.Param("id")
	userID := middleware.GetUserID(c)
	image, err := h.imageService.UpdateImage(c.Request.Context(), imageID, &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, image)
}

// DeleteImage 删除镜像
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.imageService.DeleteImage(c.Request.Context(), imageID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CollectImage 收藏镜像
func (h *ImageHandler) CollectImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.imageService.CollectImage(userID, imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UncollectImage 取消收藏镜像
func (h *ImageHandler) UncollectImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.imageService.UncollectImage(userID, imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListFavorites 获取用户收藏的镜像列表
func (h *ImageHandler) ListFavorites(c *gin.Context) {
	var req services.ImageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	images, total, err := h.imageService.ListFavorites(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  images,
		"total": total,
	})
} 
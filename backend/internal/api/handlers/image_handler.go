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

// ListImages godoc
// @Summary 获取容器镜像列表
// @Description 获取所有可用的容器镜像列表，支持分页和搜索，包含镜像名称、标签、描述等信息
// @Tags container-images
// @Accept json
// @Produce json
// @Param page query int false "页码，默认 1"
// @Param page_size query int false "每页数量，默认 10"
// @Param search query string false "搜索关键词（镜像名称、描述）"
// @Success 200 {object} map[string]interface{} "data: []ContainerImage, total: int"
// @Failure 400 {object} map[string]interface{} "error message"
// @Failure 500 {object} map[string]interface{} "error message"
// @Router /images [get]
func (h *ImageHandler) ListImages(c *gin.Context) {
	var req services.ImageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户 ID（如果已登录）
	var userID string
	if id, exists := c.Get("user_id"); exists {
		userID = id.(string)
	}

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

// GetImage godoc
// @Summary 获取容器镜像详情
// @Description 根据镜像 ID 获取容器镜像的详细信息，包括镜像配置、版本、使用说明等
// @Tags container-images
// @Accept json
// @Produce json
// @Param id path string true "容器镜像 ID"
// @Success 200 {object} services.ImageResponse
// @Failure 404 {object} map[string]interface{} "error message"
// @Router /images/{id} [get]
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

// CreateImage godoc
// @Summary 创建容器镜像
// @Description 创建一个新的容器镜像，包括基本信息、配置参数、运行环境等详细信息
// @Tags container-images
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org_id path string true "组织ID，如果不指定则为 'public'"
// @Param request body services.CreateImageRequest true "镜像信息"
// @Success 201 {object} services.ImageResponse
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /orgs/{org_id}/images [post]
func (h *ImageHandler) CreateImage(c *gin.Context) {
	var req services.CreateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	orgID := c.Param("org_id")
	if orgID == "" {
		orgID = "public" // 如果不指定组织ID，则使用 "public"
	}

	image, err := h.imageService.CreateImage(c.Request.Context(), &req, userID, orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, image)
}

// UpdateImage godoc
// @Summary 更新容器镜像信息
// @Description 更新指定容器镜像的信息，包括基本信息、配置参数等
// @Tags container-images
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "容器镜像 ID"
// @Param request body services.UpdateImageRequest true "更新的镜像信息"
// @Success 200 {object} services.ImageResponse
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /images/{id} [put]
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

// DeleteImage godoc
// @Summary 删除容器镜像
// @Description 删除指定的容器镜像及其相关配置信息
// @Tags container-images
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "容器镜像 ID"
// @Success 200 {object} map[string]interface{} "success message"
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /images/{id} [delete]
func (h *ImageHandler) DeleteImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.imageService.DeleteImage(c.Request.Context(), imageID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CollectImage godoc
// @Summary 收藏容器镜像
// @Description 将指定的容器镜像添加到个人收藏夹中
// @Tags container-images
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "容器镜像 ID"
// @Success 200 {object} map[string]interface{} "success message"
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /images/{id}/collect [post]
func (h *ImageHandler) CollectImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.imageService.CollectImage(userID, imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UncollectImage godoc
// @Summary 取消收藏容器镜像
// @Description 将指定的容器镜像从个人收藏夹中移除
// @Tags container-images
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "容器镜像 ID"
// @Success 200 {object} map[string]interface{} "success message"
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /images/{id}/collect [delete]
func (h *ImageHandler) UncollectImage(c *gin.Context) {
	imageID := c.Param("id")
	userID := middleware.GetUserID(c)

	if err := h.imageService.UncollectImage(userID, imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListFavorites godoc
// @Summary 获取收藏的容器镜像列表
// @Description 获取当前用户收藏的所有容器镜像列表
// @Tags container-images
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "data: []ContainerImage"
// @Failure 400 {object} map[string]interface{} "error message"
// @Router /favorites [get]
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

package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/models"
	"github.com/samzong/share-ai-platform/internal/utils"
)

type ImageService struct{}

type ImageListRequest struct {
	Page     int      `form:"page" binding:"omitempty,min=1"`
	PageSize int      `form:"page_size" binding:"omitempty,min=1,max=100"`
	Search   string   `form:"search"`
	Labels   []string `form:"labels"`
	Sort     string   `form:"sort" binding:"oneof=stars created_at updated_at ''"`
}

type ImageResponse struct {
	ID          string    `json:"id"`          // 镜像唯一标识符
	OrgID       string    `json:"org_id"`      // 组织ID
	Name        string    `json:"name"`        // 镜像显示名称
	Description string    `json:"description"` // 镜像描述
	Author      string    `json:"author"`      // 创建者ID
	Registry    string    `json:"registry"`    // 镜像仓库服务器
	Namespace   string    `json:"namespace"`   // 命名空间/组织
	Repository  string    `json:"repository"`  // 镜像名称
	Tag         string    `json:"tag"`         // 版本标签
	Digest      string    `json:"digest"`      // 镜像内容哈希值
	Size        int64     `json:"size"`        // 镜像大小（字节）
	ReadmePath  string    `json:"readme_path"` // README文件路径
	Stars       int       `json:"stars"`       // 收藏数
	Visibility  string    `json:"visibility"`  // 可见性：public/private
	Platform    string    `json:"platform"`    // 平台架构
	Labels      []string  `json:"labels"`      // 标签列表，用于分类和搜索
	IsStarred   bool      `json:"is_starred"`  // 当前用户是否已收藏
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
}

type CreateImageRequest struct {
	Name        string                `json:"name" binding:"required"`
	Description string                `json:"description"`
	Registry    string                `json:"registry" binding:"required"`
	Namespace   string                `json:"namespace"`
	Repository  string                `json:"repository" binding:"required"`
	Tag         string                `json:"tag" binding:"required"`
	Digest      string                `json:"digest" binding:"required"`
	Size        int64                 `json:"size"`
	ReadmeFile  *multipart.FileHeader `json:"readme_file,omitempty"`
	Visibility  string                `json:"visibility" binding:"required,oneof=public private"`
	Platform    string                `json:"platform" binding:"required"`
	Labels      []string              `json:"labels,omitempty"`
}

type LayerInfo struct {
	Digest    string `json:"digest"`              // 层摘要
	Size      int64  `json:"size"`                // 层大小
	CreatedBy string `json:"created_by,omitempty"` // 创建该层的命令
}

type UpdateImageRequest struct {
	Name        string                `form:"name" json:"name,omitempty"`
	Description string                `form:"description" json:"description,omitempty"`
	Registry    string                `form:"registry" json:"registry,omitempty"`
	Namespace   string                `form:"namespace" json:"namespace,omitempty"`
	Repository  string                `form:"repository" json:"repository,omitempty"`
	Tag         string                `form:"tag" json:"tag,omitempty"`
	ReadmeFile  *multipart.FileHeader `form:"readme_file" json:"readme_file,omitempty"`
	Visibility  string                `form:"visibility" binding:"omitempty,oneof=public private" json:"visibility,omitempty"`
	Platform    string                `form:"platform" json:"platform,omitempty"`
	Labels      []string              `form:"labels" json:"labels,omitempty"`
}

// NewImageService creates a new ImageService
func NewImageService() *ImageService {
	return &ImageService{}
}

// ListImages retrieves a list of images with pagination and filtering
func (s *ImageService) ListImages(ctx context.Context, req *ImageListRequest, userID string) ([]ImageResponse, int64, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := database.GetDB()
	rdb := database.GetRedis()

	// Try to get from cache
	cacheKey := fmt.Sprintf("images:%d:%d:%s:%v:%s", req.Page, req.PageSize, req.Search, req.Labels, req.Sort)
	if cached, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
		var response []ImageResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return response, 0, nil // Total count might be inaccurate from cache
		}
	}

	// Build query
	query := db.Model(&models.Image{})

	// Apply search filter
	if req.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// Apply labels filter
	if len(req.Labels) > 0 {
		query = query.Joins("JOIN image_labels ON images.id = image_labels.image_id").
			Where("image_labels.label IN ?", req.Labels).
			Group("images.id").
			Having("COUNT(DISTINCT image_labels.label) = ?", len(req.Labels))
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	switch req.Sort {
	case "stars":
		query = query.Order("stars DESC")
	case "created_at":
		query = query.Order("created_at DESC")
	case "updated_at":
		query = query.Order("updated_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// Execute query
	var images []models.Image
	if err := query.Find(&images).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response
	response := make([]ImageResponse, len(images))
	for i, img := range images {
		isStarred := false
		if userID != "" {
			var count int64
			if err := db.Model(&models.Collection{}).Where("user_id = ? AND image_id = ?", userID, img.ID).Count(&count).Error; err != nil {
				return nil, 0, err
			}
			isStarred = count > 0
		}

		response[i] = ImageResponse{
			ID:          img.ID,
			OrgID:       img.OrgID,
			Name:        img.Name,
			Description: img.Description,
			Author:      img.Author,
			Registry:    img.Registry,
			Namespace:   img.Namespace,
			Repository:  img.Repository,
			Tag:         img.Tag,
			Digest:      img.Digest,
			Size:        img.Size,
			ReadmePath:  img.ReadmePath,
			Stars:       img.Stars,
			Labels:      make([]string, len(img.Labels)),
			IsStarred:   isStarred,
			CreatedAt:   img.CreatedAt,
			UpdatedAt:   img.UpdatedAt,
			Visibility:  img.Visibility,
			Platform:    img.Platform,
		}

		for j, label := range img.Labels {
			response[i].Labels[j] = label.Name
		}
	}

	// Cache the results
	if len(response) > 0 {
		if cached, err := json.Marshal(response); err == nil {
			rdb.Set(ctx, cacheKey, cached, time.Minute*5)
		}
	}

	return response, total, nil
}

// GetImageByID retrieves an image by ID
func (s *ImageService) GetImageByID(ctx context.Context, id string, userID string) (*ImageResponse, error) {
	db := database.GetDB()

	var image models.Image
	if err := db.Preload("Labels").First(&image, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Check if user has starred the image
	var isStarred bool
	if userID != "" {
		var collection models.Collection
		isStarred = db.Where("user_id = ? AND image_id = ?", userID, id).First(&collection).Error == nil
	}

	response := &ImageResponse{
		ID:          image.ID,
		OrgID:       image.OrgID,
		Name:        image.Name,
		Description: image.Description,
		Author:      image.Author,
		Registry:    image.Registry,
		Namespace:   image.Namespace,
		Repository:  image.Repository,
		Tag:         image.Tag,
		Digest:      image.Digest,
		Size:        image.Size,
		ReadmePath:  image.ReadmePath,
		Stars:       image.Stars,
		Visibility:  image.Visibility,
		Platform:    image.Platform,
		Labels:      make([]string, len(image.Labels)),
		IsStarred:   isStarred,
		CreatedAt:   image.CreatedAt,
		UpdatedAt:   image.UpdatedAt,
	}

	for i, label := range image.Labels {
		response.Labels[i] = label.Name
	}

	return response, nil
}

// CollectImage adds an image to user's collection
func (s *ImageService) CollectImage(userID string, imageID string) error {
	db := database.GetDB()

	// Check if image exists
	var image models.Image
	if err := db.First(&image, "id = ?", imageID).Error; err != nil {
		return errors.New("image not found")
	}

	// Check if already collected
	var collection models.Collection
	if err := db.Where("user_id = ? AND image_id = ?", userID, imageID).First(&collection).Error; err == nil {
		return errors.New("image already collected")
	}

	// Add to collection
	collection = models.Collection{
		UserID:  userID,
		ImageID: imageID,
	}

	if err := db.Create(&collection).Error; err != nil {
		return err
	}

	// Increment stars count
	return db.Model(&image).Update("stars", image.Stars+1).Error
}

// UncollectImage removes an image from user's collection
func (s *ImageService) UncollectImage(userID string, imageID string) error {
	db := database.GetDB()

	// Check if image exists
	var image models.Image
	if err := db.First(&image, "id = ?", imageID).Error; err != nil {
		return errors.New("image not found")
	}

	// Remove from collection
	result := db.Where("user_id = ? AND image_id = ?", userID, imageID).Delete(&models.Collection{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("image not in collection")
	}

	// Decrement stars count
	return db.Model(&image).Update("stars", image.Stars-1).Error
}

// CreateImage creates a new image
func (s *ImageService) CreateImage(ctx context.Context, req *CreateImageRequest, userID string, orgID string) (*ImageResponse, error) {
	db := database.GetDB()

	// 处理 public 组织的情况
	if orgID == "public" {
		orgID = "00000000-0000-0000-0000-000000000000" // 使用特殊的 UUID 表示 public 组织
	}

	// 创建镜像记录
	image := &models.Image{
		Name:        req.Name,
		Description: req.Description,
		Author:      userID,
		Registry:    req.Registry,
		Namespace:   req.Namespace,
		Repository:  req.Repository,
		Tag:         req.Tag,
		Digest:      req.Digest,
		Size:        req.Size,
		OrgID:       orgID,
		Visibility:  req.Visibility,
		Platform:    req.Platform,
	}

	// 如果有 README 文件，上传它
	if req.ReadmeFile != nil {
		readmePath, err := utils.UploadFile(req.ReadmeFile, "readme")
		if err != nil {
			return nil, fmt.Errorf("failed to upload readme file: %v", err)
		}
		image.ReadmePath = readmePath
	}

	// 开始事务
	tx := db.Begin()
	if err := tx.Create(image).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	// 处理标签
	if len(req.Labels) > 0 {
		for _, labelName := range req.Labels {
			var label models.Label
			// 查找或创建标签
			if err := tx.Where("name = ?", labelName).FirstOrCreate(&label, models.Label{Name: labelName}).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create label: %v", err)
			}
			// 关联标签和镜像
			if err := tx.Model(image).Association("Labels").Append(&label); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to associate label: %v", err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// 返回创建的镜像信息
	return s.GetImageByID(ctx, image.ID, userID)
}

// UpdateImage updates an existing image
func (s *ImageService) UpdateImage(ctx context.Context, imageID string, req *UpdateImageRequest, userID string) (*ImageResponse, error) {
	db := database.GetDB()

	// 查找现有镜像
	var image models.Image
	if err := db.First(&image, "id = ? AND author = ?", imageID, userID).Error; err != nil {
		return nil, fmt.Errorf("image not found or not authorized")
	}

	// 开始事务
	tx := db.Begin()

	// 更新基本信息
	if req.Name != "" {
		image.Name = req.Name
	}
	if req.Description != "" {
		image.Description = req.Description
	}
	if req.Registry != "" {
		image.Registry = req.Registry
	}
	if req.Repository != "" {
		image.Repository = req.Repository
	}
	if req.Tag != "" {
		image.Tag = req.Tag
	}

	// 如果有新的 README 文件，上传它并删除旧文件
	if req.ReadmeFile != nil {
		// 删除旧文件
		if image.ReadmePath != "" {
			utils.DeleteFile(image.ReadmePath)
		}
		// 上传新文件
		readmePath, err := utils.UploadFile(req.ReadmeFile, "readme")
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to upload readme file: %v", err)
		}
		image.ReadmePath = readmePath
	}

	// 更新镜像记录
	if err := tx.Save(&image).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update image: %v", err)
	}

	// 如果提供了新的标签列表，更新标签
	if len(req.Labels) > 0 {
		// 清除现有标签
		if err := tx.Model(&image).Association("Labels").Clear(); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to clear labels: %v", err)
		}

		// 添加新标签
		for _, labelName := range req.Labels {
			var label models.Label
			// 查找或创建标签
			if err := tx.Where("name = ?", labelName).FirstOrCreate(&label, models.Label{Name: labelName}).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create label: %v", err)
			}
			// 关联标签和镜像
			if err := tx.Model(&image).Association("Labels").Append(&label); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to associate label: %v", err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// 返回更新后的镜像信息
	return s.GetImageByID(ctx, image.ID, userID)
}

// DeleteImage deletes an image
func (s *ImageService) DeleteImage(ctx context.Context, imageID string, userID string) error {
	db := database.GetDB()

	// 查找镜像
	var image models.Image
	if err := db.First(&image, "id = ? AND author = ?", imageID, userID).Error; err != nil {
		return fmt.Errorf("image not found or not authorized")
	}

	// 开始事务
	tx := db.Begin()

	// 删除相关的收藏记录
	if err := tx.Where("image_id = ?", imageID).Delete(&models.Collection{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete collections: %v", err)
	}

	// 清除标签关联
	if err := tx.Model(&image).Association("Labels").Clear(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear labels: %v", err)
	}

	// 删除 README 文件
	if image.ReadmePath != "" {
		utils.DeleteFile(image.ReadmePath)
	}

	// 删除镜像记录
	if err := tx.Delete(&image).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete image: %v", err)
	}

	// 提交事务
	return tx.Commit().Error
}

// ListFavorites retrieves a list of user's favorite images
func (s *ImageService) ListFavorites(ctx context.Context, req *ImageListRequest, userID string) ([]ImageResponse, int64, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := database.GetDB()

	// 构建查询
	query := db.Model(&models.Image{}).
		Joins("JOIN collections ON collections.image_id = images.id").
		Where("collections.user_id = ?", userID)

	// 应用搜索过滤
	if req.Search != "" {
		query = query.Where("images.name ILIKE ? OR images.description ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 应用标签过滤
	if len(req.Labels) > 0 {
		query = query.Joins("JOIN image_labels ON images.id = image_labels.image_id").
			Where("image_labels.label IN ?", req.Labels).
			Group("images.id").
			Having("COUNT(DISTINCT image_labels.label) = ?", len(req.Labels))
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用排序
	switch req.Sort {
	case "stars":
		query = query.Order("stars DESC")
	case "created_at":
		query = query.Order("created_at DESC")
	case "updated_at":
		query = query.Order("updated_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 应用分页
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	// 执行查询
	var images []models.Image
	if err := query.Preload("Labels").Find(&images).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	response := make([]ImageResponse, len(images))
	for i, img := range images {
		response[i] = ImageResponse{
			ID:          img.ID,
			OrgID:       img.OrgID,
			Name:        img.Name,
			Description: img.Description,
			Author:      img.Author,
			Registry:    img.Registry,
			Namespace:   img.Namespace,
			Repository:  img.Repository,
			Tag:         img.Tag,
			Digest:      img.Digest,
			Size:        img.Size,
			ReadmePath:  img.ReadmePath,
			Stars:       img.Stars,
			Visibility:  img.Visibility,
			Platform:    img.Platform,
			Labels:      make([]string, len(img.Labels)),
			IsStarred:   true, // 这是收藏列表，所以一定是已收藏的
			CreatedAt:   img.CreatedAt,
			UpdatedAt:   img.UpdatedAt,
		}

		for j, label := range img.Labels {
			response[i].Labels[j] = label.Name
		}
	}

	return response, total, nil
} 
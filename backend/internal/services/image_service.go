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
	Page     int      `form:"page" binding:"required,min=1"`
	PageSize int      `form:"size" binding:"required,min=1,max=100"`
	Search   string   `form:"search"`
	Tags     []string `form:"tags"`
	Sort     string   `form:"sort" binding:"oneof=stars created_at updated_at ''"`
}

type ImageResponse struct {
	ID          string             `json:"id"`
	OrgID       string             `json:"org_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Author      string             `json:"author"`
	ReadmePath  string             `json:"readme_path"`
	Stars       int                `json:"stars"`
	Tags        []string           `json:"tags"`
	Providers   []ProviderResponse `json:"providers"`
	IsStarred   bool              `json:"is_starred"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type ProviderResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	APIURL string `json:"api_url"`
}

type CreateImageRequest struct {
	Name        string                `form:"name" binding:"required"`
	Description string                `form:"description"`
	Registry    string                `form:"registry" binding:"required"`
	Repository  string                `form:"repository" binding:"required"`
	Tag         string                `form:"tag" binding:"required"`
	ReadmeFile  *multipart.FileHeader `form:"readme_file"`
	Tags        []string              `form:"tags"`
}

type UpdateImageRequest struct {
	Name        string                `form:"name"`
	Description string                `form:"description"`
	Registry    string                `form:"registry"`
	Repository  string                `form:"repository"`
	Tag         string                `form:"tag"`
	ReadmeFile  *multipart.FileHeader `form:"readme_file"`
	Tags        []string              `form:"tags"`
}

// NewImageService creates a new ImageService
func NewImageService() *ImageService {
	return &ImageService{}
}

// ListImages retrieves a list of images with pagination and filtering
func (s *ImageService) ListImages(ctx context.Context, req *ImageListRequest, userID string) ([]ImageResponse, int64, error) {
	db := database.GetDB()
	rdb := database.GetRedis()

	// Try to get from cache
	cacheKey := fmt.Sprintf("images:%d:%d:%s:%v:%s", req.Page, req.PageSize, req.Search, req.Tags, req.Sort)
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

	// Apply tags filter
	if len(req.Tags) > 0 {
		query = query.Joins("JOIN image_tags ON images.id = image_tags.image_id").
			Where("image_tags.tag IN ?", req.Tags).
			Group("images.id").
			Having("COUNT(DISTINCT image_tags.tag) = ?", len(req.Tags))
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
	if err := query.Preload("Tags").Preload("Providers").Find(&images).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response
	response := make([]ImageResponse, len(images))
	for i, img := range images {
		// Check if user has starred the image
		var isStarred bool
		if userID != "" {
			var collection models.Collection
			isStarred = db.Where("user_id = ? AND image_id = ?", userID, img.ID).First(&collection).Error == nil
		}

		response[i] = ImageResponse{
			ID:          img.ID,
			OrgID:       img.OrgID,
			Name:        img.Name,
			Description: img.Description,
			Author:      img.Author,
			ReadmePath:  img.ReadmePath,
			Stars:       img.Stars,
			Tags:        make([]string, len(img.Tags)),
			Providers:   make([]ProviderResponse, len(img.Providers)),
			IsStarred:   isStarred,
			CreatedAt:   img.CreatedAt,
			UpdatedAt:   img.UpdatedAt,
		}

		for j, tag := range img.Tags {
			response[i].Tags[j] = tag.Name
		}

		for j, provider := range img.Providers {
			response[i].Providers[j] = ProviderResponse{
				ID:     provider.ID,
				Name:   provider.Name,
				APIURL: provider.APIURL,
			}
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
	if err := db.Preload("Tags").Preload("Providers").First(&image, "id = ?", id).Error; err != nil {
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
		ReadmePath:  image.ReadmePath,
		Stars:       image.Stars,
		Tags:        make([]string, len(image.Tags)),
		Providers:   make([]ProviderResponse, len(image.Providers)),
		IsStarred:   isStarred,
		CreatedAt:   image.CreatedAt,
		UpdatedAt:   image.UpdatedAt,
	}

	for i, tag := range image.Tags {
		response.Tags[i] = tag.Name
	}

	for i, provider := range image.Providers {
		response.Providers[i] = ProviderResponse{
			ID:     provider.ID,
			Name:   provider.Name,
			APIURL: provider.APIURL,
		}
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
func (s *ImageService) CreateImage(ctx context.Context, req *CreateImageRequest, userID string) (*ImageResponse, error) {
	db := database.GetDB()

	// 创建镜像记录
	image := &models.Image{
		Name:        req.Name,
		Description: req.Description,
		Author:      userID,
		Registry:    req.Registry,
		Repository:  req.Repository,
		Tag:         req.Tag,
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
	if len(req.Tags) > 0 {
		for _, tagName := range req.Tags {
			var tag models.Tag
			// 查找或创建标签
			if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create tag: %v", err)
			}
			// 关联标签和镜像
			if err := tx.Model(image).Association("Tags").Append(&tag); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to associate tag: %v", err)
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
	if len(req.Tags) > 0 {
		// 清除现有标签
		if err := tx.Model(&image).Association("Tags").Clear(); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to clear tags: %v", err)
		}

		// 添加新标签
		for _, tagName := range req.Tags {
			var tag models.Tag
			// 查找或创建标签
			if err := tx.Where("name = ?", tagName).FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create tag: %v", err)
			}
			// 关联标签和镜像
			if err := tx.Model(&image).Association("Tags").Append(&tag); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to associate tag: %v", err)
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
	if err := tx.Model(&image).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to clear tags: %v", err)
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
	if len(req.Tags) > 0 {
		query = query.Joins("JOIN image_tags ON images.id = image_tags.image_id").
			Where("image_tags.tag IN ?", req.Tags).
			Group("images.id").
			Having("COUNT(DISTINCT image_tags.tag) = ?", len(req.Tags))
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
	if err := query.Preload("Tags").Preload("Providers").Find(&images).Error; err != nil {
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
			ReadmePath:  img.ReadmePath,
			Stars:       img.Stars,
			Tags:        make([]string, len(img.Tags)),
			Providers:   make([]ProviderResponse, len(img.Providers)),
			IsStarred:   true, // 这是收藏列表，所以一定是已收藏的
			UpdatedAt:   img.UpdatedAt,
		}

		for j, tag := range img.Tags {
			response[i].Tags[j] = tag.Name
		}

		for j, provider := range img.Providers {
			response[i].Providers[j] = ProviderResponse{
				ID:     provider.ID,
				Name:   provider.Name,
				APIURL: provider.APIURL,
			}
		}
	}

	return response, total, nil
} 
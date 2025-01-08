package services

import (
	"errors"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/models"
)

type DeployService struct{}

type DeployRequest struct {
	ImageID string                 `json:"image_id" binding:"required"`
	Params  map[string]interface{} `json:"params"`
}

type DeployResponse struct {
	ImageID string                 `json:"image_id"`
	Params  map[string]interface{} `json:"params"`
}

// NewDeployService creates a new DeployService
func NewDeployService() *DeployService {
	return &DeployService{}
}

// Deploy prepares deployment information for an image
func (s *DeployService) Deploy(req *DeployRequest) (*DeployResponse, error) {
	db := database.GetDB()

	// Verify image exists
	var image models.Image
	if err := db.First(&image, "id = ?", req.ImageID).Error; err != nil {
		return nil, errors.New("image not found")
	}

	return &DeployResponse{
		ImageID: req.ImageID,
		Params:  req.Params,
	}, nil
} 
package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/models"
)

type DeployService struct{}

type DeployRequest struct {
	ImageID    string                 `json:"image_id" binding:"required"`
	ProviderID string                 `json:"provider_id" binding:"required"`
	Params     map[string]interface{} `json:"params"`
}

type DeployResponse struct {
	ProviderName string                 `json:"provider_name"`
	APIURL       string                 `json:"api_url"`
	Params       map[string]interface{} `json:"params"`
}

// NewDeployService creates a new DeployService
func NewDeployService() *DeployService {
	return &DeployService{}
}

// GetDeployInfo retrieves deployment information for an image
func (s *DeployService) GetDeployInfo(imageID string, providerID string) (*DeployResponse, error) {
	db := database.GetDB()

	// Get image and provider information
	var imageProvider models.ImageProvider
	if err := db.Where("image_id = ? AND provider_id = ?", imageID, providerID).First(&imageProvider).Error; err != nil {
		return nil, errors.New("deployment configuration not found")
	}

	var provider models.Provider
	if err := db.First(&provider, "id = ?", providerID).Error; err != nil {
		return nil, errors.New("provider not found")
	}

	// Parse provider params
	var params map[string]interface{}
	if imageProvider.Params != "" {
		if err := json.Unmarshal([]byte(imageProvider.Params), &params); err != nil {
			return nil, fmt.Errorf("failed to parse provider params: %v", err)
		}
	}

	return &DeployResponse{
		ProviderName: provider.Name,
		APIURL:       provider.APIURL,
		Params:       params,
	}, nil
}

// Deploy prepares deployment information for an image
func (s *DeployService) Deploy(req *DeployRequest) (*DeployResponse, error) {
	db := database.GetDB()

	// Verify image exists
	var image models.Image
	if err := db.First(&image, "id = ?", req.ImageID).Error; err != nil {
		return nil, errors.New("image not found")
	}

	// Verify provider exists
	var provider models.Provider
	if err := db.First(&provider, "id = ?", req.ProviderID).Error; err != nil {
		return nil, errors.New("provider not found")
	}

	// Get or create image provider configuration
	var imageProvider models.ImageProvider
	result := db.Where("image_id = ? AND provider_id = ?", req.ImageID, req.ProviderID).First(&imageProvider)
	
	if result.Error != nil {
		// Create new configuration
		paramsJSON, err := json.Marshal(req.Params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %v", err)
		}

		imageProvider = models.ImageProvider{
			ImageID:    req.ImageID,
			ProviderID: req.ProviderID,
			Params:     string(paramsJSON),
		}

		if err := db.Create(&imageProvider).Error; err != nil {
			return nil, fmt.Errorf("failed to create deployment configuration: %v", err)
		}
	} else {
		// Update existing configuration
		paramsJSON, err := json.Marshal(req.Params)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal params: %v", err)
		}

		imageProvider.Params = string(paramsJSON)
		if err := db.Save(&imageProvider).Error; err != nil {
			return nil, fmt.Errorf("failed to update deployment configuration: %v", err)
		}
	}

	return &DeployResponse{
		ProviderName: provider.Name,
		APIURL:       provider.APIURL,
		Params:       req.Params,
	}, nil
} 
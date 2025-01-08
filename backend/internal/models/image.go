package models

import (
	"time"
)

type Image struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrgID       string    `json:"org_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Author      string    `json:"author" gorm:"type:uuid;not null"`
	ReadmePath  string    `json:"readme_path"`
	Stars       int       `json:"stars" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []Tag     `json:"tags" gorm:"many2many:image_tags;"`
	Providers   []Provider `json:"providers" gorm:"many2many:image_providers;"`
}

type Tag struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Provider struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	APIURL    string    `json:"api_url" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ImageProvider struct {
	ImageID    string         `json:"image_id" gorm:"type:uuid;not null"`
	ProviderID string         `json:"provider_id" gorm:"type:uuid;not null"`
	Params     string         `json:"params" gorm:"type:jsonb"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type Collection struct {
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	ImageID   string    `json:"image_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName - Set the table names for the models
func (Image) TableName() string {
	return "images"
}

func (Tag) TableName() string {
	return "tags"
}

func (Provider) TableName() string {
	return "providers"
}

func (ImageProvider) TableName() string {
	return "image_providers"
}

func (Collection) TableName() string {
	return "collections"
} 
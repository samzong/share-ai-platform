package models

import (
	"time"
)

// Image 表示一个容器镜像
type Image struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"` // 镜像唯一标识符
	OrgID       string    `json:"org_id" gorm:"type:uuid;not null"`                         // 组织ID
	Name        string    `json:"name" gorm:"not null"`                                     // 镜像显示名称
	Description string    `json:"description"`                                              // 镜像描述
	Author      string    `json:"author" gorm:"type:uuid;not null"`                        // 创建者ID
	Registry    string    `json:"registry" gorm:"not null"`                                // 镜像仓库服务器（例如：docker.io）
	Namespace   string    `json:"namespace" gorm:"not null"`                               // 命名空间/组织（例如：library）
	Repository  string    `json:"repository" gorm:"not null"`                              // 镜像名称（例如：nginx）
	Tag         string    `json:"tag" gorm:"not null"`                                     // 版本标签（例如：latest）
	Digest      string    `json:"digest" gorm:"not null"`                                  // 镜像内容哈希值
	Size        int64     `json:"size" gorm:"default:0"`                                   // 镜像大小（字节）
	ReadmePath  string    `json:"readme_path"`                                             // README文件路径
	Stars       int       `json:"stars" gorm:"default:0"`                                  // 收藏数（通过 Collection 表关联计算）
	Visibility  string    `json:"visibility" gorm:"type:varchar(10);not null;default:'public'"` // 可见性：public/private
	Platform    string    `json:"platform" gorm:"not null"`                                // 平台架构（例如：linux/amd64）
	Labels      []Label   `json:"labels" gorm:"many2many:image_labels;constraint:OnDelete:CASCADE;"` // 标签列表，用于分类和搜索
	CreatedAt   time.Time `json:"created_at"`                                             // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                                             // 更新时间
}

// Label 表示镜像的分类标签
type Label struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"` // 标签唯一标识符
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`                         // 标签名称
	CreatedAt time.Time `json:"created_at"`                                               // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                                               // 更新时间
}

// Collection 表示用户收藏的镜像
type Collection struct {
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`  // 用户ID
	ImageID   string    `json:"image_id" gorm:"type:uuid;not null"` // 镜像ID
	CreatedAt time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
}

// TableName - Set the table names for the models
func (Image) TableName() string {
	return "images"
}

func (Label) TableName() string {
	return "labels"
}

func (Collection) TableName() string {
	return "collections"
}
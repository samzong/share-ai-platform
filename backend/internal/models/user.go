package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User 表示系统用户
type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"type:varchar(100);not null"` // "-" means this field will not be included in JSON
	Nickname  string    `json:"nickname" gorm:"type:varchar(50)"`    // 昵称
	Avatar    string    `json:"avatar" gorm:"type:varchar(255)"`     // 头像URL
	Role      Role      `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// BeforeCreate - GORM hook that runs before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	
	// Set default role if not specified
	if u.Role == "" {
		u.Role = RoleUser
	}

	// Set default nickname if not specified
	if u.Nickname == "" {
		u.Nickname = u.Username
	}
	
	return nil
}

// ComparePassword - Compare the provided password with the hashed password
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// IsAdmin - Check if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// TableName - Set the table name for the User model
func (User) TableName() string {
	return "users"
}

// IsValidRole checks if a role is valid
func IsValidRole(role Role) bool {
	return role == RoleUser || role == RoleAdmin
} 
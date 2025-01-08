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

type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"` // "-" means this field will not be included in JSON
	Role      Role      `json:"role" gorm:"type:varchar(20);default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
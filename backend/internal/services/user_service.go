package services

import (
	"errors"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/middleware"
	"github.com/samzong/share-ai-platform/internal/models"
)

type UserService struct{}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token,omitempty"`
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{}
}

// Register creates a new user
func (s *UserService) Register(req *RegisterRequest) (*UserResponse, error) {
	db := database.GetDB()

	// Check if username already exists
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Create new user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

// Login authenticates a user
func (s *UserService) Login(req *LoginRequest) (*UserResponse, error) {
	db := database.GetDB()

	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := user.ComparePassword(req.Password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID string) (*UserResponse, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID string, username, email string) (*UserResponse, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	// Update fields if provided
	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
} 
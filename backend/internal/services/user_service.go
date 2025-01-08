package services

import (
	"errors"
	"regexp"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/middleware"
	"github.com/samzong/share-ai-platform/internal/models"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
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
	ID       string      `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Role     models.Role `json:"role"`
	Token    string      `json:"token,omitempty"`
}

// ListUsersRequest represents the request parameters for listing users
type ListUsersRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Search   string `form:"search"`
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64         `json:"total"`
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{}
}

// Register creates a new user
func (s *UserService) Register(req *RegisterRequest) (*UserResponse, error) {
	// Validate email format
	if !emailRegex.MatchString(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate password length
	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

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
		Role:     models.RoleUser, // Default role
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
		Role:     user.Role,
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
		Role:     user.Role,
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
		Role:     user.Role,
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
		// Check if new username is already taken
		var existingUser models.User
		if err := db.Where("username = ? AND id != ?", username, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("username already exists")
		}
		user.Username = username
	}
	
	if email != "" {
		// Check if new email is already taken
		var existingUser models.User
		if err := db.Where("email = ? AND id != ?", email, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already exists")
		}
		user.Email = email
	}

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

// UpdateUserRole updates a user's role (admin only)
func (s *UserService) UpdateUserRole(adminID, userID string, newRole models.Role) error {
	db := database.GetDB()

	// Check if the admin user exists and is actually an admin
	var admin models.User
	if err := db.First(&admin, "id = ?", adminID).Error; err != nil {
		return errors.New("admin not found")
	}

	// Ensure the user is an admin
	if admin.Role != models.RoleAdmin {
		return errors.New("insufficient permissions")
	}

	// Don't allow changing own role
	if adminID == userID {
		return errors.New("cannot change own role")
	}

	// Update the user's role
	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Validate the new role
	if !models.IsValidRole(newRole) {
		return errors.New("invalid role")
	}

	// Update the role
	user.Role = newRole
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// ListUsers retrieves a paginated list of users with optional search
func (s *UserService) ListUsers(req *ListUsersRequest) (*ListUsersResponse, error) {
	db := database.GetDB()
	var users []models.User
	var total int64

	query := db.Model(&models.User{})

	// Apply search if provided
	if req.Search != "" {
		searchQuery := "%" + req.Search + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ?", searchQuery, searchQuery)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}
	}

	return &ListUsersResponse{
		Users: userResponses,
		Total: total,
	}, nil
} 
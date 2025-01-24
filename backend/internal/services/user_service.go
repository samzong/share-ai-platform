package services

import (
	"context"
	"errors"
	"mime/multipart"
	"regexp"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/middleware"
	"github.com/samzong/share-ai-platform/internal/models"
	"github.com/samzong/share-ai-platform/internal/utils"
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

type UpdateProfileRequest struct {
	Nickname string                `form:"nickname"`
	Avatar   *multipart.FileHeader `form:"avatar"`
}

type UserResponse struct {
	ID       string      `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Nickname string      `json:"nickname"`
	Avatar   string      `json:"avatar"`
	Role     models.Role `json:"role"`
	Token    string      `json:"token,omitempty"`
}

type ListUsersRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
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
		Nickname: user.Nickname,
		Avatar:   utils.GetFileURL(user.Avatar),
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
		Nickname: user.Nickname,
		Avatar:   utils.GetFileURL(user.Avatar),
		Role:     user.Role,
		Token:    token,
	}, nil
}

// Logout invalidates a user's token
func (s *UserService) Logout(ctx context.Context, userID string) error {
	// 将token加入黑名单
	token := middleware.GetTokenFromContext(ctx)
	if token != "" {
		return database.GetRedis().Set(ctx, "blacklist:"+token, userID, middleware.TokenExpiration).Err()
	}
	return nil
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
		Nickname: user.Nickname,
		Avatar:   utils.GetFileURL(user.Avatar),
		Role:     user.Role,
	}, nil
}

// UpdateProfile updates user's profile information
func (s *UserService) UpdateProfile(userID string, req *UpdateProfileRequest) (*UserResponse, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Update nickname if provided
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	// Handle avatar upload if provided
	if req.Avatar != nil {
		// Delete old avatar if exists
		if user.Avatar != "" {
			utils.DeleteFile(user.Avatar)
		}

		// Upload new avatar
		avatarPath, err := utils.UploadFile(req.Avatar, "avatars")
		if err != nil {
			return nil, err
		}
		user.Avatar = avatarPath
	}

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   utils.GetFileURL(user.Avatar),
		Role:     user.Role,
	}, nil
}

// UpdateUser updates user's username and email
func (s *UserService) UpdateUser(userID string, username string, email string) (*UserResponse, error) {
	if !emailRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	db := database.GetDB()
	user := &models.User{}
	if err := db.First(user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	user.Username = username
	user.Email = email

	if err := db.Save(user).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Role:     user.Role,
	}, nil
}

// UpdateUserRole updates a user's role (admin only)
func (s *UserService) UpdateUserRole(adminID string, userID string, role models.Role) error {
	db := database.GetDB()
	// 验证管理员权限
	admin := &models.User{}
	if err := db.First(admin, "id = ?", adminID).Error; err != nil {
		return err
	}
	if admin.Role != models.RoleAdmin {
		return errors.New("permission denied: requires admin role")
	}

	// 更新用户角色
	user := &models.User{}
	if err := db.First(user, "id = ?", userID).Error; err != nil {
		return err
	}

	user.Role = role
	if err := db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

// ListUsers returns a paginated list of users (admin only)
func (s *UserService) ListUsers(req *ListUsersRequest) (*struct {
	Total int64          `json:"total"`
	Users []UserResponse `json:"users"`
}, error) {
	db := database.GetDB()
	var users []models.User
	var total int64

	offset := (req.Page - 1) * req.PageSize

	if err := db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Role:     user.Role,
		}
	}

	return &struct {
		Total int64          `json:"total"`
		Users []UserResponse `json:"users"`
	}{
		Total: total,
		Users: userResponses,
	}, nil
}

package services

import (
	"testing"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*UserService, func()) {
	db := database.SetupTestDB()

	// Auto migrate the schema
	err := db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	// Clear all records
	err = db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error
	assert.NoError(t, err)

	service := NewUserService()

	return service, func() {
		database.CleanupTestDB(db)
	}
}

func TestUserService_Register(t *testing.T) {
	service, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *RegisterRequest
		wantErr bool
	}{
		{
			name: "valid registration",
			req: &RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			req: &RegisterRequest{
				Username: "testuser2",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "short password",
			req: &RegisterRequest{
				Username: "testuser3",
				Email:    "test3@example.com",
				Password: "123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Register(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Username, resp.Username)
				assert.Equal(t, tt.req.Email, resp.Email)
				assert.Equal(t, models.RoleUser, resp.Role)
				assert.NotEmpty(t, resp.Token)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	service, cleanup := setupTest(t)
	defer cleanup()

	// First register a user
	registerReq := &RegisterRequest{
		Username: "logintest",
		Email:    "login@example.com",
		Password: "password123",
	}
	_, err := service.Register(registerReq)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		req     *LoginRequest
		wantErr bool
	}{
		{
			name: "valid login",
			req: &LoginRequest{
				Username: "logintest",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "invalid password",
			req: &LoginRequest{
				Username: "logintest",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "non-existent user",
			req: &LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Login(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Username, resp.Username)
				assert.NotEmpty(t, resp.Token)
			}
		})
	}
}

func TestUserService_UpdateUserRole(t *testing.T) {
	service, cleanup := setupTest(t)
	defer cleanup()

	db := database.GetDB()

	// Register an admin user
	adminReq := &RegisterRequest{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin123",
	}
	adminResp, err := service.Register(adminReq)
	assert.NoError(t, err)

	// Register a normal user
	userReq := &RegisterRequest{
		Username: "user",
		Email:    "user@example.com",
		Password: "user123",
	}
	userResp, err := service.Register(userReq)
	assert.NoError(t, err)

	// Set admin role directly in the database
	err = db.Model(&models.User{}).Where("id = ?", adminResp.ID).Update("role", models.RoleAdmin).Error
	assert.NoError(t, err)

	// Verify admin role was set correctly
	var adminUser models.User
	err = db.First(&adminUser, "id = ?", adminResp.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, models.RoleAdmin, adminUser.Role)

	tests := []struct {
		name    string
		adminID string
		userID  string
		newRole models.Role
		wantErr bool
	}{
		{
			name:    "admin can change user role",
			adminID: adminResp.ID,
			userID:  userResp.ID,
			newRole: models.RoleAdmin,
			wantErr: false,
		},
		{
			name:    "non-admin cannot change role",
			adminID: userResp.ID, // Using normal user as admin
			userID:  adminResp.ID,
			newRole: models.RoleUser,
			wantErr: true,
		},
		{
			name:    "admin cannot change own role",
			adminID: adminResp.ID,
			userID:  adminResp.ID,
			newRole: models.RoleUser,
			wantErr: true,
		},
		{
			name:    "invalid role",
			adminID: adminResp.ID,
			userID:  userResp.ID,
			newRole: "invalid_role",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateUserRole(tt.adminID, tt.userID, tt.newRole)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the role was actually changed
				var user models.User
				err := db.First(&user, "id = ?", tt.userID).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.newRole, user.Role)
			}
		})
	}
}

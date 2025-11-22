package services

import (
	"context"
	"time"

	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService defines the interface for user operations
type UserService interface {
	Create(ctx context.Context, email, password string, role models.UserRole, firstName, lastName string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) (*models.User, error)
	Delete(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id string, oldPassword, newPassword string) error
	ValidatePassword(ctx context.Context, email, password string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, id string) error
	HasPermission(ctx context.Context, userID string, permission string) (bool, error)
}

type userService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) Create(ctx context.Context, email, password string, role models.UserRole, firstName, lastName string) (*models.User, error) {
	// Check if user already exists
	var existing models.User
	if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.Conflict("User with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternal, "Failed to hash password")
	}

	user := &models.User{
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		FirstName: firstName,
		LastName:  lastName,
		IsActive:  true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.DatabaseError("creating user", err)
	}

	// Don't return password hash
	user.Password = ""
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Teacher").Preload("Student").First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("User", id)
		}
		return nil, errors.DatabaseError("finding user", err)
	}
	user.Password = "" // Never expose password
	return &user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("User not found")
		}
		return nil, errors.DatabaseError("finding user by email", err)
	}
	return &user, nil
}

func (s *userService) Update(ctx context.Context, id string, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("User", id)
		}
		return nil, errors.DatabaseError("finding user", err)
	}

	// Don't allow password updates through this method
	delete(updates, "password")

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, errors.DatabaseError("updating user", err)
	}

	user.Password = ""
	return &user, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return errors.DatabaseError("deleting user", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("User", id)
	}
	return nil
}

func (s *userService) ChangePassword(ctx context.Context, id string, oldPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFoundWithID("User", id)
		}
		return errors.DatabaseError("finding user", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.Unauthorized("Invalid current password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(errors.ErrCodeInternal, "Failed to hash password")
	}

	if err := s.db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return errors.DatabaseError("updating password", err)
	}

	return nil
}

func (s *userService) ValidatePassword(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.Unauthorized("User account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	user.Password = ""
	return user, nil
}

func (s *userService) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}

func (s *userService) HasPermission(ctx context.Context, userID string, permission string) (bool, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return false, err
	}

	// Admin has all permissions
	if user.Role == models.RoleAdmin {
		return true, nil
	}

	// Check role-specific permissions
	var count int64
	err := s.db.Table("role_permissions").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role = ? AND CONCAT(permissions.resource, ':', permissions.action) = ?", user.Role, permission).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

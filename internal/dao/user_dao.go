package dao

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/models"
	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/log"
)

// UserDAO handles user data access operations
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO creates a new UserDAO instance
func NewUserDAO() *UserDAO {
	return &UserDAO{
		db: DB,
	}
}

// Create creates a new user in the database
func (dao *UserDAO) Create(user *table.User) error {
	if err := dao.db.Create(user).Error; err != nil {
		log.Error(err).Str("username", user.Username).Msg("Failed to create user")
		return fmt.Errorf("failed to create user: %w", err)
	}
	log.Info().Str("user_id", user.ID).Str("username", user.Username).Msg("User created successfully")
	return nil
}

// FindByID finds a user by ID
func (dao *UserDAO) FindByID(id string) (*table.User, error) {
	var user table.User
	if err := dao.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		log.Error(err).Str("user_id", id).Msg("Failed to find user by ID")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (dao *UserDAO) FindByUsername(username string) (*table.User, error) {
	var user table.User
	if err := dao.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		log.Error(err).Str("username", username).Msg("Failed to find user by username")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (dao *UserDAO) FindByEmail(email string) (*table.User, error) {
	var user table.User
	if err := dao.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		log.Error(err).Str("email", email).Msg("Failed to find user by email")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// FindByUsernameOrEmail finds a user by either username or email
func (dao *UserDAO) FindByUsernameOrEmail(usernameOrEmail string) (*table.User, error) {
	var user table.User
	err := dao.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		log.Error(err).Str("identifier", usernameOrEmail).Msg("Failed to find user by username or email")
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

// Update updates a user in the database
func (dao *UserDAO) Update(user *table.User) error {
	if err := dao.db.Save(user).Error; err != nil {
		log.Error(err).Str("user_id", user.ID).Msg("Failed to update user")
		return fmt.Errorf("failed to update user: %w", err)
	}
	log.Info().Str("user_id", user.ID).Msg("User updated successfully")
	return nil
}

// UpdateProfile updates user profile information
func (dao *UserDAO) UpdateProfile(userID string, updateData *dto.UserUpdateRequest) (*table.User, error) {
	var user table.User
	if err := dao.db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Update only the provided fields
	updateMap := make(map[string]interface{})
	if updateData.FullName != "" {
		updateMap["full_name"] = updateData.FullName
	}
	if updateData.AvatarURL != "" {
		updateMap["avatar_url"] = updateData.AvatarURL
	}
	if updateData.Bio != "" {
		updateMap["bio"] = updateData.Bio
	}

	if len(updateMap) == 0 {
		return &user, nil // No updates needed
	}

	if err := dao.db.Model(&user).Updates(updateMap).Error; err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to update user profile")
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	log.Info().Str("user_id", userID).Msg("User profile updated successfully")
	return &user, nil
}

// UpdatePassword updates a user's password
func (dao *UserDAO) UpdatePassword(userID string, newPassword string) error {
	var user table.User
	if err := dao.db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	businessUser := models.NewUserBusiness(&user)
	if err := businessUser.SetPassword(newPassword); err != nil {
		return fmt.Errorf("failed to set new password: %w", err)
	}

	if err := dao.db.Model(&user).Update("password_hash", user.PasswordHash).Error; err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to update user password")
		return fmt.Errorf("failed to update user password: %w", err)
	}

	log.Info().Str("user_id", userID).Msg("User password updated successfully")
	return nil
}

// Delete deletes a user from the database (soft delete)
func (dao *UserDAO) Delete(userID string) error {
	if err := dao.db.Delete(&table.User{}, "id = ?", userID).Error; err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}
	log.Info().Str("user_id", userID).Msg("User deleted successfully")
	return nil
}

// List returns a paginated list of users
func (dao *UserDAO) List(page, pageSize int) ([]table.User, int64, error) {
	var users []table.User
	var total int64

	// Get total count
	if err := dao.db.Model(&table.User{}).Count(&total).Error; err != nil {
		log.Error(err).Msg("Failed to count users")
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := dao.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		log.Error(err).Msg("Failed to list users")
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// FindActiveUsers finds all active users
func (dao *UserDAO) FindActiveUsers() ([]table.User, error) {
	var users []table.User
	if err := dao.db.Where("is_active = ?", true).Find(&users).Error; err != nil {
		log.Error(err).Msg("Failed to find active users")
		return nil, fmt.Errorf("failed to find active users: %w", err)
	}
	return users, nil
}

// FindAdmins finds all admin users
func (dao *UserDAO) FindAdmins() ([]table.User, error) {
	var users []table.User
	if err := dao.db.Where("role = ?", "admin").Find(&users).Error; err != nil {
		log.Error(err).Msg("Failed to find admin users")
		return nil, fmt.Errorf("failed to find admin users: %w", err)
	}
	return users, nil
}

// Count returns the total number of users
func (dao *UserDAO) Count() (int64, error) {
	var count int64
	if err := dao.db.Model(&table.User{}).Count(&count).Error; err != nil {
		log.Error(err).Msg("Failed to count users")
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// ExistsByUsername checks if a username already exists
func (dao *UserDAO) ExistsByUsername(username string) (bool, error) {
	var count int64
	if err := dao.db.Model(&table.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		log.Error(err).Str("username", username).Msg("Failed to check username existence")
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail checks if an email already exists
func (dao *UserDAO) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := dao.db.Model(&table.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		log.Error(err).Str("email", email).Msg("Failed to check email existence")
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}
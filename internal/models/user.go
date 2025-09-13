package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/table"
)

// UserBusiness represents user business logic methods
// This extends the table.User with business logic methods
type UserBusiness struct {
	*table.User
}

// NewUserBusiness creates a new user business instance
func NewUserBusiness(user *table.User) *UserBusiness {
	return &UserBusiness{User: user}
}

// ToResponse converts User to UserResponse (hides sensitive data)
func (u *UserBusiness) ToResponse() dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		AvatarURL: u.AvatarURL,
		Bio:       u.Bio,
		Role:      u.Role,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// SetPassword hashes and sets the user's password
func (u *UserBusiness) SetPassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	// Generate bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword verifies the user's password
func (u *UserBusiness) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// UpdateLastLogin updates the user's last login time
func (u *UserBusiness) UpdateLastLogin() error {
	now := time.Now().UnixMilli()
	u.LastLogin = &now
	return nil
}

// IsAdmin checks if the user has admin role
func (u *UserBusiness) IsAdmin() bool {
	return u.Role == "admin"
}

// IsActiveUser checks if the user account is active
func (u *UserBusiness) IsActiveUser() bool {
	return u.IsActive
}
package service

import (
	"fmt"

	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/models"
	"github.com/sanmu2018/word-hero/internal/table"
)

// UserService handles user business logic
type UserService struct {
	userDAO *dao.UserDAO
}

// NewUserService creates a new UserService instance
func NewUserService(userDAO *dao.UserDAO) *UserService {
	return &UserService{
		userDAO: userDAO,
	}
}

// GetUserProfile returns a user's profile
func (s *UserService) GetUserProfile(userID string) (*dto.UserResponse, error) {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	businessUser := models.NewUserBusiness(user)
	profile := businessUser.ToResponse()
	return &profile, nil
}

// UpdateUserProfile updates a user's profile
func (s *UserService) UpdateUserProfile(userID string, req *dto.UserUpdateRequest) (*dto.UserResponse, error) {
	user, err := s.userDAO.UpdateProfile(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	businessUser := models.NewUserBusiness(user)
	profile := businessUser.ToResponse()
	return &profile, nil
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(userID string) (*table.User, error) {
	return s.userDAO.FindByID(userID)
}

// GetUserByUsername returns a user by username
func (s *UserService) GetUserByUsername(username string) (*table.User, error) {
	return s.userDAO.FindByUsername(username)
}

// GetUserByEmail returns a user by email
func (s *UserService) GetUserByEmail(email string) (*table.User, error) {
	return s.userDAO.FindByEmail(email)
}

// ListUsers returns a paginated list of users
func (s *UserService) ListUsers(page, pageSize int) ([]dto.UserResponse, int64, error) {
	users, total, err := s.userDAO.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to response format
	userResponses := make([]dto.UserResponse, len(users))
	for i := range users {
		businessUser := models.NewUserBusiness(&users[i])
		userResponses[i] = businessUser.ToResponse()
	}

	return userResponses, total, nil
}

// DeleteUser deletes a user (soft delete)
func (s *UserService) DeleteUser(userID string) error {
	return s.userDAO.Delete(userID)
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(userID string) error {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.IsActive = true
	return s.userDAO.Update(user)
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(userID string) error {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.IsActive = false
	return s.userDAO.Update(user)
}

// PromoteToAdmin promotes a user to admin role
func (s *UserService) PromoteToAdmin(userID string) error {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.Role = "admin"
	return s.userDAO.Update(user)
}

// DemoteFromUser demotes an admin to user role
func (s *UserService) DemoteFromUser(userID string) error {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.Role = "user"
	return s.userDAO.Update(user)
}

// GetActiveUsers returns all active users
func (s *UserService) GetActiveUsers() ([]dto.UserResponse, error) {
	users, err := s.userDAO.FindActiveUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	// Convert to response format
	userResponses := make([]dto.UserResponse, len(users))
	for i := range users {
		businessUser := models.NewUserBusiness(&users[i])
		userResponses[i] = businessUser.ToResponse()
	}

	return userResponses, nil
}

// GetAdminUsers returns all admin users
func (s *UserService) GetAdminUsers() ([]dto.UserResponse, error) {
	users, err := s.userDAO.FindAdmins()
	if err != nil {
		return nil, fmt.Errorf("failed to get admin users: %w", err)
	}

	// Convert to response format
	userResponses := make([]dto.UserResponse, len(users))
	for i := range users {
		businessUser := models.NewUserBusiness(&users[i])
		userResponses[i] = businessUser.ToResponse()
	}

	return userResponses, nil
}

// GetUserCount returns the total number of users
func (s *UserService) GetUserCount() (int64, error) {
	return s.userDAO.Count()
}

// CheckUsernameAvailability checks if a username is available
func (s *UserService) CheckUsernameAvailability(username string) (bool, error) {
	return s.userDAO.ExistsByUsername(username)
}

// CheckEmailAvailability checks if an email is available
func (s *UserService) CheckEmailAvailability(email string) (bool, error) {
	return s.userDAO.ExistsByEmail(email)
}
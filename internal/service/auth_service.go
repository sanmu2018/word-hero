package service

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/models"
	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/internal/utils"
	"github.com/sanmu2018/word-hero/log"
)

// AuthService handles authentication business logic
type AuthService struct {
	userDAO   *dao.UserDAO
	jwtUtils  *utils.JWTUtils
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userDAO *dao.UserDAO, jwtUtils *utils.JWTUtils) *AuthService {
	return &AuthService{
		userDAO:  userDAO,
		jwtUtils: jwtUtils,
	}
}

// Register registers a new user
func (s *AuthService) Register(req *dto.UserRegisterRequest) (*table.User, string, error) {
	// Validate input
	if err := s.validateRegistrationInput(req.Username, req.Email, req.Password); err != nil {
		return nil, "", err
	}

	// Check if username already exists
	exists, err := s.userDAO.ExistsByUsername(req.Username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check username availability: %w", err)
	}
	if exists {
		return nil, "", errors.New("username already exists")
	}

	// Check if email already exists
	exists, err = s.userDAO.ExistsByEmail(req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check email availability: %w", err)
	}
	if exists {
		return nil, "", errors.New("email already exists")
	}

	// Create new user
	user := &table.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     "user",
		IsActive: true,
	}

	// Hash password
	if err := s.hashPassword(user, req.Password); err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Save user to database
	if err := s.userDAO.Create(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtUtils.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	log.Info().
		Str("user_id", user.ID).
		Str("username", user.Username).
		Str("email", user.Email).
		Msg("User registered successfully")

	return user, token, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(req *dto.UserLoginRequest) (*table.User, string, error) {
	// Find user by username or email
	user, err := s.userDAO.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", errors.New("account is deactivated")
	}

	// Verify password
	if !s.verifyPassword(user, req.Password) {
		return nil, "", errors.New("invalid username or password")
	}

	// Update last login time
	businessUser := models.NewUserBusiness(user)
	if err := businessUser.UpdateLastLogin(); err != nil {
		log.Warn().Err(err).Str("user_id", user.ID).Msg("Failed to update last login time")
	}

	// Generate JWT token
	token, err := s.jwtUtils.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	log.Info().
		Str("user_id", user.ID).
		Str("username", user.Username).
		Msg("User logged in successfully")

	return user, token, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *AuthService) ValidateToken(tokenString string) (*table.User, error) {
	claims, err := s.jwtUtils.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	user, err := s.userDAO.FindByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	return user, nil
}

// RefreshToken generates a new token for a valid existing token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	return s.jwtUtils.RefreshToken(tokenString)
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID string, req *dto.ChangePasswordRequest) error {
	// Find user
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if !s.verifyPassword(user, req.CurrentPassword) {
		return errors.New("current password is incorrect")
	}

	// Hash and update new password
	if err := s.hashPassword(user, req.NewPassword); err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := s.userDAO.UpdatePassword(userID, user.PasswordHash); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	log.Info().Str("user_id", userID).Msg("User password changed successfully")
	return nil
}

// validateRegistrationInput validates registration input data
func (s *AuthService) validateRegistrationInput(username, email, password string) error {
	// Validate username format
	if !s.isValidUsername(username) {
		return errors.New("username must be 3-50 characters and contain only letters, numbers, and underscores")
	}

	// Validate email format
	if !s.isValidEmail(email) {
		return errors.New("invalid email format")
	}

	// Validate password strength
	if !s.isValidPassword(password) {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

// isValidUsername checks if username format is valid
func (s *AuthService) isValidUsername(username string) bool {
	// Username should be 3-50 characters, alphanumeric and underscores only
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,50}$`, username)
	return match
}

// isValidEmail checks if email format is valid
func (s *AuthService) isValidEmail(email string) bool {
	// Simple email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidPassword checks if password meets strength requirements
func (s *AuthService) isValidPassword(password string) bool {
	return len(password) >= 6
}

// hashPassword hashes a password using bcrypt
func (s *AuthService) hashPassword(user *table.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)
	return nil
}

// verifyPassword verifies a password against the hash
func (s *AuthService) verifyPassword(user *table.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
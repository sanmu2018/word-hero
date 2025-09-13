package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sanmu2018/word-hero/internal/conf"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTUtils handles JWT token operations
type JWTUtils struct {
	config *conf.JWTConfig
}

// NewJWTUtils creates a new JWTUtils instance
func NewJWTUtils(config *conf.JWTConfig) *JWTUtils {
	return &JWTUtils{
		config: config,
	}
}

// GenerateToken generates a JWT token for a user
func (j *JWTUtils) GenerateToken(userID string, username, email, role string) (string, error) {
	// Parse expiration duration
	duration, err := time.ParseDuration(j.config.ExpiresIn)
	if err != nil {
		return "", fmt.Errorf("invalid token expiration duration: %w", err)
	}

	// Create claims
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    "word-hero",
			Subject:   userID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTUtils) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token with the same claims but new expiration
func (j *JWTUtils) RefreshToken(tokenString string) (string, error) {
	// Validate existing token
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to validate existing token: %w", err)
	}

	// Generate new token
	return j.GenerateToken(claims.UserID, claims.Username, claims.Email, claims.Role)
}

// GetUserIDFromToken extracts user ID from token
func (j *JWTUtils) GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}
	return claims.UserID, nil
}
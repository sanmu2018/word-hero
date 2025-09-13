package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sanmu2018/word-hero/internal/service"
	"github.com/sanmu2018/word-hero/internal/table"
)

// AuthMiddleware handles authentication for protected routes
type AuthMiddleware struct {
	authService *service.AuthService
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware requires authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Bearer token is required",
			})
			c.Abort()
			return
		}

		// Validate token
		user, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)
		c.Set("user_role", user.Role)

		c.Next()
	}
}

// RequireAdmin middleware requires admin role
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if user is authenticated
		if _, exists := c.Get("user"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authentication required",
			})
			c.Abort()
			return
		}

		// Check if user has admin role
		userRole, exists := c.Get("user_role")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware optionally authenticates user
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString != authHeader {
				user, err := m.authService.ValidateToken(tokenString)
				if err == nil {
					c.Set("user", user)
					c.Set("user_id", user.ID)
					c.Set("username", user.Username)
					c.Set("user_role", user.Role)
				}
			}
		}
		c.Next()
	}
}

// GetUserFromContext gets the current user from context
func GetUserFromContext(c *gin.Context) (*table.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found in context")
	}
	return user.(*table.User), nil
}

// GetUserIDFromContext gets the current user ID from context
func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("user ID not found in context")
	}
	return userID.(string), nil
}

// GetUsernameFromContext gets the current username from context
func GetUsernameFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("username not found in context")
	}
	return username.(string), nil
}

// GetUserRoleFromContext gets the current user role from context
func GetUserRoleFromContext(c *gin.Context) (string, error) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", errors.New("user role not found in context")
	}
	return userRole.(string), nil
}

// IsAdmin checks if the current user is an admin
func IsAdmin(c *gin.Context) bool {
	role, err := GetUserRoleFromContext(c)
	if err != nil {
		return false
	}
	return role == "admin"
}

// IsAuthenticated checks if the current user is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, err := GetUserFromContext(c)
	return err == nil
}
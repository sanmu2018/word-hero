package dto

// UserRegisterRequest represents a user registration request
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	FullName string `json:"full_name" binding:"max=100"`
}

// UserLoginRequest represents a user login request
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents a user response (without sensitive data)
type UserResponse struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FullName  string  `json:"full_name"`
	AvatarURL string  `json:"avatar_url"`
	Bio       string  `json:"bio"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
	LastLogin *int64  `json:"last_login,omitempty"`
	CreatedAt int64   `json:"createdAt"`
	UpdatedAt int64   `json:"updatedAt"`
}

// UserUpdateRequest represents a user profile update request
type UserUpdateRequest struct {
	FullName  string `json:"full_name" binding:"max=100"`
	AvatarURL string `json:"avatar_url" binding:"max=255"`
	Bio       string `json:"bio"`
	Email     string `json:"email" binding:"omitempty,email,max=100"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=100"`
}
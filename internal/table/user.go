package table

import (
	"gorm.io/gorm"

	"github.com/sanmu2018/word-hero/internal/utils"
)

// User represents the users table in database
type User struct {
	ID           string    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:100;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	FullName     string    `json:"full_name" gorm:"size:100"`
	AvatarURL    string    `json:"avatar_url" gorm:"size:255"`
	Bio          string    `json:"bio" gorm:"type:text"`
	Role         string    `json:"role" gorm:"size:20;default:'user'"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	LastLogin    *int64    `json:"last_login,omitempty"`
	CreatedAt    int64     `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt    int64     `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM hook - called before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID for ID if not provided
	if u.ID == "" {
		u.ID = utils.GenerateUUID()
	}

	if u.Role == "" {
		u.Role = "user"
	}
	if u.IsActive == false {
		u.IsActive = true
	}
	return nil
}
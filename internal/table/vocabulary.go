package table

import (
	"gorm.io/gorm"

	"github.com/sanmu2018/word-hero/internal/utils"
)

// Word represents the words table in database
type Word struct {
	ID           string  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	English      string  `json:"english" gorm:"size:200;not null"`
	Chinese      string  `json:"chinese" gorm:"size:500;not null"`
	Phonetic     string  `json:"phonetic,omitempty" gorm:"size:100"`
	Example      string  `json:"example,omitempty" gorm:"type:text"`
	Definition   string  `json:"definition,omitempty" gorm:"type:text"`
	Difficulty   string  `json:"difficulty,omitempty" gorm:"size:20"`
	Category     string  `json:"category,omitempty" gorm:"size:50"`
	CreatedAt    int64   `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt    int64   `gorm:"autoUpdateTime:milli" json:"updatedAt"`
}

// TableName returns the table name for Word model
func (Word) TableName() string {
	return "words"
}

// BeforeCreate GORM hook - called before creating a new word
func (w *Word) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID for ID if not provided
	if w.ID == "" {
		w.ID = utils.GenerateUUID()
	}
	return nil
}
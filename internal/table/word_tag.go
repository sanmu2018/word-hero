package table

import (
	"time"

	"gorm.io/gorm"

	"github.com/sanmu2018/word-hero/internal/utils"
)

// WordTag represents the word_tags table in database
type WordTag struct {
	ID        string  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WordID    string  `json:"wordId" gorm:"type:uuid;not null;index:idx_word_tags_word_id"`
	UserID    string  `json:"userId" gorm:"type:uuid;not null;index:idx_word_tags_user_id"`
	Known     *int64  `json:"known"` // NULL means not known, non-null means known with timestamp
	CreatedAt int64   `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt int64   `gorm:"autoUpdateTime:milli" json:"updatedAt"`
}

// TableName returns the table name for WordTag model
func (WordTag) TableName() string {
	return "word_tags"
}

// BeforeCreate GORM hook - called before creating a new word tag
func (wt *WordTag) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID for ID if not provided
	if wt.ID == "" {
		wt.ID = utils.GenerateUUID()
	}

	return nil
}

// MarkAsKnown marks the word as known with current timestamp
func (wt *WordTag) MarkAsKnown() {
	timestamp := time.Now().UnixMilli()
	wt.Known = &timestamp
}

// MarkAsUnknown marks the word as unknown (sets Known to NULL)
func (wt *WordTag) MarkAsUnknown() {
	wt.Known = nil
}

// IsKnown checks if the word is marked as known
func (wt *WordTag) IsKnown() bool {
	return wt.Known != nil
}

// GetKnownTimestamp returns the timestamp when the word was marked as known
func (wt *WordTag) GetKnownTimestamp() int64 {
	if wt.Known == nil {
		return 0
	}
	return *wt.Known
}

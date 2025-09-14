package dao

import (
	"errors"
	"fmt"

	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/log"
	"gorm.io/gorm"
)

// WordTagDAO handles data access operations for word tags
type WordTagDAO struct {
	db *gorm.DB
}

// NewWordTagDAO creates a new WordTagDAO instance
func NewWordTagDAO() *WordTagDAO {
	return &WordTagDAO{
		db: DB,
	}
}

// Create creates a new word tag record
func (dao *WordTagDAO) Create(wordTag *table.WordTag) error {
	if err := dao.db.Create(wordTag).Error; err != nil {
		log.Error(err).Str("word_id", wordTag.WordID).Msg("Failed to create word tag")
		return fmt.Errorf("failed to create word tag: %w", err)
	}

	log.Info().Str("word_tag_id", wordTag.ID).Str("word_id", wordTag.WordID).Msg("Word tag created successfully")
	return nil
}

// GetByID retrieves a word tag by ID
func (dao *WordTagDAO) GetByID(id string) (*table.WordTag, error) {
	var wordTag table.WordTag
	if err := dao.db.First(&wordTag, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("word tag not found")
		}
		return nil, fmt.Errorf("failed to get word tag: %w", err)
	}
	return &wordTag, nil
}

// GetByWordID retrieves a word tag by word ID (deprecated - use GetByWordIDAndUserID)
func (dao *WordTagDAO) GetByWordID(wordID string) (*table.WordTag, error) {
	var wordTag table.WordTag
	if err := dao.db.Where("word_id = ?", wordID).First(&wordTag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("word tag not found for word")
		}
		return nil, fmt.Errorf("failed to get word tag by word ID: %w", err)
	}
	return &wordTag, nil
}

// GetByWordIDAndUserID retrieves a word tag by word ID and user ID
func (dao *WordTagDAO) GetByWordIDAndUserID(wordID, userID string) (*table.WordTag, error) {
	var wordTag table.WordTag
	if err := dao.db.Where("word_id = ? AND user_id = ?", wordID, userID).First(&wordTag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("word tag not found for word and user")
		}
		return nil, fmt.Errorf("failed to get word tag by word ID and user ID: %w", err)
	}
	return &wordTag, nil
}

// GetOrCreateByWordID gets an existing word tag or creates a new one for a user
func (dao *WordTagDAO) GetOrCreateByWordID(wordID, userID string) (*table.WordTag, error) {
	wordTag, err := dao.GetByWordIDAndUserID(wordID, userID)
	if err == nil {
		return wordTag, nil
	}

	// Create new word tag if not found
	newWordTag := &table.WordTag{
		WordID: wordID,
		UserID: userID,
		Known:  nil, // Initially not known
	}

	if err := dao.Create(newWordTag); err != nil {
		return nil, err
	}

	return newWordTag, nil
}

// Update updates an existing word tag
func (dao *WordTagDAO) Update(wordTag *table.WordTag) error {
	if err := dao.db.Save(wordTag).Error; err != nil {
		log.Error(err).Str("word_tag_id", wordTag.ID).Msg("Failed to update word tag")
		return fmt.Errorf("failed to update word tag: %w", err)
	}

	log.Info().Str("word_tag_id", wordTag.ID).Msg("Word tag updated successfully")
	return nil
}

// Delete deletes a word tag by ID
func (dao *WordTagDAO) Delete(id string) error {
	if err := dao.db.Delete(&table.WordTag{}, "id = ?", id).Error; err != nil {
		log.Error(err).Str("word_tag_id", id).Msg("Failed to delete word tag")
		return fmt.Errorf("failed to delete word tag: %w", err)
	}

	log.Info().Str("word_tag_id", id).Msg("Word tag deleted successfully")
	return nil
}

// MarkWordAsKnown marks a word as known by a user
func (dao *WordTagDAO) MarkWordAsKnown(wordID, userID string) error {
	wordTag, err := dao.GetOrCreateByWordID(wordID, userID)
	if err != nil {
		return fmt.Errorf("failed to get or create word tag: %w", err)
	}

	// Mark as known
	wordTag.MarkAsKnown()

	// Update the record
	return dao.Update(wordTag)
}

// RemoveWordMark removes a user's mark from a word
func (dao *WordTagDAO) RemoveWordMark(wordID, userID string) error {
	wordTag, err := dao.GetByWordIDAndUserID(wordID, userID)
	if err != nil {
		return fmt.Errorf("failed to get word tag: %w", err)
	}

	// Mark as unknown
	wordTag.MarkAsUnknown()

	// Update the record
	return dao.Update(wordTag)
}

// IsWordMarkedAsKnown checks if a word is marked as known by a user
func (dao *WordTagDAO) IsWordMarkedAsKnown(wordID, userID string) (bool, error) {
	wordTag, err := dao.GetByWordIDAndUserID(wordID, userID)
	if err != nil {
		if err.Error() == "word tag not found for word and user" {
			return false, nil // Word not tagged by this user
		}
		return false, fmt.Errorf("failed to get word tag: %w", err)
	}

	return wordTag.IsKnown(), nil
}

// GetKnownWords returns all words marked as known by a user
func (dao *WordTagDAO) GetKnownWords(userID string, baseList *BaseList) ([]string, int64, error) {
	var wordTags []table.WordTag
	var total int64

	// Build query with pagination and user filtering
	query := dao.db.Model(&table.WordTag{}).
		Where("user_id = ? AND known IS NOT NULL", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count known words: %w", err)
	}

	// Apply pagination and sorting
	query, err := PageList(query, baseList)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to apply pagination: %w", err)
	}

	// Get paginated word tags
	if err := query.Find(&wordTags).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get known words: %w", err)
	}

	// Extract word IDs
	var wordIDs []string
	for _, wordTag := range wordTags {
		wordIDs = append(wordIDs, wordTag.WordID)
	}

	return wordIDs, total, nil
}

// GetKnownWordsCount returns the count of words marked as known by a user
func (dao *WordTagDAO) GetKnownWordsCount(userID string) (int64, error) {
	var count int64
	if err := dao.db.Model(&table.WordTag{}).
		Where("user_id = ? AND known IS NOT NULL", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count known words: %w", err)
	}
	return count, nil
}

// BulkRemoveWordMarks removes marks from multiple words for a user in one operation
func (dao *WordTagDAO) BulkRemoveWordMarks(wordIDs []string, userID string) (int, error) {
	if len(wordIDs) == 0 {
		return 0, fmt.Errorf("no words to remove marks from")
	}

	// Build the SQL query to set known to NULL for multiple words for a specific user
	result := dao.db.Model(&table.WordTag{}).
		Where("word_id IN ? AND user_id = ?", wordIDs, userID).
		Update("known", nil)

	if result.Error != nil {
		log.Error(result.Error).
			Int("word_count", len(wordIDs)).
			Str("user_id", userID).
			Msg("Failed to bulk remove word marks")
		return 0, fmt.Errorf("failed to bulk remove word marks: %w", result.Error)
	}

	log.Info().
		Int("word_count", len(wordIDs)).
		Str("user_id", userID).
		Int64("affected_rows", result.RowsAffected).
		Msg("Bulk remove word marks completed")

	return int(result.RowsAffected), nil
}

// RemoveAllWordMarks removes all word marks for a user (sets all known to NULL for the user)
func (dao *WordTagDAO) RemoveAllWordMarks(userID string) (int, error) {
	result := dao.db.Model(&table.WordTag{}).
		Where("user_id = ? AND known IS NOT NULL", userID).
		Update("known", nil)

	if result.Error != nil {
		log.Error(result.Error).Str("user_id", userID).Msg("Failed to remove all word marks")
		return 0, fmt.Errorf("failed to remove all word marks: %w", result.Error)
	}

	log.Info().
		Str("user_id", userID).
		Int64("affected_rows", result.RowsAffected).
		Msg("All word marks removed for user")

	return int(result.RowsAffected), nil
}

// GetStats returns word tag statistics
func (dao *WordTagDAO) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total word tags count
	var totalTags int64
	if err := dao.db.Model(&table.WordTag{}).Count(&totalTags).Error; err != nil {
		return nil, fmt.Errorf("failed to count word tags: %w", err)
	}
	stats["total_word_tags"] = totalTags

	// Known words count
	var knownWords int64
	if err := dao.db.Model(&table.WordTag{}).
		Where("known IS NOT NULL").
		Count(&knownWords).Error; err != nil {
		log.Warn().Err(err).Msg("Failed to count known words")
		stats["known_words"] = 0
	} else {
		stats["known_words"] = knownWords
		stats["total_user_marks"] = knownWords // Each word can only be known once now
	}

	// Most recently known words
	var recentWords []table.WordTag
	if err := dao.db.Where("known IS NOT NULL").
		Order("known DESC").
		Limit(10).
		Find(&recentWords).Error; err != nil {
		log.Warn().Err(err).Msg("Failed to get recent known words")
		stats["recent_words"] = []table.WordTag{}
	} else {
		stats["recent_words"] = recentWords
	}

	return stats, nil
}
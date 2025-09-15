package dao

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/log"
	"gorm.io/gorm"
)

// WordDAO handles data access operations for words
type WordDAO struct {
	db *gorm.DB
}

// NewWordDAO creates a new WordDAO instance
func NewWordDAO() *WordDAO {
	return &WordDAO{
		db: DB,
	}
}

// Create creates a new word in the database
func (dao *WordDAO) Create(word *table.Word) error {
	if err := dao.db.Create(word).Error; err != nil {
		log.Error(err).Str("english", word.English).Msg("Failed to create word")
		return fmt.Errorf("failed to create word: %w", err)
	}

	log.Info().Str("word_id", word.ID).Str("english", word.English).Msg("Word created successfully")
	return nil
}

// GetByID retrieves a word by ID
func (dao *WordDAO) GetByID(id string) (*table.Word, error) {
	var word table.Word
	if err := dao.db.First(&word, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("word not found")
		}
		return nil, fmt.Errorf("failed to get word: %w", err)
	}
	return &word, nil
}

// GetByIDs retrieves multiple words by their IDs in a single query
func (dao *WordDAO) GetByIDs(ids []string) ([]table.Word, error) {
	if len(ids) == 0 {
		return []table.Word{}, nil
	}

	var words []table.Word
	if err := dao.db.Where("id IN ?", ids).Find(&words).Error; err != nil {
		return nil, fmt.Errorf("failed to get words by IDs: %w", err)
	}

	return words, nil
}

// GetByEnglish retrieves a word by English text
func (dao *WordDAO) GetByEnglish(english string) (*table.Word, error) {
	var word table.Word
	if err := dao.db.Where("english = ?", english).First(&word).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("word not found")
		}
		return nil, fmt.Errorf("failed to get word by english: %w", err)
	}
	return &word, nil
}

// Update updates an existing word
func (dao *WordDAO) Update(word *table.Word) error {
	if err := dao.db.Save(word).Error; err != nil {
		log.Error(err).Str("word_id", word.ID).Msg("Failed to update word")
		return fmt.Errorf("failed to update word: %w", err)
	}

	log.Info().Str("word_id", word.ID).Msg("Word updated successfully")
	return nil
}

// Delete deletes a word by ID
func (dao *WordDAO) Delete(id string) error {
	if err := dao.db.Delete(&table.Word{}, "id = ?", id).Error; err != nil {
		log.Error(err).Str("word_id", id).Msg("Failed to delete word")
		return fmt.Errorf("failed to delete word: %w", err)
	}

	log.Info().Str("word_id", id).Msg("Word deleted successfully")
	return nil
}

// GetWordsByPage returns words for a specific page with pagination using BaseList
func (dao *WordDAO) GetWordsByPage(baseList *BaseList) ([]table.Word, int64, error) {
	// baseList can be nil, meaning no pagination (return all data)
	// No default values are set - pagination is completely optional

	var words []table.Word
	var total int64

	// Get total count
	if err := dao.db.Model(&table.Word{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count words: %w", err)
	}

	// Build query with pagination
	query := dao.db.Model(&table.Word{})

	// Apply pagination and sorting
	query, err := PageList(query, baseList)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to apply pagination: %w", err)
	}

	// Get paginated words
	if err := query.Find(&words).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get words: %w", err)
	}

	return words, total, nil
}

// GetWordsByPageLegacy 保持向后兼容的旧版本方法
func (dao *WordDAO) GetWordsByPageLegacy(pageNumber, pageSize int) ([]table.Word, int64, error) {
	baseList := &BaseList{
		PageNum:  pageNumber,
		PageSize: pageSize,
	}
	return dao.GetWordsByPage(baseList)
}

// GetAllWords returns all words (use carefully for large datasets)
func (dao *WordDAO) GetAllWords() ([]table.Word, error) {
	var words []table.Word
	if err := dao.db.Order("english ASC").Find(&words).Error; err != nil {
		return nil, fmt.Errorf("failed to get all words: %w", err)
	}
	return words, nil
}

// SearchWords searches for words matching the query in English and Chinese
func (dao *WordDAO) SearchWords(query string) ([]table.Word, error) {
	if len(query) < 2 {
		return []table.Word{}, nil
	}

	var words []table.Word
	searchPattern := "%" + strings.ToLower(query) + "%"

	if err := dao.db.Where(
		"LOWER(english) LIKE ? OR LOWER(chinese) LIKE ?",
		searchPattern, searchPattern,
	).Order("english ASC").Find(&words).Error; err != nil {
		return nil, fmt.Errorf("failed to search words: %w", err)
	}

	log.Info().Str("query", query).Int("results", len(words)).Msg("Search completed")
	return words, nil
}

// SearchWordsWithRegex searches for words using regex patterns
func (dao *WordDAO) SearchWordsWithRegex(query string) ([]table.Word, error) {
	if len(query) < 2 {
		return []table.Word{}, nil
	}

	// Compile regex pattern
	pattern, err := regexp.Compile("(?i)" + query) // (?i) for case-insensitive
	if err != nil {
		return []table.Word{}, fmt.Errorf("invalid regex pattern: %w", err)
	}

	allWords, err := dao.GetAllWords()
	if err != nil {
		return nil, err
	}

	var results []table.Word
	for _, word := range allWords {
		englishMatch := pattern.MatchString(word.English)
		chineseMatch := pattern.MatchString(word.Chinese)

		if englishMatch || chineseMatch {
			results = append(results, word)
		}
	}

	log.Info().Str("query", query).Int("results", len(results)).Msg("Regex search completed")
	return results, nil
}

// GetWordsByChinese finds words by Chinese text
func (dao *WordDAO) GetWordsByChinese(chinese string) ([]table.Word, error) {
	var words []table.Word
	searchPattern := "%" + chinese + "%"

	if err := dao.db.Where("chinese LIKE ?", searchPattern).Order("english ASC").Find(&words).Error; err != nil {
		return nil, fmt.Errorf("failed to get words by chinese: %w", err)
	}

	return words, nil
}

// GetRandomWords returns a random selection of words
func (dao *WordDAO) GetRandomWords(count int) ([]table.Word, error) {
	if count <= 0 {
		return []table.Word{}, nil
	}

	var words []table.Word
	if err := dao.db.Order("RANDOM()").Limit(count).Find(&words).Error; err != nil {
		return nil, fmt.Errorf("failed to get random words: %w", err)
	}

	return words, nil
}

// GetWordCount returns the total number of words
func (dao *WordDAO) GetWordCount() (int64, error) {
	var count int64
	if err := dao.db.Model(&table.Word{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count words: %w", err)
	}
	return count, nil
}

// GetWordsByCategory returns words filtered by category
func (dao *WordDAO) GetWordsByCategory(category string) ([]table.Word, error) {
	var words []table.Word
	if err := dao.db.Where("category = ?", category).Order("english ASC").Find(&words).Error; err != nil {
		return nil, fmt.Errorf("failed to get words by category: %w", err)
	}
	return words, nil
}

// GetWordsByDifficulty returns words filtered by difficulty level
func (dao *WordDAO) GetWordsByDifficulty(difficulty string) ([]table.Word, error) {
	var words []table.Word
	if err := dao.db.Where("difficulty = ?", difficulty).Order("english ASC").Find(&words).Error; err != nil {
		return nil, fmt.Errorf("failed to get words by difficulty: %w", err)
	}
	return words, nil
}

// BulkImport imports multiple words in a batch
func (dao *WordDAO) BulkImport(words []table.Word) error {
	if len(words) == 0 {
		return fmt.Errorf("no words to import")
	}

	log.Info().Int("count", len(words)).Msg("Starting bulk import of words")

	// Use batch insertion for better performance
	batchSize := 100
	for i := 0; i < len(words); i += batchSize {
		end := i + batchSize
		if end > len(words) {
			end = len(words)
		}

		batch := words[i:end]
		if err := dao.db.CreateInBatches(batch, batchSize).Error; err != nil {
			log.Error(err).Int("batch_start", i).Int("batch_end", end).Msg("Failed to import batch")
			return fmt.Errorf("failed to import batch %d-%d: %w", i, end, err)
		}

		log.Debug().Int("batch_start", i).Int("batch_end", end).Msg("Batch imported successfully")
	}

	log.Info().Int("total_imported", len(words)).Msg("Bulk import completed successfully")
	return nil
}

// DeleteAllWords deletes all words from the database (use with caution)
func (dao *WordDAO) DeleteAllWords() error {
	log.Warn().Msg("Deleting all words from database")

	if err := dao.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&table.Word{}).Error; err != nil {
		log.Error(err).Msg("Failed to delete all words")
		return fmt.Errorf("failed to delete all words: %w", err)
	}

	log.Info().Msg("All words deleted successfully")
	return nil
}

// GetStats returns database statistics for words
func (dao *WordDAO) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total word count
	totalCount, err := dao.GetWordCount()
	if err != nil {
		return nil, err
	}
	stats["total_words"] = totalCount

	// Category distribution
	var categoryStats []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	if err := dao.db.Model(&table.Word{}).
		Select("category, COUNT(*) as count").
		Where("category != ''").
		Group("category").
		Find(&categoryStats).Error; err != nil {
		log.Warn().Err(err).Msg("Failed to get category stats")
	} else {
		stats["categories"] = categoryStats
	}

	// Difficulty distribution
	var difficultyStats []struct {
		Difficulty string `json:"difficulty"`
		Count      int64  `json:"count"`
	}
	if err := dao.db.Model(&table.Word{}).
		Select("difficulty, COUNT(*) as count").
		Where("difficulty != ''").
		Group("difficulty").
		Find(&difficultyStats).Error; err != nil {
		log.Warn().Err(err).Msg("Failed to get difficulty stats")
	} else {
		stats["difficulties"] = difficultyStats
	}

	return stats, nil
}

// IsEmpty checks if the words table is empty
func (dao *WordDAO) IsEmpty() (bool, error) {
	var count int64
	if err := dao.db.Model(&table.Word{}).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if words table is empty: %w", err)
	}
	return count == 0, nil
}

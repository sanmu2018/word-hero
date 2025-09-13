package service

import (
	"fmt"

	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/log"
)

// VocabularyService handles vocabulary-related business logic
type VocabularyService struct {
	wordDAO *dao.WordDAO
}

// NewVocabularyService creates a new vocabulary service instance
func NewVocabularyService(wordDAO *dao.WordDAO) *VocabularyService {
	log.Info().Msg("Creating vocabulary service with database backend")

	return &VocabularyService{
		wordDAO: wordDAO,
	}
}

// GetWordsByPage returns words for a specific page using BaseList
func (vs *VocabularyService) GetWordsByPage(baseList *dao.BaseList) (*dto.VocabularyPage, error) {
	// baseList can be nil, meaning no pagination (return all data)
	// No default values are set - pagination is completely optional

	// Get words and total count from database
	words, totalCount, err := vs.wordDAO.GetWordsByPage(baseList)
	if err != nil {
		return nil, fmt.Errorf("failed to get words by page: %w", err)
	}

	// Calculate pagination info only if pagination is requested
	var pageNumber, pageSize, totalPages int
	if baseList != nil {
		pageNumber = baseList.PageNum
		pageSize = baseList.PageSize
		if pageSize > 0 {
			totalPages = int((totalCount + int64(pageSize) - 1) / int64(pageSize))
		}
	}

	return &dto.VocabularyPage{
		Words:      words,
		TotalCount: int(totalCount),
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetWordsByPageLegacy 保持向后兼容的旧版本方法
func (vs *VocabularyService) GetWordsByPageLegacy(pageNumber, pageSize int) (*dto.VocabularyPage, error) {
	baseList := &dao.BaseList{
		PageNum:  pageNumber,
		PageSize: pageSize,
	}
	return vs.GetWordsByPage(baseList)
}

// SearchWords searches for words matching the query in English and Chinese
func (vs *VocabularyService) SearchWords(query string) ([]table.Word, error) {
	if len(query) < 2 {
		return []table.Word{}, nil
	}

	words, err := vs.wordDAO.SearchWords(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search words: %w", err)
	}

	log.Info().Str("query", query).Int("results", len(words)).Msg("Database search completed")
	return words, nil
}

// SearchWordsWithRegex searches for words using regex patterns
func (vs *VocabularyService) SearchWordsWithRegex(query string) ([]table.Word, error) {
	if len(query) < 2 {
		return []table.Word{}, nil
	}

	words, err := vs.wordDAO.SearchWordsWithRegex(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search words with regex: %w", err)
	}

	log.Info().Str("query", query).Int("results", len(words)).Msg("Regex search completed")
	return words, nil
}

// GetWordByEnglish finds a word by its English text
func (vs *VocabularyService) GetWordByEnglish(english string) (*table.Word, bool) {
	word, err := vs.wordDAO.GetByEnglish(english)
	if err != nil {
		return nil, false
	}
	return word, true
}

// GetWordsByChinese finds words by Chinese text
func (vs *VocabularyService) GetWordsByChinese(chinese string) ([]table.Word, error) {
	words, err := vs.wordDAO.GetWordsByChinese(chinese)
	if err != nil {
		return nil, fmt.Errorf("failed to get words by chinese: %w", err)
	}
	return words, nil
}

// GetRandomWords returns a random selection of words
func (vs *VocabularyService) GetRandomWords(count int) ([]table.Word, error) {
	if count <= 0 {
		return []table.Word{}, nil
	}

	words, err := vs.wordDAO.GetRandomWords(count)
	if err != nil {
		return nil, fmt.Errorf("failed to get random words: %w", err)
	}

	return words, nil
}

// GetWordCount returns the total number of words
func (vs *VocabularyService) GetWordCount() int {
	count, err := vs.wordDAO.GetWordCount()
	if err != nil {
		log.Error(err).Msg("Failed to get word count")
		return 0
	}
	return int(count)
}

// GetStats returns application statistics
func (vs *VocabularyService) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// Get total word count
	totalWords, err := vs.wordDAO.GetWordCount()
	if err != nil {
		log.Error(err).Msg("Failed to get word count for stats")
		totalWords = 0
	}
	stats["total_words"] = totalWords

	// Get additional database statistics
	dbStats, err := vs.wordDAO.GetStats()
	if err != nil {
		log.Error(err).Msg("Failed to get database stats")
	} else {
		for k, v := range dbStats {
			stats[k] = v
		}
	}

	// Add data source information
	stats["data_source"] = "PostgreSQL Database"
	stats["migration_status"] = "completed"

	return stats
}

// ValidateWordList checks if the word data is accessible
func (vs *VocabularyService) ValidateWordList() error {
	isEmpty, err := vs.wordDAO.IsEmpty()
	if err != nil {
		return fmt.Errorf("failed to check if words table is empty: %w", err)
	}

	if isEmpty {
		return fmt.Errorf("words table is empty - please run data migration")
	}

	count, err := vs.wordDAO.GetWordCount()
	if err != nil {
		return fmt.Errorf("failed to get word count: %w", err)
	}

	log.Info().Int64("totalWords", count).Msg("Word data validation completed")
	return nil
}

// RefreshWordList is no longer needed for database backend
// This method is kept for compatibility but logs a deprecation warning
func (vs *VocabularyService) RefreshWordList() error {
	log.Warn().Msg("RefreshWordList is deprecated for database backend - word data is now managed in database")
	return nil
}

// GetWordsByCategory returns words filtered by category
func (vs *VocabularyService) GetWordsByCategory(category string) ([]table.Word, error) {
	words, err := vs.wordDAO.GetWordsByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to get words by category: %w", err)
	}
	return words, nil
}

// GetWordsByDifficulty returns words filtered by difficulty level
func (vs *VocabularyService) GetWordsByDifficulty(difficulty string) ([]table.Word, error) {
	words, err := vs.wordDAO.GetWordsByDifficulty(difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to get words by difficulty: %w", err)
	}
	return words, nil
}

// GetAllCategories returns all unique categories in the database
func (vs *VocabularyService) GetAllCategories() ([]string, error) {
	stats, err := vs.wordDAO.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	var categories []string
	if categoryStats, ok := stats["categories"]; ok {
		if catStats, ok := categoryStats.([]interface{}); ok {
			for _, catStat := range catStats {
				if statMap, ok := catStat.(map[string]interface{}); ok {
					if category, ok := statMap["category"].(string); ok {
						categories = append(categories, category)
					}
				}
			}
		}
	}

	return categories, nil
}

// GetAllDifficultyLevels returns all unique difficulty levels in the database
func (vs *VocabularyService) GetAllDifficultyLevels() ([]string, error) {
	stats, err := vs.wordDAO.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get difficulty levels: %w", err)
	}

	var difficulties []string
	if diffStats, ok := stats["difficulties"]; ok {
		if diffList, ok := diffStats.([]interface{}); ok {
			for _, diffStat := range diffList {
				if statMap, ok := diffStat.(map[string]interface{}); ok {
					if difficulty, ok := statMap["difficulty"].(string); ok {
						difficulties = append(difficulties, difficulty)
					}
				}
			}
		}
	}

	return difficulties, nil
}

// CreateWord adds a new word to the database
func (vs *VocabularyService) CreateWord(word *table.Word) error {
	return vs.wordDAO.Create(word)
}

// UpdateWord updates an existing word in the database
func (vs *VocabularyService) UpdateWord(word *table.Word) error {
	return vs.wordDAO.Update(word)
}

// DeleteWord deletes a word from the database
func (vs *VocabularyService) DeleteWord(id string) error {
	return vs.wordDAO.Delete(id)
}

// GetWordByID retrieves a word by ID
func (vs *VocabularyService) GetWordByID(id string) (*table.Word, error) {
	return vs.wordDAO.GetByID(id)
}

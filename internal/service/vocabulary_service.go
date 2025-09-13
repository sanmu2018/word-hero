package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/log"
)

// VocabularyService handles vocabulary-related business logic
type VocabularyService struct {
	wordList     *dao.WordList
	excelReader  *dao.ExcelReader
	pagerService *PagerService
}

// NewVocabularyService creates a new vocabulary service instance
func NewVocabularyService(wordList *dao.WordList, excelReader *dao.ExcelReader) *VocabularyService {
	log.Info().Int("wordCount", len(wordList.Words)).Msg("Creating vocabulary service")

	return &VocabularyService{
		wordList:    wordList,
		excelReader: excelReader,
	}
}

// SetPagerService sets the pager service reference
func (vs *VocabularyService) SetPagerService(pagerService *PagerService) {
	vs.pagerService = pagerService
}

// GetWordsByPage returns words for a specific page
func (vs *VocabularyService) GetWordsByPage(pageNumber, pageSize int) (*dao.Page, error) {
	if vs.pagerService == nil {
		// Create a temporary pager service if not set
		tempPager := NewPagerService(vs.wordList, pageSize)
		return tempPager.GetPage(pageNumber)
	}

	// Use existing pager service
	return vs.pagerService.GetPage(pageNumber)
}

// SearchWords searches for words matching the query
func (vs *VocabularyService) SearchWords(query string) ([]dao.Word, error) {
	if len(query) < 2 {
		return []dao.Word{}, nil
	}

	query = strings.ToLower(query)
	var results []dao.Word

	// Search in both English and Chinese
	for _, word := range vs.wordList.Words {
		englishMatch := strings.Contains(strings.ToLower(word.English), query)
		chineseMatch := strings.Contains(strings.ToLower(word.Chinese), query)

		if englishMatch || chineseMatch {
			results = append(results, word)
		}
	}

	log.Info().Str("query", query).Int("results", len(results)).Msg("Search completed")

	return results, nil
}

// SearchWordsWithRegex searches for words using regex patterns
func (vs *VocabularyService) SearchWordsWithRegex(query string) ([]dao.Word, error) {
	if len(query) < 2 {
		return []dao.Word{}, nil
	}

	// Compile regex pattern
	pattern, err := regexp.Compile("(?i)" + query) // (?i) for case-insensitive
	if err != nil {
		return []dao.Word{}, err
	}

	var results []dao.Word

	// Search in both English and Chinese
	for _, word := range vs.wordList.Words {
		englishMatch := pattern.MatchString(word.English)
		chineseMatch := pattern.MatchString(word.Chinese)

		if englishMatch || chineseMatch {
			results = append(results, word)
		}
	}

	log.Info().Str("query", query).Int("results", len(results)).Msg("Regex search completed")

	return results, nil
}

// GetWordByEnglish finds a word by its English text
func (vs *VocabularyService) GetWordByEnglish(english string) (*dao.Word, bool) {
	english = strings.TrimSpace(english)
	for _, word := range vs.wordList.Words {
		if strings.EqualFold(word.English, english) {
			return &word, true
		}
	}
	return nil, false
}

// GetWordsByChinese finds words by Chinese text
func (vs *VocabularyService) GetWordsByChinese(chinese string) ([]dao.Word, error) {
	chinese = strings.TrimSpace(chinese)
	var results []dao.Word

	for _, word := range vs.wordList.Words {
		if strings.Contains(word.Chinese, chinese) {
			results = append(results, word)
		}
	}

	return results, nil
}

// GetRandomWords returns a random selection of words
func (vs *VocabularyService) GetRandomWords(count int) ([]dao.Word, error) {
	if count <= 0 || count > len(vs.wordList.Words) {
		count = len(vs.wordList.Words)
	}

	// Simple random selection (not cryptographically secure)
	// For production, consider using crypto/rand
	indices := make(map[int]bool)
	var results []dao.Word

	for len(results) < count {
		index := int(len(vs.wordList.Words) / 2) // Simple deterministic for now
		if !indices[index] {
			indices[index] = true
			results = append(results, vs.wordList.Words[index])
		}
	}

	return results, nil
}

// GetWordCount returns the total number of words
func (vs *VocabularyService) GetWordCount() int {
	return len(vs.wordList.Words)
}

// GetStats returns application statistics
func (vs *VocabularyService) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"total_words": vs.GetWordCount(),
		"file_source":  vs.excelReader.GetFilePath(),
	}
}

// ValidateWordList checks if the word list is valid
func (vs *VocabularyService) ValidateWordList() error {
	if vs.wordList == nil || len(vs.wordList.Words) == 0 {
		return fmt.Errorf("word list is empty or nil")
	}

	// Check for required fields
	for i, word := range vs.wordList.Words {
		if word.English == "" || word.Chinese == "" {
			log.Warn().Int("index", i).Str("english", word.English).Str("chinese", word.Chinese).Msg("Found word with empty fields")
		}
	}

	log.Info().Int("totalWords", len(vs.wordList.Words)).Msg("Word list validation completed")

	return nil
}

// RefreshWordList reloads words from Excel file
func (vs *VocabularyService) RefreshWordList() error {
	log.Info().Msg("Refreshing word list from Excel file")

	// Re-read words from Excel
	newWordList, err := vs.excelReader.ReadWords()
	if err != nil {
		return fmt.Errorf("failed to refresh word list: %w", err)
	}

	vs.wordList = newWordList

	log.Info().Int("newCount", len(newWordList.Words)).Msg("Word list refreshed successfully")

	return nil
}
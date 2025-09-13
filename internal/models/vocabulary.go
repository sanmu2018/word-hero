package models

import (
	"fmt"
	"strings"

	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/table"
)

// VocabularyBusiness represents vocabulary business logic
type VocabularyBusiness struct {
	words []table.Word
}

// NewVocabularyBusiness creates a new vocabulary business instance
func NewVocabularyBusiness(words []table.Word) *VocabularyBusiness {
	return &VocabularyBusiness{words: words}
}

// GetWords returns all vocabulary words
func (v *VocabularyBusiness) GetWords() []table.Word {
	return v.words
}

// GetWordCount returns the total number of words
func (v *VocabularyBusiness) GetWordCount() int {
	return len(v.words)
}

// SearchWords searches for words containing the query in English or Chinese
func (v *VocabularyBusiness) SearchWords(query string) []table.Word {
	var results []table.Word
	queryLower := strings.ToLower(query)

	for _, word := range v.words {
		if strings.Contains(strings.ToLower(word.English), queryLower) ||
		   strings.Contains(strings.ToLower(word.Chinese), queryLower) {
			results = append(results, word)
		}
	}
	return results
}

// GetWordsByPage returns words for a specific page with pagination metadata
func (v *VocabularyBusiness) GetWordsByPage(pageNumber, pageSize int) (*dto.VocabularyPage, error) {
	totalWords := len(v.words)
	totalPages := (totalWords + pageSize - 1) / pageSize

	if pageNumber < 1 || pageNumber > totalPages {
		return nil, fmt.Errorf("page number %d is out of range (1-%d)", pageNumber, totalPages)
	}

	startIndex := (pageNumber - 1) * pageSize
	endIndex := startIndex + pageSize

	if endIndex > totalWords {
		endIndex = totalWords
	}

	pageWords := v.words[startIndex:endIndex]

	return &dto.VocabularyPage{
		Words:      pageWords,
		TotalCount: totalWords,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetStats returns vocabulary statistics
func (v *VocabularyBusiness) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["total_words"] = len(v.words)
	stats["words_by_page"] = map[string]interface{}{
		"25":   (len(v.words) + 24) / 25,
		"50":   (len(v.words) + 49) / 50,
		"100":  (len(v.words) + 99) / 100,
	}
	return stats
}
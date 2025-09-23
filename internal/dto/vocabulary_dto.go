package dto

import (
	"github.com/sanmu2018/word-hero/internal/table"
)

// VocabularyList represents a list of vocabulary words
type VocabularyList struct {
	Words []table.Word `json:"words"`
}

// VocabularyPage represents a paginated result of vocabulary words
type VocabularyPage struct {
	Words      []table.Word `json:"words"`
	TotalCount int64        `json:"totalCount"`
	PageNumber int          `json:"pageNumber"`
	PageSize   int          `json:"pageSize"`
	TotalPages int          `json:"totalPages"`
}

// VocabularyStats represents vocabulary statistics
type VocabularyStats struct {
	TotalWords  int                    `json:"totalWords"`
	TotalPages  int                    `json:"totalPages"`
	WordsByPage map[string]interface{} `json:"wordsByPage"`
}

// WordWithMarkStatus represents a word with its mark status for a user
type WordWithMarkStatus struct {
	table.Word
	IsMarked  bool  `json:"isMarked"`
	MarkCount int   `json:"markCount"`
	MarkedAt  int64 `json:"markedAt,omitempty"`
}

// VocabularyPageWithMarks represents a paginated result with mark status
type VocabularyPageWithMarks struct {
	Words      []WordWithMarkStatus `json:"words"`
	TotalCount int                  `json:"totalCount"`
	PageNumber int                  `json:"pageNumber"`
	PageSize   int                  `json:"pageSize"`
	TotalPages int                  `json:"totalPages"`
}

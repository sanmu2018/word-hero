package dto

import (
	"github.com/sanmu2018/word-hero/internal/table"
)

// WordList represents a list of vocabulary words (legacy structure for compatibility)
type WordList struct {
	Words []table.Word `json:"words"`
}

// Page represents a paginated result of vocabulary words (legacy structure for compatibility)
type Page struct {
	Words      []table.Word `json:"words"`
	TotalCount int          `json:"total_count"`
	PageNumber int          `json:"page_number"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}
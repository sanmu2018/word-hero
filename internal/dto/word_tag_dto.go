package dto

import (
	"github.com/sanmu2018/word-hero/internal/table"
)

// WordMarkRequest represents a request to mark a word
type WordMarkRequest struct {
	WordID string `json:"wordId" binding:"required,uuid"`
	UserID string `json:"userId"`
}

// WordMarkResponse represents a response for word mark operations
type WordMarkResponse struct {
	WordID    string `json:"wordId"`
	IsMarked  bool   `json:"isMarked"`
	MarkCount int    `json:"markCount"`
	Message   string `json:"message"`
}

// UserProgressRequest represents a request for user progress
type UserProgressRequest struct {
	UserID string `json:"userId" binding:"required,uuid"`
}

// UserProgressResponse represents user learning progress
type UserProgressResponse struct {
	UserID         string  `json:"userId"`
	KnownWords     int64   `json:"knownWords"`
	TotalWords     int64   `json:"totalWords"`
	ProgressRate   float64 `json:"progressRate"`
	RecentActivity int     `json:"recentActivity"`
}

// KnownWordsRequest represents a request to get known words
type KnownWordsRequest struct {
	UserID string `json:"userId" binding:"required,uuid"`
}

// KnownWordsResponse represents a response for known words
type KnownWordsResponse struct {
	WordIDs    []string `json:"wordIds"`
	TotalCount int64    `json:"totalCount"`
}

// WordMarkStatusRequest represents a request to get mark status for specific words
type WordMarkStatusRequest struct {
	WordIDs []string `json:"wordIds" binding:"required,min=1"`
}

// WordMarkStatusResponse represents a response for word mark status
type WordMarkStatusResponse struct {
	WordMarkStatuses []WordMarkStatus `json:"wordMarkStatuses"`
}

// WordMarkStatus represents the mark status of a word for a user
type WordMarkStatus struct {
	WordID    string `json:"wordId"`
	IsMarked  bool   `json:"isMarked"`
	MarkCount int    `json:"markCount"`
	MarkedAt  int64  `json:"markedAt,omitempty"`
}

// BatchWordMarkRequest represents a request to mark multiple words
type BatchWordMarkRequest struct {
	WordIDs []string `json:"wordIds" binding:"required,min=1"`
	UserID  string   `json:"userId" binding:"required,uuid"`
}

// BatchWordMarkResponse represents a response for batch word mark operations
type BatchWordMarkResponse struct {
	SuccessCount int                `json:"successCount"`
	FailedCount  int                `json:"failedCount"`
	Results      []WordMarkResponse `json:"results"`
	Errors       map[string]string  `json:"errors,omitempty"`
}

// WordTagStats represents statistics for word tags
type WordTagStats struct {
	TotalWordTags  int64           `json:"totalWordTags"`
	TotalUserMarks int64           `json:"totalUserMarks"`
	TopWords       []table.WordTag `json:"topWords"`
}

// KnownWordInfo represents information about a known word
type KnownWordInfo struct {
	WordID    string `json:"wordId"`
	KnownAt   int64  `json:"knownAt"`
}

// UserWordStats represents detailed statistics for a user's word learning
type UserWordStats struct {
	UserID           string           `json:"userId"`
	KnownWordsCount  int64            `json:"knownWordsCount"`
	TotalWordsCount  int64            `json:"totalWordsCount"`
	ProgressRate     float64          `json:"progressRate"`
	RecentMarks      []KnownWordInfo  `json:"recentMarks"`
	KnownWordsByDate map[string]int   `json:"knownWordsByDate"`
	TopCategories    map[string]int   `json:"topCategories,omitempty"`
}

// ForgetWordsRequest represents a request to forget specific words
type ForgetWordsRequest struct {
	WordIDs []string `json:"wordIds" binding:"required,min=1"`
}

// ForgetWordsResponse represents response for forget words operation
type ForgetWordsResponse struct {
	WordIDs        []string `json:"wordIds"`
	ForgottenCount int      `json:"forgottenCount"`
	Message        string   `json:"message"`
}

// ForgetAllRequest represents a request to forget all words
type ForgetAllRequest struct {
	Confirm bool `json:"confirm" binding:"required"`
}

// ForgetAllResponse represents response for all forget operation
type ForgetAllResponse struct {
	ForgottenCount int    `json:"forgottenCount"`
	Message        string `json:"message"`
}

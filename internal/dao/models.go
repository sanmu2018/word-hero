package dao

// Word represents a vocabulary word with its translation
type Word struct {
	English string `json:"english"`
	Chinese string `json:"chinese"`
}

// WordList represents a collection of words
type WordList struct {
	Words []Word `json:"words"`
}

// Page represents a single page of words
type Page struct {
	Number     int    `json:"number"`
	TotalPages int    `json:"total_pages"`
	Words      []Word `json:"words"`
	PageSize   int    `json:"page_size"`
}
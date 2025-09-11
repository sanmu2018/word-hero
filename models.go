package main

// Word represents a vocabulary word with its translation
type Word struct {
	English  string
	Chinese  string
}

// WordList represents a collection of words
type WordList struct {
	Words []Word
}

// Page represents a single page of words
type Page struct {
	Number      int
	TotalPages  int
	Words       []Word
	PageSize    int
}
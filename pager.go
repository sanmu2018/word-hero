package main

import (
	"fmt"
	"math"
)

// Pager handles pagination of word lists
type Pager struct {
	wordList *WordList
	pageSize int
}

// NewPager creates a new pager instance
func NewPager(wordList *WordList, pageSize int) *Pager {
	return &Pager{
		wordList: wordList,
		pageSize: pageSize,
	}
}

// GetTotalPages returns the total number of pages
func (p *Pager) GetTotalPages() int {
	if len(p.wordList.Words) == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(p.wordList.Words)) / float64(p.pageSize)))
}

// GetPage returns a specific page of words
func (p *Pager) GetPage(pageNumber int) (*Page, error) {
	totalPages := p.GetTotalPages()
	
	if pageNumber < 1 || pageNumber > totalPages {
		return nil, fmt.Errorf("page number %d is out of range (1-%d)", pageNumber, totalPages)
	}

	startIndex := (pageNumber - 1) * p.pageSize
	endIndex := startIndex + p.pageSize
	
	if endIndex > len(p.wordList.Words) {
		endIndex = len(p.wordList.Words)
	}

	pageWords := p.wordList.Words[startIndex:endIndex]

	return &Page{
		Number:     pageNumber,
		TotalPages: totalPages,
		Words:      pageWords,
		PageSize:   p.pageSize,
	}, nil
}

// GetFirstPage returns the first page
func (p *Pager) GetFirstPage() (*Page, error) {
	return p.GetPage(1)
}

// GetLastPage returns the last page
func (p *Pager) GetLastPage() (*Page, error) {
	return p.GetPage(p.GetTotalPages())
}

// GetNextPage returns the next page
func (p *Pager) GetNextPage(currentPage int) (*Page, error) {
	return p.GetPage(currentPage + 1)
}

// GetPreviousPage returns the previous page
func (p *Pager) GetPreviousPage(currentPage int) (*Page, error) {
	if currentPage <= 1 {
		return nil, fmt.Errorf("already on first page")
	}
	return p.GetPage(currentPage - 1)
}

// HasNextPage checks if there is a next page
func (p *Pager) HasNextPage(currentPage int) bool {
	return currentPage < p.GetTotalPages()
}

// HasPreviousPage checks if there is a previous page
func (p *Pager) HasPreviousPage(currentPage int) bool {
	return currentPage > 1
}

// GetWordCount returns the total number of words
func (p *Pager) GetWordCount() int {
	return len(p.wordList.Words)
}

// GetPageRange returns the start and end word indices for a page
func (p *Pager) GetPageRange(pageNumber int) (start, end int, err error) {
	totalPages := p.GetTotalPages()
	
	if pageNumber < 1 || pageNumber > totalPages {
		return 0, 0, fmt.Errorf("page number %d is out of range (1-%d)", pageNumber, totalPages)
	}

	start = (pageNumber - 1) * p.pageSize + 1
	end = pageNumber * p.pageSize
	
	if end > p.GetWordCount() {
		end = p.GetWordCount()
	}

	return start, end, nil
}
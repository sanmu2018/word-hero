package service

import (
	"fmt"
	"math"

	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/log"
)

// PagerService handles pagination logic for vocabulary words
type PagerService struct {
	wordList *dao.WordList
	pageSize int
}

// NewPagerService creates a new pager service instance
func NewPagerService(wordList *dao.WordList, pageSize int) *PagerService {
	log.Info().Int("wordCount", len(wordList.Words)).Int("pageSize", pageSize).Msg("Creating pager service")

	return &PagerService{
		wordList: wordList,
		pageSize: pageSize,
	}
}

// GetPageSize returns the current page size
func (ps *PagerService) GetPageSize() int {
	return ps.pageSize
}

// GetTotalPages returns the total number of pages
func (ps *PagerService) GetTotalPages() int {
	if len(ps.wordList.Words) == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(ps.wordList.Words)) / float64(ps.pageSize)))
}

// GetPage returns a specific page of words with metadata
func (ps *PagerService) GetPage(pageNumber int) (*dao.Page, error) {
	totalPages := ps.GetTotalPages()

	if pageNumber < 1 || pageNumber > totalPages {
		return nil, fmt.Errorf("page number %d is out of range (1-%d)", pageNumber, totalPages)
	}

	startIndex := (pageNumber - 1) * ps.pageSize
	endIndex := startIndex + ps.pageSize

	if endIndex > len(ps.wordList.Words) {
		endIndex = len(ps.wordList.Words)
	}

	pageWords := ps.wordList.Words[startIndex:endIndex]

	return &dao.Page{
		Number:     pageNumber,
		TotalPages: totalPages,
		Words:      pageWords,
		PageSize:   ps.pageSize,
	}, nil
}

// GetPageData returns page data with additional metadata for API responses
func (ps *PagerService) GetPageData(pageNumber int) (map[string]interface{}, error) {
	page, err := ps.GetPage(pageNumber)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"items": page.Words,
		"total": ps.GetWordCount(),
	}

	return data, nil
}

// GetFirstPage returns the first page
func (ps *PagerService) GetFirstPage() (*dao.Page, error) {
	return ps.GetPage(1)
}

// GetLastPage returns the last page
func (ps *PagerService) GetLastPage() (*dao.Page, error) {
	return ps.GetPage(ps.GetTotalPages())
}

// GetNextPage returns the next page
func (ps *PagerService) GetNextPage(currentPage int) (*dao.Page, error) {
	return ps.GetPage(currentPage + 1)
}

// GetPreviousPage returns the previous page
func (ps *PagerService) GetPreviousPage(currentPage int) (*dao.Page, error) {
	if currentPage <= 1 {
		return nil, fmt.Errorf("already on first page")
	}
	return ps.GetPage(currentPage - 1)
}

// HasNextPage checks if there is a next page
func (ps *PagerService) HasNextPage(currentPage int) bool {
	return currentPage < ps.GetTotalPages()
}

// HasPreviousPage checks if there is a previous page
func (ps *PagerService) HasPreviousPage(currentPage int) bool {
	return currentPage > 1
}

// GetWordCount returns the total number of words
func (ps *PagerService) GetWordCount() int {
	return len(ps.wordList.Words)
}

// GetPageRange returns the start and end word indices for a page
func (ps *PagerService) GetPageRange(pageNumber int) (start, end int, err error) {
	totalPages := ps.GetTotalPages()

	if pageNumber < 1 || pageNumber > totalPages {
		return 0, 0, fmt.Errorf("page number %d is out of range (1-%d)", pageNumber, totalPages)
	}

	start = (pageNumber - 1) * ps.pageSize + 1
	end = pageNumber * ps.pageSize

	if end > ps.GetWordCount() {
		end = ps.GetWordCount()
	}

	return start, end, nil
}

// UpdatePageSize updates the page size and returns the first page with new size
func (ps *PagerService) UpdatePageSize(newPageSize int) (*dao.Page, error) {
	if newPageSize <= 0 {
		return nil, fmt.Errorf("page size must be positive")
	}

	ps.pageSize = newPageSize
	log.Info().Int("newPageSize", newPageSize).Msg("Updated page size")

	return ps.GetFirstPage()
}
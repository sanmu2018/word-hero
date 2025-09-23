package service

import (
	"fmt"
	"math"

	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/log"
	"github.com/sanmu2018/word-hero/pkg/pke"
)

// PagerService handles pagination logic for vocabulary words
type PagerService struct {
	vocabularyService *VocabularyService
}

// NewPagerService creates a new pager service instance
func NewPagerService() *PagerService {
	log.Info().Msg("Creating pager service")

	return &PagerService{}
}

// SetVocabularyService sets the vocabulary service reference
func (ps *PagerService) SetVocabularyService(vocabularyService *VocabularyService) {
	ps.vocabularyService = vocabularyService
}

// GetTotalPages returns the total number of pages for given page size
func (ps *PagerService) GetTotalPages(pageSize int) int {
	if ps.vocabularyService == nil {
		return 0
	}
	if pageSize <= 0 {
		pageSize = 12 // Default page size
	}
	totalWords := ps.vocabularyService.GetWordCount()
	if totalWords == 0 {
		return 0
	}
	return int(math.Ceil(float64(totalWords) / float64(pageSize)))
}

// GetPage returns a specific page of words with metadata
func (ps *PagerService) GetPage(pageNumber, pageSize int) (*dto.Page, error) {
	if ps.vocabularyService == nil {
		return nil, fmt.Errorf("vocabulary service not initialized")
	}

	baseList := &dao.BaseList{
		PageNum:  pageNumber,
		PageSize: pageSize,
	}

	vocabPage, err := ps.vocabularyService.GetWordsByPage(baseList)
	if err != nil {
		return nil, fmt.Errorf("failed to get page: %w", err)
	}

	return &dto.Page{
		Words:      vocabPage.Words,
		TotalCount: vocabPage.TotalCount,
		PageNumber: vocabPage.PageNumber,
		PageSize:   vocabPage.PageSize,
		TotalPages: vocabPage.TotalPages,
	}, nil
}

// GetPageData returns page data with additional metadata for API responses
func (ps *PagerService) GetPageData(pageNumber, pageSize int) (*pke.BaseListResp, error) {
	page, err := ps.GetPage(pageNumber, pageSize)
	if err != nil {
		return nil, err
	}

	data := &pke.BaseListResp{
		Items: page.Words,
		Total: page.TotalCount,
	}
	return data, nil
}

// GetFirstPage returns the first page
func (ps *PagerService) GetFirstPage(pageSize int) (*dto.Page, error) {
	return ps.GetPage(1, pageSize)
}

// GetLastPage returns the last page
func (ps *PagerService) GetLastPage(pageSize int) (*dto.Page, error) {
	return ps.GetPage(ps.GetTotalPages(pageSize), pageSize)
}

// GetNextPage returns the next page
func (ps *PagerService) GetNextPage(currentPage, pageSize int) (*dto.Page, error) {
	return ps.GetPage(currentPage+1, pageSize)
}

// GetPreviousPage returns the previous page
func (ps *PagerService) GetPreviousPage(currentPage, pageSize int) (*dto.Page, error) {
	if currentPage <= 1 {
		return nil, fmt.Errorf("already on first page")
	}
	return ps.GetPage(currentPage-1, pageSize)
}

// HasNextPage checks if there is a next page
func (ps *PagerService) HasNextPage(currentPage, pageSize int) bool {
	return currentPage < ps.GetTotalPages(pageSize)
}

// HasPreviousPage checks if there is a previous page
func (ps *PagerService) HasPreviousPage(currentPage int) bool {
	return currentPage > 1
}

// GetWordCount returns the total number of words
func (ps *PagerService) GetWordCount() int {
	if ps.vocabularyService == nil {
		return 0
	}
	return ps.vocabularyService.GetWordCount()
}

// GetPageRange returns the start and end word indices for a page
func (ps *PagerService) GetPageRange(pageNumber, pageSize int) (start, end int, err error) {
	if pageSize <= 0 {
		pageSize = 25 // Default page size
	}

	totalPages := ps.GetTotalPages(pageSize)

	if pageNumber < 1 || pageNumber > totalPages {
		return 0, 0, fmt.Errorf("page number %d is out of range (1-%d)", pageNumber, totalPages)
	}

	start = (pageNumber-1)*pageSize + 1
	end = pageNumber * pageSize

	if end > ps.GetWordCount() {
		end = ps.GetWordCount()
	}

	return start, end, nil
}

// UpdatePageSize returns the first page with the specified size
// Note: Since PagerService no longer stores pageSize, this method is simplified
func (ps *PagerService) UpdatePageSize(newPageSize int) (*dto.Page, error) {
	if newPageSize <= 0 {
		return nil, fmt.Errorf("page size must be positive")
	}

	log.Info().Int("newPageSize", newPageSize).Msg("Using new page size")

	return ps.GetFirstPage(newPageSize)
}

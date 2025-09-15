package service

import (
	"fmt"

	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/log"
)

// WordTagService handles word tagging business logic
type WordTagService struct {
	wordTagDAO        *dao.WordTagDAO
	wordDAO           *dao.WordDAO
	userDAO           *dao.UserDAO
	vocabularyService *VocabularyService
}

// NewWordTagService creates a new WordTagService instance
func NewWordTagService(wordTagDAO *dao.WordTagDAO, wordDAO *dao.WordDAO, userDAO *dao.UserDAO, vocabularyService *VocabularyService) *WordTagService {
	log.Info().Msg("Creating word tag service")

	return &WordTagService{
		wordTagDAO:        wordTagDAO,
		wordDAO:           wordDAO,
		userDAO:           userDAO,
		vocabularyService: vocabularyService,
	}
}

// MarkWordAsKnown marks a word as known
func (s *WordTagService) MarkWordAsKnown(req *dto.WordMarkRequest) (*dto.WordMarkResponse, error) {
	// Validate user exists
	user, err := s.userDAO.FindByID(req.UserID)
	if err != nil {
		log.Error(err).Str("user_id", req.UserID).Msg("Failed to find user")
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validate word exists
	word, err := s.wordDAO.GetByID(req.WordID)
	if err != nil {
		log.Error(err).Str("word_id", req.WordID).Msg("Failed to find word")
		return nil, fmt.Errorf("word not found: %w", err)
	}

	// Mark word as known
	if err := s.wordTagDAO.MarkWordAsKnown(req.WordID, req.UserID); err != nil {
		log.Error(err).Str("user_id", req.UserID).Str("word_id", req.WordID).Msg("Failed to mark word as known")
		return nil, fmt.Errorf("failed to mark word as known: %w", err)
	}

	log.Info().
		Str("user_id", req.UserID).
		Str("username", user.Username).
		Str("word_id", req.WordID).
		Str("english", word.English).
		Msg("Word marked as known")

	return &dto.WordMarkResponse{
		WordID:    req.WordID,
		IsMarked:  true,
		MarkCount: 1, // In new design, always 1 if known
		Message:   "单词已标记为认识",
	}, nil
}

// RemoveWordMark removes a word's mark (marks as unknown)
func (s *WordTagService) RemoveWordMark(req *dto.WordMarkRequest) (*dto.WordMarkResponse, error) {
	// Validate user exists
	user, err := s.userDAO.FindByID(req.UserID)
	if err != nil {
		log.Error(err).Str("user_id", req.UserID).Msg("Failed to find user")
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validate word exists
	word, err := s.wordDAO.GetByID(req.WordID)
	if err != nil {
		log.Error(err).Str("word_id", req.WordID).Msg("Failed to find word")
		return nil, fmt.Errorf("word not found: %w", err)
	}

	// Remove word mark
	if err := s.wordTagDAO.RemoveWordMark(req.WordID, req.UserID); err != nil {
		log.Error(err).Str("user_id", req.UserID).Str("word_id", req.WordID).Msg("Failed to remove word mark")
		return nil, fmt.Errorf("failed to remove word mark: %w", err)
	}

	log.Info().
		Str("user_id", req.UserID).
		Str("username", user.Username).
		Str("word_id", req.WordID).
		Str("english", word.English).
		Msg("Word mark removed")

	return &dto.WordMarkResponse{
		WordID:    req.WordID,
		IsMarked:  false,
		MarkCount: 0, // In new design, always 0 if unknown
		Message:   "单词标记已移除",
	}, nil
}

// GetWordMarkStatus checks if a word is marked as known
func (s *WordTagService) GetWordMarkStatus(wordID, userID string) (*dto.WordMarkStatus, error) {
	// Check if user exists
	if _, err := s.userDAO.FindByID(userID); err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to find user")
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if word exists
	if _, err := s.wordDAO.GetByID(wordID); err != nil {
		log.Error(err).Str("word_id", wordID).Msg("Failed to find word")
		return nil, fmt.Errorf("word not found: %w", err)
	}

	// Check mark status
	isMarked, err := s.wordTagDAO.IsWordMarkedAsKnown(wordID, userID)
	if err != nil {
		log.Error(err).Str("word_id", wordID).Str("user_id", userID).Msg("Failed to check word mark status")
		return nil, fmt.Errorf("failed to check word mark status: %w", err)
	}

	var markedAt int64
	if isMarked {
		// Get the mark timestamp
		wordTag, err := s.wordTagDAO.GetByWordIDAndUserID(wordID, userID)
		if err == nil {
			markedAt = wordTag.GetKnownTimestamp()
		}
	}

	return &dto.WordMarkStatus{
		WordID:   wordID,
		IsMarked: isMarked,
		MarkCount: func() int {
			if isMarked {
				return 1
			} else {
				return 0
			}
		}(),
		MarkedAt: markedAt,
	}, nil
}

// GetBatchWordMarkStatus gets mark status for multiple words
func (s *WordTagService) GetBatchWordMarkStatus(userID string, wordIDs []string) (*dto.WordMarkStatusResponse, error) {
	// Check if user exists
	if _, err := s.userDAO.FindByID(userID); err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to find user")
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get mark status for all words
	wordMarkStatuses := make([]dto.WordMarkStatus, 0, len(wordIDs))

	for _, wordID := range wordIDs {
		// Check if word exists
		if _, err := s.wordDAO.GetByID(wordID); err != nil {
			log.Warn().Str("word_id", wordID).Msg("Word not found, skipping")
			continue
		}

		// Check mark status
		isMarked, err := s.wordTagDAO.IsWordMarkedAsKnown(wordID, userID)
		if err != nil {
			log.Error(err).Str("word_id", wordID).Str("user_id", userID).Msg("Failed to check word mark status")
			continue
		}

		var markedAt int64
		if isMarked {
			// Get the mark timestamp
			wordTag, err := s.wordTagDAO.GetByWordIDAndUserID(wordID, userID)
			if err == nil {
				markedAt = wordTag.GetKnownTimestamp()
			}
		}

		wordMarkStatuses = append(wordMarkStatuses, dto.WordMarkStatus{
			WordID:   wordID,
			IsMarked: isMarked,
			MarkCount: func() int {
				if isMarked {
					return 1
				} else {
					return 0
				}
			}(),
			MarkedAt: markedAt,
		})
	}

	return &dto.WordMarkStatusResponse{
		WordMarkStatuses: wordMarkStatuses,
	}, nil
}

// GetUserProgress returns user's learning progress
func (s *WordTagService) GetUserProgress(userID string) (*dto.UserProgressResponse, error) {
	// Validate user exists
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to find user")
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get known words count
	knownWords, err := s.wordTagDAO.GetKnownWordsCount(userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to get known words count")
		return nil, fmt.Errorf("failed to get known words count: %w", err)
	}

	// Get total words count
	totalWords, err := s.wordDAO.GetWordCount()
	if err != nil {
		log.Error(err).Msg("Failed to get total words count")
		return nil, fmt.Errorf("failed to get total words count: %w", err)
	}

	// Calculate progress rate
	var progressRate float64
	if totalWords > 0 {
		progressRate = float64(knownWords) / float64(totalWords) * 100
	}

	// Get recent activity (for simplicity, using known words count as recent activity)
	recentActivity := int(knownWords)

	log.Info().
		Str("user_id", userID).
		Str("username", user.Username).
		Int64("known_words", knownWords).
		Int64("total_words", totalWords).
		Float64("progress_rate", progressRate).
		Msg("Retrieved user progress")

	return &dto.UserProgressResponse{
		UserID:         userID,
		KnownWords:     knownWords,
		TotalWords:     totalWords,
		ProgressRate:   progressRate,
		RecentActivity: recentActivity,
	}, nil
}

// GetKnownWords returns paginated known words
func (s *WordTagService) GetKnownWords(req *dto.KnownWordsRequest) (*dto.KnownWordsResponse, error) {
	// Validate user exists
	if _, err := s.userDAO.FindByID(req.UserID); err != nil {
		log.Error(err).Str("user_id", req.UserID).Msg("Failed to find user")
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get known words
	wordIDs, totalCount, err := s.wordTagDAO.GetKnownWords(req.UserID, nil)
	if err != nil {
		log.Error(err).Str("user_id", req.UserID).Msg("Failed to get known words")
		return nil, fmt.Errorf("failed to get known words: %w", err)
	}

	return &dto.KnownWordsResponse{
		WordIDs:    wordIDs,
		TotalCount: totalCount,
	}, nil
}

// ForgetWords forgets specific words by their IDs
func (s *WordTagService) ForgetWords(userID string, req *dto.ForgetWordsRequest) (*dto.ForgetWordsResponse, error) {

	// Forget all specified words
	forgottenCount, err := s.wordTagDAO.BulkRemoveWordMarks(req.WordIDs, userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Int("word_count", len(req.WordIDs)).Msg("Failed to forget words")
		return nil, fmt.Errorf("failed to forget words: %w", err)
	}

	return &dto.ForgetWordsResponse{
		WordIDs:        req.WordIDs,
		ForgottenCount: forgottenCount,
		Message:        fmt.Sprintf("已忘光 %d 个已认识单词", forgottenCount),
	}, nil
}

// ForgetAllWords forgets all known words
func (s *WordTagService) ForgetAllWords(userID string, req *dto.ForgetAllRequest) (*dto.ForgetAllResponse, error) {

	// Check confirmation
	if !req.Confirm {
		return nil, fmt.Errorf("confirmation required to forget all words")
	}

	// Remove all word marks
	forgottenCount, err := s.wordTagDAO.RemoveAllWordMarks(userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to forget all words")
		return nil, fmt.Errorf("failed to forget all words: %w", err)
	}

	return &dto.ForgetAllResponse{
		ForgottenCount: forgottenCount,
		Message:        fmt.Sprintf("已忘光全部 %d 个已认识单词", forgottenCount),
	}, nil
}

// GetWordTagStats returns statistics for word tags
func (s *WordTagService) GetWordTagStats() (*dto.WordTagStats, error) {
	stats, err := s.wordTagDAO.GetStats()
	if err != nil {
		log.Error(err).Msg("Failed to get word tag stats")
		return nil, fmt.Errorf("failed to get word tag stats: %w", err)
	}

	// Convert to DTO format
	wordTagStats := &dto.WordTagStats{
		TotalWordTags:  stats["total_word_tags"].(int64),
		TotalUserMarks: stats["total_user_marks"].(int64),
	}

	if recentWords, ok := stats["recent_words"].([]table.WordTag); ok {
		wordTagStats.TopWords = recentWords
	}

	return wordTagStats, nil
}

// GetUserWordStats returns detailed word statistics for a user
func (s *WordTagService) GetUserWordStats(userID string) (*dto.UserWordStats, error) {
	// Validate user exists
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to find user")
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get basic counts
	knownWords, err := s.wordTagDAO.GetKnownWordsCount(userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to get known words count")
		return nil, fmt.Errorf("failed to get known words count: %w", err)
	}

	totalWords, err := s.wordDAO.GetWordCount()
	if err != nil {
		log.Error(err).Msg("Failed to get total words count")
		return nil, fmt.Errorf("failed to get total words count: %w", err)
	}

	// Calculate progress rate
	var progressRate float64
	if totalWords > 0 {
		progressRate = float64(knownWords) / float64(totalWords) * 100
	}

	// Get recent marks (last 10)
	baseList := &dao.BaseList{PageNum: 1, PageSize: 10}
	wordIDs, _, err := s.wordTagDAO.GetKnownWords(userID, baseList)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get recent marks")
	}

	// Get word tags for recent words
	var recentMarks []dto.KnownWordInfo
	for _, wordID := range wordIDs {
		wordTag, err := s.wordTagDAO.GetByWordIDAndUserID(wordID, userID)
		if err == nil && wordTag.IsKnown() {
			recentMarks = append(recentMarks, dto.KnownWordInfo{
				WordID:  wordID,
				KnownAt: wordTag.GetKnownTimestamp(),
			})
		}
	}

	log.Info().
		Str("user_id", userID).
		Str("username", user.Username).
		Int64("known_words", knownWords).
		Float64("progress_rate", progressRate).
		Msg("Retrieved user word statistics")

	return &dto.UserWordStats{
		UserID:           userID,
		KnownWordsCount:  knownWords,
		TotalWordsCount:  totalWords,
		ProgressRate:     progressRate,
		RecentMarks:      recentMarks,
		KnownWordsByDate: make(map[string]int),
		TopCategories:    make(map[string]int),
	}, nil
}

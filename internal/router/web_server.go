package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/service"
	"github.com/sanmu2018/word-hero/log"
)

// WebServer handles the web application with layered architecture
type WebServer struct {
	vocabularyService *service.VocabularyService
	pagerService      *service.PagerService
	templateDir       string
	mu                sync.RWMutex
}

// NewWebServer creates a new web server instance
func NewWebServer(vocabularyService *service.VocabularyService, pagerService *service.PagerService, templateDir string) *WebServer {
	log.Info().Str("templateDir", templateDir).Msg("Creating web server")

	// Set up service dependencies
	vocabularyService.SetPagerService(pagerService)

	return &WebServer{
		vocabularyService: vocabularyService,
		pagerService:      pagerService,
		templateDir:       templateDir,
	}
}

// PageData represents the data passed to templates
type PageData struct {
	Words       []dao.Word `json:"words"`
	CurrentPage int       `json:"current_page"`
	TotalPages  int       `json:"total_pages"`
	TotalWords  int       `json:"total_words"`
	PageSize    int       `json:"page_size"`
	HasPrev     bool      `json:"has_prev"`
	HasNext     bool      `json:"has_next"`
	PrevPage    int       `json:"prev_page"`
	NextPage    int       `json:"next_page"`
	StartIndex  int       `json:"start_index"`
	EndIndex    int       `json:"end_index"`
}

// APIResponse represents the JSON response for API calls
type APIResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

// Start starts the web server
func (ws *WebServer) Start(port int) error {
	// Parse templates
	templates := template.Must(template.ParseGlob(ws.templateDir + "/*.html"))

	// Set up static file serving
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Set up routes with middleware
	http.HandleFunc("/", ws.loggingMiddleware(ws.homeHandler(templates)))
	http.HandleFunc("/api/words", ws.loggingMiddleware(ws.apiWordsHandler))
	http.HandleFunc("/api/page/", ws.loggingMiddleware(ws.apiPageHandler))
	http.HandleFunc("/api/search", ws.loggingMiddleware(ws.apiSearchHandler))
	http.HandleFunc("/api/stats", ws.loggingMiddleware(ws.apiStatsHandler))

	log.Info().Int("port", port).Msg("Starting web server")
	log.Info().Str("url", fmt.Sprintf("http://localhost:%d", port)).Msg("Open in browser")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// loggingMiddleware logs HTTP requests with performance metrics
func (ws *WebServer) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next(w, r)

		// Log the request
		duration := time.Since(start)
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Str("ip", r.RemoteAddr).
			Dur("duration", duration).
			Msg("HTTP request")
	}
}

// homeHandler handles the main page
func (ws *WebServer) homeHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Prepare minimal template data - no word data
		data := PageData{
			Words:       []dao.Word{}, // Empty array for initial load
			CurrentPage: 1,
			TotalPages:  ws.pagerService.GetTotalPages(),
			TotalWords:  ws.pagerService.GetWordCount(),
			PageSize:    ws.pagerService.GetPageSize(),
			HasPrev:     false,
			HasNext:     ws.pagerService.HasNextPage(1),
			PrevPage:    1,
			NextPage:    2,
			StartIndex:  1,
			EndIndex:    ws.pagerService.GetPageSize(),
		}

		// Render template without word data
		err := templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			log.Error(err).Msg("Failed to execute template")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// apiWordsHandler handles API requests for words with pagination
func (ws *WebServer) apiWordsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get pagination parameters
	page := 1
	pageSize := 24 // Default page size

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	log.Debug().Int("page", page).Int("pageSize", pageSize).Msg("API words request")

	// Get page data using service layer
	responseData, err := ws.pagerService.GetPageData(page)
	if err != nil {
		log.Error(err).Int("page", page).Msg("Failed to get page data for API")
		json.NewEncoder(w).Encode(APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	// Always update page size in pager service to ensure consistency
	_, err = ws.pagerService.UpdatePageSize(pageSize)
	if err != nil {
		log.Error(err).Int("pageSize", pageSize).Msg("Failed to update page size")
		json.NewEncoder(w).Encode(APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}
	// Get updated page data with new page size
	responseData, err = ws.pagerService.GetPageData(page)
	if err != nil {
		log.Error(err).Int("page", page).Msg("Failed to get updated page data")
		json.NewEncoder(w).Encode(APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{
		Code: 0,
		Data: responseData,
	})
}

// apiPageHandler handles API requests for specific pages
func (ws *WebServer) apiPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract page number from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		json.NewEncoder(w).Encode(APIResponse{
			Code: 150321309,
			Msg:  "Invalid page number",
		})
		return
	}

	pageNum, err := strconv.Atoi(pathParts[2])
	if err != nil || pageNum < 1 {
		json.NewEncoder(w).Encode(APIResponse{
			Code: 150321309,
			Msg:  "Invalid page number",
		})
		return
	}

	// Get page data using service layer
	responseData, err := ws.pagerService.GetPageData(pageNum)
	if err != nil {
		log.Error(err).Int("page", pageNum).Msg("Failed to get page data for API")
		json.NewEncoder(w).Encode(APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{
		Code: 0,
		Data: responseData,
	})
}

// apiSearchHandler handles search requests using service layer
func (ws *WebServer) apiSearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		log.Warn().Msg("Search query is empty")
		json.NewEncoder(w).Encode(APIResponse{
			Code: 150321309,
			Msg:  "Search query is required",
		})
		return
	}

	log.Debug().Str("query", query).Msg("Search request")

	// Use service layer for search
	results, err := ws.vocabularyService.SearchWords(query)
	if err != nil {
		log.Error(err).Str("query", query).Msg("Search failed")
		json.NewEncoder(w).Encode(APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	log.Debug().Str("query", query).Int("results", len(results)).Msg("Search completed")

	json.NewEncoder(w).Encode(APIResponse{
		Code: 0,
		Data: map[string]interface{}{
			"query":   query,
			"results": results,
			"count":   len(results),
		},
	})
}

// apiStatsHandler handles statistics requests using service layer
func (ws *WebServer) apiStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get stats from service layer
	stats := ws.vocabularyService.GetStats()

	// Add pagination stats
	stats["totalPages"] = ws.pagerService.GetTotalPages()
	stats["pageSize"] = ws.pagerService.GetPageSize()

	json.NewEncoder(w).Encode(APIResponse{
		Code: 0,
		Data: stats,
	})
}
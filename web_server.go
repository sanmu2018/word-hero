package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// WebServer handles the web application
type WebServer struct {
	wordList    *WordList
	pager       *Pager
	templateDir string
	mu          sync.RWMutex
}

// NewWebServer creates a new web server instance
func NewWebServer(wordList *WordList, pager *Pager, templateDir string) *WebServer {
	return &WebServer{
		wordList:    wordList,
		pager:       pager,
		templateDir: templateDir,
	}
}

// PageData represents the data passed to templates
type PageData struct {
	Words       []Word
	CurrentPage int
	TotalPages  int
	TotalWords  int
	PageSize    int
	HasPrev     bool
	HasNext     bool
	PrevPage    int
	NextPage    int
	StartIndex  int
	EndIndex    int
}

// APIResponse represents the JSON response for API calls
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Start starts the web server
func (ws *WebServer) Start(port int) error {
	// Parse templates
	templates := template.Must(template.ParseGlob(ws.templateDir + "/*.html"))

	// Set up static file serving
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Set up routes with logging middleware
	http.HandleFunc("/", ws.loggingMiddleware(ws.homeHandler(templates)))
	http.HandleFunc("/api/words", ws.loggingMiddleware(ws.apiWordsHandler))
	http.HandleFunc("/api/page/", ws.loggingMiddleware(ws.apiPageHandler))
	http.HandleFunc("/api/search", ws.loggingMiddleware(ws.apiSearchHandler))
	http.HandleFunc("/api/stats", ws.loggingMiddleware(ws.apiStatsHandler))

	log.Info().Int("port", port).Msg("Starting web server")
	log.Info().Str("url", fmt.Sprintf("http://localhost:%d", port)).Msg("Open in browser")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// loggingMiddleware logs HTTP requests
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
		// Get page number from query parameter
		pageNum := 1
		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				pageNum = p
			}
		}

		// Get page size from query parameter or cookie
		pageSize := 24 // Default page size
		if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
				pageSize = ps
			}
		}

		// Create temporary pager with requested page size
		tempPager := NewPager(ws.wordList, pageSize)

		// Get page data
		page, err := tempPager.GetPage(pageNum)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Calculate indices
		startIndex, endIndex, _ := tempPager.GetPageRange(pageNum)

		// Prepare template data
		data := PageData{
			Words:       page.Words,
			CurrentPage: page.Number,
			TotalPages:  page.TotalPages,
			TotalWords:  tempPager.GetWordCount(),
			PageSize:    page.PageSize,
			HasPrev:     tempPager.HasPreviousPage(pageNum),
			HasNext:     tempPager.HasNextPage(pageNum),
			PrevPage:    pageNum - 1,
			NextPage:    pageNum + 1,
			StartIndex:  startIndex,
			EndIndex:    endIndex,
		}

		// Render template
		err = templates.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// apiWordsHandler handles API requests for words
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

	// Create temporary pager with requested page size
	tempPager := NewPager(ws.wordList, pageSize)
	pageData, err := tempPager.GetPage(page)
	if err != nil {
		log.Warn().Err(err).Int("page", page).Msg("Failed to get page data")
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	startIndex, endIndex, _ := tempPager.GetPageRange(page)

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"words":       pageData.Words,
			"currentPage": pageData.Number,
			"totalPages":  pageData.TotalPages,
			"totalWords":  tempPager.GetWordCount(),
			"pageSize":    pageData.PageSize,
			"startIndex":  startIndex,
			"endIndex":    endIndex,
			"hasPrev":     tempPager.HasPreviousPage(page),
			"hasNext":     tempPager.HasNextPage(page),
			"prevPage":    page - 1,
			"nextPage":    page + 1,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// apiPageHandler handles API requests for specific pages
func (ws *WebServer) apiPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract page number from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "Invalid page number",
		})
		return
	}

	pageNum, err := strconv.Atoi(pathParts[2])
	if err != nil || pageNum < 1 {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "Invalid page number",
		})
		return
	}

	page, err := ws.pager.GetPage(pageNum)
	if err != nil {
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	startIndex, endIndex, _ := ws.pager.GetPageRange(pageNum)

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"words":       page.Words,
			"currentPage": page.Number,
			"totalPages":  page.TotalPages,
			"totalWords":  ws.pager.GetWordCount(),
			"pageSize":    page.PageSize,
			"startIndex":  startIndex,
			"endIndex":    endIndex,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// apiSearchHandler handles search requests
func (ws *WebServer) apiSearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		log.Warn().Msg("Search query is empty")
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   "Search query is required",
		})
		return
	}

	log.Debug().Str("query", query).Msg("Search request")

	var results []Word
	ws.mu.RLock()
	for _, word := range ws.wordList.Words {
		if strings.Contains(strings.ToLower(word.English), query) ||
			strings.Contains(strings.ToLower(word.Chinese), query) {
			results = append(results, word)
		}
	}
	ws.mu.RUnlock()

	log.Debug().Str("query", query).Int("results", len(results)).Msg("Search completed")

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"query":   query,
			"results": results,
			"count":   len(results),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// apiStatsHandler handles statistics requests
func (ws *WebServer) apiStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"totalWords": ws.pager.GetWordCount(),
			"totalPages": ws.pager.GetTotalPages(),
			"pageSize":   ws.pager.pageSize,
			"fileSource": "words/IELTS.xlsx",
		},
	}

	json.NewEncoder(w).Encode(response)
}

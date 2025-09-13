package router

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/service"
	"github.com/sanmu2018/word-hero/log"
)

// WebServer handles the web application with layered architecture
type WebServer struct {
	vocabularyService *service.VocabularyService
	pagerService      *service.PagerService
	templateDir       string
	engine            *gin.Engine
}

// NewWebServer creates a new web server instance
func NewWebServer(vocabularyService *service.VocabularyService, pagerService *service.PagerService, templateDir string) *WebServer {
	log.Info().Str("templateDir", templateDir).Msg("Creating web server")

	// Set up service dependencies
	vocabularyService.SetPagerService(pagerService)

	// Create Gin engine
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Set up routes
	ws := &WebServer{
		vocabularyService: vocabularyService,
		pagerService:      pagerService,
		templateDir:       templateDir,
		engine:            engine,
	}

	// Add middleware
	engine.Use(gin.Logger(), gin.Recovery(), ws.loggingMiddleware())

	// Set up templates
	templates := template.Must(template.ParseGlob(templateDir + "/*.html"))
	engine.SetHTMLTemplate(templates)

	// Set up static files
	engine.Static("/static", "web/static")

	// Register routes
	ws.setupRoutes()

	return ws
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

// setupRoutes configures all routes for the Gin engine
func (ws *WebServer) setupRoutes() {
	// Web routes
	ws.engine.GET("/", ws.homeHandler)

	// API routes
	api := ws.engine.Group("/api")
	{
		api.GET("/words", ws.apiWordsHandler)
		api.GET("/page/:pageNumber", ws.apiPageHandler)
		api.GET("/search", ws.apiSearchHandler)
		api.GET("/stats", ws.apiStatsHandler)
	}
}

// Start starts the web server
func (ws *WebServer) Start(port int) error {
	log.Info().Int("port", port).Msg("Starting web server")
	log.Info().Str("url", fmt.Sprintf("http://localhost:%d", port)).Msg("Open in browser")
	return ws.engine.Run(fmt.Sprintf(":%d", port))
}

// loggingMiddleware logs HTTP requests with performance metrics
func (ws *WebServer) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log the request
		duration := time.Since(start)
		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("ip", c.ClientIP()).
			Int("status", c.Writer.Status()).
			Dur("duration", duration).
			Msg("HTTP request")
	}
}

// homeHandler handles the main page
func (ws *WebServer) homeHandler(c *gin.Context) {
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
	c.HTML(http.StatusOK, "index.html", data)
}

// apiWordsHandler handles API requests for words with pagination
func (ws *WebServer) apiWordsHandler(c *gin.Context) {
	// Get pagination parameters
	page := 1
	pageSize := 24 // Default page size

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	log.Debug().Int("page", page).Int("pageSize", pageSize).Msg("API words request")

	// Get page data using service layer
	responseData, err := ws.pagerService.GetPageData(page)
	if err != nil {
		log.Error(err).Int("page", page).Msg("Failed to get page data for API")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	// Always update page size in pager service to ensure consistency
	_, err = ws.pagerService.UpdatePageSize(pageSize)
	if err != nil {
		log.Error(err).Int("pageSize", pageSize).Msg("Failed to update page size")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}
	// Get updated page data with new page size
	responseData, err = ws.pagerService.GetPageData(page)
	if err != nil {
		log.Error(err).Int("page", page).Msg("Failed to get updated page data")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: responseData,
	})
}

// apiPageHandler handles API requests for specific pages
func (ws *WebServer) apiPageHandler(c *gin.Context) {
	// Extract page number from URL parameter
	pageNum, err := strconv.Atoi(c.Param("pageNumber"))
	if err != nil || pageNum < 1 {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid page number",
		})
		return
	}

	// Get page data using service layer
	responseData, err := ws.pagerService.GetPageData(pageNum)
	if err != nil {
		log.Error(err).Int("page", pageNum).Msg("Failed to get page data for API")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: responseData,
	})
}

// apiSearchHandler handles search requests using service layer
func (ws *WebServer) apiSearchHandler(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		log.Warn().Msg("Search query is empty")
		c.JSON(http.StatusOK, APIResponse{
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
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	log.Debug().Str("query", query).Int("results", len(results)).Msg("Search completed")

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: map[string]interface{}{
			"query":   query,
			"results": results,
			"count":   len(results),
		},
	})
}

// apiStatsHandler handles statistics requests using service layer
func (ws *WebServer) apiStatsHandler(c *gin.Context) {
	// Get stats from service layer
	stats := ws.vocabularyService.GetStats()

	// Add pagination stats
	stats["totalPages"] = ws.pagerService.GetTotalPages()
	stats["pageSize"] = ws.pagerService.GetPageSize()

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: stats,
	})
}
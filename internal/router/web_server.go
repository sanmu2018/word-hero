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
	"github.com/sanmu2018/word-hero/internal/dto"
	"github.com/sanmu2018/word-hero/internal/middleware"
	"github.com/sanmu2018/word-hero/internal/models"
	"github.com/sanmu2018/word-hero/internal/service"
	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/log"
)

// WebServer handles the web application with layered architecture
type WebServer struct {
	vocabularyService *service.VocabularyService
	pagerService      *service.PagerService
	authService       *service.AuthService
	userService       *service.UserService
	wordTagService    *service.WordTagService
	authMiddleware    *middleware.AuthMiddleware
	templateDir       string
	engine            *gin.Engine
}

// NewWebServer creates a new web server instance
func NewWebServer(vocabularyService *service.VocabularyService, pagerService *service.PagerService, authService *service.AuthService, userService *service.UserService, wordTagService *service.WordTagService, authMiddleware *middleware.AuthMiddleware, templateDir string) *WebServer {
	log.Info().Str("templateDir", templateDir).Msg("Creating web server")

	// Create Gin engine
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Set up routes
	ws := &WebServer{
		vocabularyService: vocabularyService,
		pagerService:      pagerService,
		authService:       authService,
		userService:       userService,
		wordTagService:    wordTagService,
		authMiddleware:    authMiddleware,
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
	Words       []table.Word `json:"words"`
	CurrentPage int          `json:"currentPage"`
	TotalPages  int          `json:"totalPages"`
	TotalWords  int          `json:"totalWords"`
	PageSize    int          `json:"pageSize"`
	HasPrev     bool         `json:"hasPrev"`
	HasNext     bool         `json:"hasNext"`
	PrevPage    int          `json:"prevPage"`
	NextPage    int          `json:"nextPage"`
	StartIndex  int          `json:"startIndex"`
	EndIndex    int          `json:"endIndex"`
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
		// Public vocabulary endpoints
		api.GET("/words", ws.apiWordsHandler)
		api.GET("/search", ws.apiSearchHandler)
		api.GET("/stats", ws.apiStatsHandler)

		// Authentication endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", ws.apiRegisterHandler)
			auth.POST("/login", ws.apiLoginHandler)
			auth.POST("/logout", ws.apiLogoutHandler)
			auth.GET("/me", ws.authMiddleware.RequireAuth(), ws.apiGetCurrentUserHandler)
			auth.PUT("/profile", ws.authMiddleware.RequireAuth(), ws.apiUpdateProfileHandler)
			auth.POST("/change-password", ws.authMiddleware.RequireAuth(), ws.apiChangePasswordHandler)
		}

		// Protected user endpoints
		user := api.Group("/user")
		user.Use(ws.authMiddleware.RequireAuth())
		{
			user.GET("/profile", ws.apiGetUserProfileHandler)
		}

		// Word tag endpoints
		wordTags := api.Group("/word-tags")
		// All word tag operations require authentication
		wordTags.Use(ws.authMiddleware.RequireAuth())
		{
			wordTags.POST("/mark", ws.apiMarkWordHandler)
			wordTags.DELETE("/unmark", ws.apiUnmarkWordHandler)
			wordTags.GET("/status/:wordId", ws.apiGetWordMarkStatusHandler)
			wordTags.GET("/known", ws.apiGetKnownWordsHandler)
			wordTags.GET("/progress", ws.apiGetUserProgressHandler)
			wordTags.GET("/stats", ws.apiGetWordTagStatsHandler)
			wordTags.POST("/forget-words", ws.apiForgetWordsHandler)
			wordTags.POST("/forget-all", ws.apiForgetAllHandler)
		}
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
	defaultPageSize := 12
	data := PageData{
		Words:       []table.Word{}, // Empty array for initial load
		CurrentPage: 1,
		TotalPages:  ws.pagerService.GetTotalPages(defaultPageSize),
		TotalWords:  ws.pagerService.GetWordCount(),
		PageSize:    defaultPageSize,
		HasPrev:     false,
		HasNext:     ws.pagerService.HasNextPage(1, defaultPageSize),
		PrevPage:    1,
		NextPage:    2,
		StartIndex:  1,
		EndIndex:    defaultPageSize,
	}

	// Render template without word data
	c.HTML(http.StatusOK, "index.html", data)
}

// apiWordsHandler handles API requests for words with pagination
func (ws *WebServer) apiWordsHandler(c *gin.Context) {
	// Get pagination parameters
	page := 1
	pageSize := 12 // Default page size

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
	responseData, err := ws.pagerService.GetPageData(page, pageSize)
	if err != nil {
		log.Error(err).Int("page", page).Msg("Failed to get page data for API")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	// Return page data with the requested page size
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
	defaultPageSize := 12
	stats["totalPages"] = ws.pagerService.GetTotalPages(defaultPageSize)
	stats["pageSize"] = defaultPageSize

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: stats,
	})
}

// apiRegisterHandler handles user registration
func (ws *WebServer) apiRegisterHandler(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid request format",
		})
		return
	}

	log.Debug().Str("username", req.Username).Str("email", req.Email).Msg("Registration request")

	user, token, err := ws.authService.Register(&req)
	if err != nil {
		log.Error(err).Str("username", req.Username).Msg("Registration failed")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	log.Info().Str("user_id", user.ID).Str("username", user.Username).Msg("User registered successfully")

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: map[string]interface{}{
			"user":  models.NewUserBusiness(user).ToResponse(),
			"token": token,
		},
		Msg: "Registration successful",
	})
}

// apiLoginHandler handles user login
func (ws *WebServer) apiLoginHandler(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid request format",
		})
		return
	}

	log.Debug().Str("username", req.Username).Msg("Login request")

	user, token, err := ws.authService.Login(&req)
	if err != nil {
		log.Error(err).Str("username", req.Username).Msg("Login failed")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	log.Info().Str("user_id", user.ID).Str("username", user.Username).Msg("User logged in successfully")

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: map[string]interface{}{
			"user":  models.NewUserBusiness(user).ToResponse(),
			"token": token,
		},
		Msg: "Login successful",
	})
}

// apiLogoutHandler handles user logout
func (ws *WebServer) apiLogoutHandler(c *gin.Context) {
	// For JWT, logout is typically handled client-side by token removal
	// Server-side logout could involve token blacklisting if needed
	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Msg:  "Logout successful",
	})
}

// apiGetCurrentUserHandler gets the current authenticated user
func (ws *WebServer) apiGetCurrentUserHandler(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: models.NewUserBusiness(user).ToResponse(),
	})
}

// apiUpdateProfileHandler updates user profile
func (ws *WebServer) apiUpdateProfileHandler(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "User not authenticated",
		})
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid request format",
		})
		return
	}

	log.Debug().Str("user_id", userID).Msg("Profile update request")

	user, err := ws.userService.UpdateUserProfile(userID, &req)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Profile update failed")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	log.Info().Str("user_id", userID).Msg("User profile updated successfully")

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: user,
		Msg:  "Profile updated successfully",
	})
}

// apiChangePasswordHandler handles password change
func (ws *WebServer) apiChangePasswordHandler(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "User not authenticated",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid request format",
		})
		return
	}

	log.Debug().Str("user_id", userID).Msg("Password change request")

	if err := ws.authService.ChangePassword(userID, &req); err != nil {
		log.Error(err).Str("user_id", userID).Msg("Password change failed")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	log.Info().Str("user_id", userID).Msg("User password changed successfully")

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Msg:  "Password changed successfully",
	})
}

// apiGetUserProfileHandler gets user profile
func (ws *WebServer) apiGetUserProfileHandler(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "User not authenticated",
		})
		return
	}

	log.Debug().Str("user_id", userID).Msg("Profile request")

	user, err := ws.userService.GetUserProfile(userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to get user profile")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: user,
	})
}

// apiMarkWordHandler marks a word as known by a user
func (ws *WebServer) apiMarkWordHandler(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		log.Error(err).Msg("User not authenticated")
		c.JSON(http.StatusOK, APIResponse{
			Code: 401,
			Msg:  "请先登录后再进行操作",
		})
		return
	}

	var req dto.WordMarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid request format",
		})
		return
	}

	// Set user ID from context
	req.UserID = userID

	response, err := ws.wordTagService.MarkWordAsKnown(&req)
	if err != nil {
		log.Error(err).Str("user_id", userID).Str("word_id", req.WordID).Msg("Failed to mark word as known")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: response,
		Msg:  response.Message,
	})
}

// apiUnmarkWordHandler removes a user's mark from a word
func (ws *WebServer) apiUnmarkWordHandler(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		log.Error(err).Msg("User not authenticated")
		c.JSON(http.StatusOK, APIResponse{
			Code: 401,
			Msg:  "请先登录后再进行操作",
		})
		return
	}

	var req dto.WordMarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid request format",
		})
		return
	}

	// Set user ID from context
	req.UserID = userID

	response, err := ws.wordTagService.RemoveWordMark(&req)
	if err != nil {
		log.Error(err).Str("user_id", userID).Str("word_id", req.WordID).Msg("Failed to remove word mark")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: response,
		Msg:  response.Message,
	})
}

// apiGetWordMarkStatusHandler gets the mark status of a word for a user
func (ws *WebServer) apiGetWordMarkStatusHandler(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		log.Error(err).Msg("User not authenticated")
		c.JSON(http.StatusOK, APIResponse{
			Code: 401,
			Msg:  "请先登录后再进行操作",
		})
		return
	}

	wordID := c.Param("wordId")
	if wordID == "" {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Word ID is required",
		})
		return
	}

	status, err := ws.wordTagService.GetWordMarkStatus(wordID, userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Str("word_id", wordID).Msg("Failed to get word mark status")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: status,
	})
}

// apiGetKnownWordsHandler gets known words for a user
func (ws *WebServer) apiGetKnownWordsHandler(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		log.Error(err).Msg("User not authenticated")
		c.JSON(http.StatusOK, APIResponse{
			Code: 401,
			Msg:  "请先登录后再进行操作",
		})
		return
	}

	baseList := &dao.BaseList{PageNum: 1, PageSize: 1000}
	response, err := ws.vocabularyService.GetKnownWordsByUser(userID, baseList)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to get known words")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: response,
	})
}

// apiGetUserProgressHandler gets user's learning progress
func (ws *WebServer) apiGetUserProgressHandler(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		log.Error(err).Msg("User not authenticated")
		c.JSON(http.StatusOK, APIResponse{
			Code: 401,
			Msg:  "请先登录后再进行操作",
		})
		return
	}

	response, err := ws.wordTagService.GetUserProgress(userID)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to get user progress")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: response,
	})
}

// apiGetWordTagStatsHandler gets word tag statistics
func (ws *WebServer) apiGetWordTagStatsHandler(c *gin.Context) {
	// This endpoint doesn't require user authentication for general stats
	// TODO: Re-enable proper authentication after testing if needed

	response, err := ws.wordTagService.GetWordTagStats()
	if err != nil {
		log.Error(err).Msg("Failed to get word tag stats")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: response,
	})
}


// apiForgetWordsHandler handles forgetting specific words
func (ws *WebServer) apiForgetWordsHandler(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		log.Error(err).Msg("User not authenticated")
		c.JSON(http.StatusOK, APIResponse{
			Code: 401,
			Msg:  "请先登录后再进行操作",
		})
		return
	}

	var req dto.ForgetWordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid request format",
		})
		return
	}

	log.Debug().Str("user_id", userID).Int("word_count", len(req.WordIDs)).Msg("Forget words request")

	response, err := ws.wordTagService.ForgetWords(userID, &req)
	if err != nil {
		log.Error(err).Str("user_id", userID).Int("word_count", len(req.WordIDs)).Msg("Failed to forget words")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: response,
		Msg:  response.Message,
	})
}

// apiForgetAllHandler handles forgetting all words
func (ws *WebServer) apiForgetAllHandler(c *gin.Context) {
	// Get user ID from authentication middleware
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		log.Error(err).Msg("User not authenticated")
		c.JSON(http.StatusOK, APIResponse{
			Code: 401,
			Msg:  "请先登录后再进行操作",
		})
		return
	}

	var req dto.ForgetAllRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  "Invalid request format",
		})
		return
	}

	log.Debug().Str("user_id", userID).Bool("confirm", req.Confirm).Msg("Forget all request")

	response, err := ws.wordTagService.ForgetAllWords(userID, &req)
	if err != nil {
		log.Error(err).Str("user_id", userID).Msg("Failed to forget all words")
		c.JSON(http.StatusOK, APIResponse{
			Code: 150321309,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code: 0,
		Data: response,
		Msg:  response.Message,
	})
}

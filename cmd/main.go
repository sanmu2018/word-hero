package main

import (
	"os"
	"strconv"

	"github.com/sanmu2018/word-hero/internal/conf"
	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/middleware"
	"github.com/sanmu2018/word-hero/internal/router"
	"github.com/sanmu2018/word-hero/internal/service"
	"github.com/sanmu2018/word-hero/internal/utils"
	"github.com/sanmu2018/word-hero/log"
)

func main() {
	// Load configuration
	config, err := conf.LoadConfig()
	if err != nil {
		log.Error(err).Msg("Failed to load configuration")
		os.Exit(1)
	}

	port := config.Server.Port

	// Initialize logger
	log.Info().Msg("=== Word Hero Web Application ===")
	log.Info().Msg("Initializing web server...")
	log.Info().Int("port", port).Msg("Using port")

	// Initialize database
	log.Info().Msg("Initializing database connection...")
	if err := dao.InitDatabase(&config.Database); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	// Run database migrations
	log.Info().Msg("Running database migrations...")
	if err := dao.RunMigrations(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize Word DAO
	wordDAO := dao.NewWordDAO()

	// Check if word data is available
	log.Info().Msg("Validating word data availability...")
	if isEmpty, err := wordDAO.IsEmpty(); err != nil {
		log.Error(err).Msg("Failed to check word data availability")
		log.Fatal().Msg("Application cannot start without word data")
	} else if isEmpty {
		log.Error(err).Msg("Word data validation failed - table is empty")
		log.Info().Msg("Please run the migration tool to import word data from Excel file")
		log.Fatal().Msg("Application cannot start without word data")
	}

	// Get word count for statistics
	totalWords, err := wordDAO.GetWordCount()
	if err != nil {
		log.Error(err).Msg("Failed to get word count")
	} else {
		log.Info().Int64("totalWords", totalWords).Msg("Word data validated")
	}

	// Initialize authentication services
	log.Info().Msg("Initializing authentication services...")
	userDAO := dao.NewUserDAO()
	jwtUtils := utils.NewJWTUtils(&config.JWT)
	authService := service.NewAuthService(userDAO, jwtUtils)
	userService := service.NewUserService(userDAO)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize service layer
	pagerService := service.NewPagerService()
	vocabularyService := service.NewVocabularyService(wordDAO)

	// Set service dependencies
	pagerService.SetVocabularyService(vocabularyService)

	// Initialize router layer
	webServer := router.NewWebServer(vocabularyService, pagerService, authService, userService, authMiddleware, "web/templates")

	// Show database info
	log.Info().Str("database", config.Database.DBName).Msg("Database Information:")
	log.Info().Str("host", config.Database.Host).Str("port", strconv.Itoa(config.Database.Port)).Msg("Database connection")
	log.Info().Int64("totalWords", totalWords).Msg("Available vocabulary words")

	log.Info().Msg("Press Ctrl+C to stop the server")

	// Start the web server
	if err := webServer.Start(port); err != nil {
		log.Fatal().Str("err", err.Error()).Msg("Failed to start web server")
	}
}

func Init() error {
	return nil
}

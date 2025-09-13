package main

import (
	"fmt"
	"os"

	"github.com/sanmu2018/word-hero/internal/conf"
	"github.com/sanmu2018/word-hero/internal/dao"
	"github.com/sanmu2018/word-hero/internal/router"
	"github.com/sanmu2018/word-hero/internal/service"
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

	// Check if Excel file exists
	if _, err := os.Stat(config.App.ExcelFile); err != nil {
		log.Error(err).Str("file", config.App.ExcelFile).Msg("Excel file not found")
		os.Exit(1)
	}

	// Initialize data access layer
	excelReader := dao.NewExcelReader(config.App.ExcelFile)

	// Validate file
	if err := excelReader.ValidateFile(); err != nil {
		log.Error(err).Msg("Error validating file")
		os.Exit(1)
	}

	// Read words from Excel file
	log.Info().Msg("Reading vocabulary data...")
	wordList, err := excelReader.ReadWords()
	if err != nil {
		log.Error(err).Msg("Error reading Excel file")
		os.Exit(1)
	}

	log.Info().Int("count", len(wordList.Words)).Msg("Successfully loaded vocabulary words")

	// Initialize service layer
	pagerService := service.NewPagerService(wordList, config.App.PageSize)
	vocabularyService := service.NewVocabularyService(wordList, excelReader)

	// Initialize router layer
	webServer := router.NewWebServer(vocabularyService, pagerService, "web/templates")

	// Show file info
	info, err := excelReader.GetFileInfo()
	if err != nil {
		log.Error(err).Msg("Could not get file info")
	} else {
		log.Info().Msg("File Information:")
		fmt.Println(info) // Keep fmt.Println for formatted file info
	}

	log.Info().Msg("Press Ctrl+C to stop the server")

	// Start the web server
	if err := webServer.Start(port); err != nil {
		log.Fatal().Str("err", err.Error()).Msg("Failed to start web server")
	}
}

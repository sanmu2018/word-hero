package main

import (
	"fmt"
	"os"

	"github.com/sanmu2018/word-hero/log"
)

func main() {
	// Load configuration
	config, err := LoadConfig()
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
		log.Info().Msg("Please ensure the IELTS.xlsx file is in the words/ directory.")
		os.Exit(1)
	}

	// Initialize Excel reader
	reader := NewExcelReader(config.App.ExcelFile)

	// Validate file
	if err := reader.ValidateFile(); err != nil {
		log.Error(err).Msg("Error validating file")
		os.Exit(1)
	}

	// Read words from Excel file
	log.Info().Msg("Reading vocabulary data...")
	wordList, err := reader.ReadWords()
	if err != nil {
		log.Error(err).Msg("Error reading Excel file")
		os.Exit(1)
	}

	log.Info().Int("count", len(wordList.Words)).Msg("Successfully loaded vocabulary words")

	// Initialize pager
	pager := NewPager(wordList, config.App.PageSize)

	// Initialize web server
	webServer := NewWebServer(wordList, pager, "web/templates")

	// Show file info
	info, err := reader.GetFileInfo()
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
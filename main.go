package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sanmu2018/word-hero/log"
)

const (
	PageSize  = 25
	ExcelFile = "words/IELTS.xlsx"
	WebPort   = 8082
)

func main() {
	// Initialize logger
	log.Info().Msg("=== Word Hero Web Application ===")
	log.Info().Msg("Initializing web server...")

	// Check if Excel file exists
	if _, err := os.Stat(ExcelFile); err != nil {
		log.Error(err).Str("file", ExcelFile).Msg("Excel file not found")
		log.Error(errors.New("")).Msg("Please ensure the IELTS.xlsx file is in the words/ directory.")
		os.Exit(1)
	}

	// Initialize Excel reader
	reader := NewExcelReader(ExcelFile)

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
	pager := NewPager(wordList, PageSize)

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
	if err := webServer.Start(WebPort); err != nil {
		log.Fatal().Str("err", err.Error()).Msg("Failed to start web server")
	}
}

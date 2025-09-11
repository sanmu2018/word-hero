package main

import (
	"fmt"
	"log"
	"os"
)

const (
	PageSize = 25
	ExcelFile = "words/IELTS.xlsx"
	WebPort = 8082
)

func main() {
	// Check if Excel file exists
	if _, err := os.Stat(ExcelFile); err != nil {
		fmt.Printf("Error: Excel file not found at %s\n", ExcelFile)
		fmt.Println("Please ensure the IELTS.xlsx file is in the words/ directory.")
		os.Exit(1)
	}

	fmt.Println("=== Word Hero Web Application ===")
	fmt.Println("Initializing web server...")

	// Initialize Excel reader
	reader := NewExcelReader(ExcelFile)
	
	// Validate file
	if err := reader.ValidateFile(); err != nil {
		fmt.Printf("Error validating file: %v\n", err)
		os.Exit(1)
	}

	// Read words from Excel file
	fmt.Println("Reading vocabulary data...")
	wordList, err := reader.ReadWords()
	if err != nil {
		fmt.Printf("Error reading Excel file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully loaded %d vocabulary words\n", len(wordList.Words))

	// Initialize pager
	pager := NewPager(wordList, PageSize)
	
	// Initialize web server
	webServer := NewWebServer(wordList, pager, "web/templates")
	
	// Show file info
	info, err := reader.GetFileInfo()
	if err != nil {
		fmt.Printf("Warning: Could not get file info: %v\n", err)
	} else {
		fmt.Println("\nFile Information:")
		fmt.Println(info)
	}

	fmt.Printf("\nStarting web server on port %d\n", WebPort)
	fmt.Printf("Open http://localhost:%d in your browser\n", WebPort)
	fmt.Println("Press Ctrl+C to stop the server")
	
	// Start the web server
	if err := webServer.Start(WebPort); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
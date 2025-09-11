package main

import (
	"fmt"
	"os"
)

const (
	PageSize = 25
	ExcelFile = "words/IELTS.xlsx"
)

func main() {
	// Check if Excel file exists
	if _, err := os.Stat(ExcelFile); err != nil {
		fmt.Printf("Error: Excel file not found at %s\n", ExcelFile)
		fmt.Println("Please ensure the IELTS.xlsx file is in the words/ directory.")
		os.Exit(1)
	}

	// Initialize Excel reader
	reader := NewExcelReader(ExcelFile)
	
	// Validate file
	if err := reader.ValidateFile(); err != nil {
		fmt.Printf("Error validating file: %v\n", err)
		os.Exit(1)
	}

	// Read words from Excel file
	wordList, err := reader.ReadWords()
	if err != nil {
		fmt.Printf("Error reading Excel file: %v\n", err)
		os.Exit(1)
	}

	// Initialize pager
	pager := NewPager(wordList, PageSize)
	
	// Initialize UI
	ui := NewUI(pager)
	
	// Show welcome message
	ui.ShowWelcome(len(wordList.Words))
	
	// Start with first page
	currentPage := 1
	
	// Main application loop
	for {
		// Get current page
		page, err := pager.GetPage(currentPage)
		if err != nil {
			fmt.Printf("Error getting page: %v\n", err)
			break
		}
		
		// Display the page
		ui.DisplayPage(page)
		
		// Get user input
		input, err := ui.GetUserInput()
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			break
		}
		
		// Handle navigation
		newPage, shouldQuit, err := ui.HandleNavigation(input, currentPage)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			// Wait for user to acknowledge error
			fmt.Println("Press Enter to continue...")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		
		if shouldQuit {
			break
		}
		
		currentPage = newPage
	}
	
	// Show goodbye message
	ui.ShowGoodbye()
}
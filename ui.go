package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// UI handles user interface interactions
type UI struct {
	pager *Pager
}

// NewUI creates a new UI instance
func NewUI(pager *Pager) *UI {
	return &UI{
		pager: pager,
	}
}

// DisplayPage displays a page of words
func (ui *UI) DisplayPage(page *Page) {
	ui.clearScreen()
	
	fmt.Println("=== Word Hero - IELTS Vocabulary ===")
	fmt.Printf("Page %d of %d (Total words: %d)\n", page.Number, page.TotalPages, ui.pager.GetWordCount())
	fmt.Println(strings.Repeat("=", 50))
	
	start, end, _ := ui.pager.GetPageRange(page.Number)
	fmt.Printf("Showing words %d-%d\n\n", start, end)
	
	// Display words
	for i, word := range page.Words {
		fmt.Printf("%2d. %-20s | %s\n", i+1, word.English, word.Chinese)
	}
	
	fmt.Println(strings.Repeat("=", 50))
	ui.displayNavigation(page)
}

// displayNavigation shows navigation options
func (ui *UI) displayNavigation(page *Page) {
	fmt.Println("\nNavigation:")
	fmt.Println("[n] Next page    [p] Previous page    [f] First page    [l] Last page")
	fmt.Println("[g] Go to page  [s] Show stats      [q] Quit")
	
	if ui.pager.HasPreviousPage(page.Number) {
		fmt.Print("You can go to previous page (p) or first page (f). ")
	}
	if ui.pager.HasNextPage(page.Number) {
		fmt.Print("You can go to next page (n) or last page (l). ")
	}
	
	fmt.Print("\n> ")
}

// clearScreen clears the console screen
func (ui *UI) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// GetUserInput gets user input for navigation
func (ui *UI) GetUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.ToLower(input)), nil
}

// ShowStats displays statistics
func (ui *UI) ShowStats() {
	ui.clearScreen()
	fmt.Println("=== Word Hero - Statistics ===")
	fmt.Printf("Total words: %d\n", ui.pager.GetWordCount())
	fmt.Printf("Total pages: %d\n", ui.pager.GetTotalPages())
	fmt.Printf("Words per page: %d\n", ui.pager.pageSize)
	fmt.Println("\nPress Enter to continue...")
	
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

// GoToPage allows user to jump to a specific page
func (ui *UI) GoToPage() (int, error) {
	ui.clearScreen()
	fmt.Println("=== Go to Page ===")
	fmt.Printf("Enter page number (1-%d): ", ui.pager.GetTotalPages())
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	
	pageNum, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, fmt.Errorf("invalid page number")
	}
	
	return pageNum, nil
}

// ShowWelcome displays welcome message
func (ui *UI) ShowWelcome(wordCount int) {
	ui.clearScreen()
	fmt.Println("=== Welcome to Word Hero ===")
	fmt.Printf("Loaded %d IELTS vocabulary words\n", wordCount)
	fmt.Printf("Displaying %d words per page\n", ui.pager.pageSize)
	fmt.Printf("Total pages: %d\n", ui.pager.GetTotalPages())
	fmt.Println("\nPress Enter to start...")
	
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

// ShowGoodbye displays goodbye message
func (ui *UI) ShowGoodbye() {
	ui.clearScreen()
	fmt.Println("=== Thank you for using Word Hero ===")
	fmt.Println("Keep learning and improving your vocabulary!")
}

// HandleNavigation handles user navigation input
func (ui *UI) HandleNavigation(input string, currentPage int) (int, bool, error) {
	switch input {
	case "n":
		if ui.pager.HasNextPage(currentPage) {
			return currentPage + 1, false, nil
		}
		return currentPage, false, fmt.Errorf("already on last page")
	
	case "p":
		if ui.pager.HasPreviousPage(currentPage) {
			return currentPage - 1, false, nil
		}
		return currentPage, false, fmt.Errorf("already on first page")
	
	case "f":
		return 1, false, nil
	
	case "l":
		return ui.pager.GetTotalPages(), false, nil
	
	case "g":
		pageNum, err := ui.GoToPage()
		if err != nil {
			return currentPage, false, err
		}
		return pageNum, false, nil
	
	case "s":
		ui.ShowStats()
		return currentPage, false, nil
	
	case "q":
		return currentPage, true, nil
	
	default:
		return currentPage, false, fmt.Errorf("invalid command")
	}
}
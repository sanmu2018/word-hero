# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Word Hero is a dual-platform Go application for learning IELTS vocabulary, offering both command-line and web interfaces. It reads vocabulary data from `words/IELTS.xlsx` (3673 words) and displays it in a paginated format (25 words per page) with interactive navigation and search functionality.

## Project Architecture

### Core Components
- **models.go**: Defines Word, WordList, and Page data structures
- **excel_reader.go**: Handles reading vocabulary from Excel files using xlsx library
- **pager.go**: Implements pagination logic with 25 words per page
- **ui.go**: Provides interactive command-line interface and navigation
- **main.go**: Command-line application entry point
- **web_main.go**: Web application entry point
- **web_server.go**: Web server and API handlers
- **web/templates/**: HTML templates for web interface
- **web/static/**: CSS, JavaScript, and other static assets

### Data Flow
1. Application starts and reads vocabulary data from `words/IELTS.xlsx`
2. Data is parsed from Excel (column 3: English, column 8: Chinese)
3. Data is stored in WordList structure
4. Pager handles pagination (25 words per page)
5. UI displays words and handles user navigation
6. User can navigate through pages using keyboard commands or web interface

## Development Environment

### Prerequisites
- Go installed on the system
- Excel file: `words/IELTS.xlsx` (3673 IELTS vocabulary words)
- External dependency: `github.com/tealeg/xlsx/v3` for Excel file reading

### Build and Run Commands

#### Command Line Version
- `go build -o word-hero.exe main.go excel_reader.go models.go pager.go ui.go` - Build CLI app
- `go run main.go excel_reader.go models.go pager.go ui.go` - Run CLI directly
- `./word-hero.exe` - Run compiled CLI version

#### Web Version
- `go build -o word-hero-web.exe web_main.go web_server.go excel_reader.go models.go pager.go` - Build web app
- `go run web_main.go web_server.go excel_reader.go models.go pager.go` - Run web directly
- `./word-hero-web.exe` - Run web server on port 8080

### Data Requirements
- Vocabulary data: `words/IELTS.xlsx`
- Worksheet: "雅思真经词汇"
- Data structure: Column 3 contains English words, Column 8 contains Chinese translations
- Header row: Automatically skipped
- Total vocabulary: 3673 words

### Testing
- Command line: `echo -e "n\nq" | ./word-hero.exe` for automated testing
- Web server: Test at `http://localhost:8080` in browser

## User Interface

### Command Line Interface
#### Navigation Commands
- `n` - Next page
- `p` - Previous page
- `f` - First page
- `l` - Last page
- `g` - Go to specific page
- `s` - Show statistics
- `q` - Quit application

#### Display Format
- Shows current page number and total pages
- Displays 25 words per page with English and Chinese translations
- Clean, formatted output with clear navigation instructions

### Web Interface
#### Features
- Modern, responsive design
- Real-time search functionality
- AJAX-based navigation
- Keyboard shortcuts (arrow keys, Ctrl+F, Esc)
- Statistics modal
- Help modal
- Mobile-friendly

#### API Endpoints
- `GET /` - Main interface
- `GET /api/words` - Paginated vocabulary data
- `GET /api/page/{number}` - Specific page
- `GET /api/search?q={query}` - Search functionality
- `GET /api/stats` - Application statistics

## Current Implementation Status

### ✅ Completed Features
- Dual-platform application (CLI + Web)
- Direct Excel file reading (3673 vocabulary words)
- Modular architecture with separate concerns
- Pagination system (25 words per page, 147 total pages)
- Interactive command-line interface
- Modern web interface with search functionality
- Complete API endpoints
- Error handling and user feedback
- Responsive design for mobile devices
- Keyboard shortcuts and accessibility

### 🔄 Future Enhancements
- Learning progress tracking and spaced repetition
- Audio pronunciation support
- Vocabulary testing and quiz modes
- User accounts and data synchronization
- Dark mode support
- Export/import functionality
- Mobile app development

## Code Style and Conventions

### Naming Conventions
- Use CamelCase for exported types and functions
- Use camelCase for private variables and functions
- Clear, descriptive names for all identifiers

### Error Handling
- Return errors with descriptive messages
- Handle errors gracefully in the UI
- Provide user-friendly error messages

### File Organization
- Each major component in its own file
- Clear separation of concerns
- Main application logic in main.go

## Important Notes
- The application reads Excel files directly using `github.com/tealeg/xlsx/v3`
- Data is read from `words/IELTS.xlsx` worksheet "雅思真经词汇"
- Column 3 contains English words, Column 8 contains Chinese translations
- Both CLI and Web versions are fully functional
- Web server runs on port 8080 by default
- Total vocabulary: 3673 words across 147 pages
- Both versions support the same core vocabulary data
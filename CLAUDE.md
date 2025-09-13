# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Word Hero is a web-based Go application for learning IELTS vocabulary, offering a modern web interface. It reads vocabulary data from `words/IELTS.xlsx` (3673 words) and displays it in a paginated format (25 words per page) with interactive navigation and search functionality.

## Project Architecture

### Core Components
- **models.go**: Defines Word, WordList, and Page data structures
- **excel_reader.go**: Handles reading vocabulary from Excel files using xlsx library
- **pager.go**: Implements pagination logic with 25 words per page
- **main.go**: Web application entry point
- **web_server.go**: Web server and API handlers
- **web/templates/**: HTML templates for web interface
- **web/static/**: CSS, JavaScript, and other static assets

### Data Flow
1. Application starts and reads vocabulary data from `words/IELTS.xlsx`
2. Data is parsed from Excel (column 3: English, column 8: Chinese)
3. Data is stored in WordList structure
4. Pager handles pagination (25 words per page)
5. Web server provides REST API endpoints
6. Frontend displays words and handles user navigation

## Development Environment

### Prerequisites
- Go installed on the system
- Excel file: `words/IELTS.xlsx` (3673 IELTS vocabulary words)
- External dependency: `github.com/tealeg/xlsx/v3` for Excel file reading

### Build and Run Commands

#### Web Version
- `go build -o word-hero.exe main.go web_server.go excel_reader.go models.go pager.go` - Build web app
- `go run main.go web_server.go excel_reader.go models.go pager.go` - Run web directly
- `./word-hero.exe` - Run web server on port 8082

### Data Requirements
- Vocabulary data: `words/IELTS.xlsx`
- Worksheet: "雅思真经词汇"
- Data structure: Column 3 contains English words, Column 8 contains Chinese translations
- Header row: Automatically skipped
- Total vocabulary: 3673 words

### Testing
- Web server: Test at `http://localhost:8082` in browser

## User Interface

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
- Web-based application
- Direct Excel file reading (3673 vocabulary words)
- Modular architecture with separate concerns
- Pagination system (25 words per page, 147 total pages)
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

## Development Standards

### Go Code Standards

#### Code Structure
- **File Size Limit**: Individual Go files must not exceed 1000 lines. If a file approaches this limit, consider refactoring into smaller, focused modules.
- **Package Organization**: Each `.go` file should focus on a single responsibility
- **Function Length**: Keep functions concise and focused on a single task
- **Exported Functions**: All exported functions must have Go-style comments explaining purpose, parameters, and return values

#### Code Quality
- **Error Handling**: Always handle errors appropriately; never ignore errors without explicit reason
- **Logging**: Use structured logging with appropriate levels (debug, info, warn, error)
- **Dependencies**: Minimize external dependencies; prefer standard library solutions when possible
- **Memory Management**: Be mindful of memory usage, especially with large datasets

#### Conventions
```go
// Good example of function documentation
func ProcessVocabularyData(data *WordList) (*Page, error) {
    // Function implementation
    return &Page{}, nil
}

// Variable naming
var totalCount int        // Clear, descriptive
var words []string       // Plural for slices
var isActive bool        // Boolean prefix with is/has/can
var maxRetries int       // Descriptive naming
```

#### Performance Considerations
- Avoid unnecessary memory allocations in hot paths
- Use appropriate data structures for the use case
- Consider lazy loading for large datasets
- Implement proper connection pooling for external resources

### HTML/CSS/JavaScript Standards

#### HTML Structure
- **File Size Limit**: Individual HTML files should not exceed 1000 lines
- **Semantic HTML**: Use appropriate HTML5 tags for structure
- **Accessibility**: Include proper ARIA labels and keyboard navigation support
- **Template Organization**: Break large templates into smaller, reusable components

#### CSS Standards
- **File Size Limit**: CSS files should not exceed 1000 lines; organize into multiple files if needed
- **Naming**: Use consistent naming conventions (kebab-case for classes)
- **Responsive Design**: All components must be mobile-friendly
- **Performance**: Minimize CSS bundle size and avoid redundant styles

#### JavaScript Standards
- **File Size Limit**: JavaScript files must not exceed 1000 lines
- **ES6+ Features**: Use modern JavaScript features appropriately
- **Error Handling**: Implement comprehensive error handling with user feedback
- **Performance**: Avoid memory leaks and optimize for performance
- **Security**: Sanitize user input and avoid XSS vulnerabilities

#### Code Quality Examples
```javascript
// Good function documentation
/**
 * Updates the statistics modal with current learning progress
 * @param {Object} statistics - Statistics data object
 * @returns {void}
 */
function updateStatisticsModal(statistics) {
    // Implementation
}

// Error handling pattern
try {
    const data = JSON.parse(response);
    updateUI(data);
} catch (error) {
    console.error('Failed to parse response:', error);
    showError('数据处理失败，请重试');
}
```

### Code Review Guidelines

#### Before Committing
- All code must follow the documented standards
- File size limits must be respected (1000 lines max per file)
- Functions should be focused and single-purpose
- Error handling must be comprehensive
- Code should be self-documenting with clear naming

#### Refactoring Requirements
- If any file exceeds 800 lines, start planning refactoring
- If any file exceeds 1000 lines, refactoring is mandatory before committing
- Break large functions into smaller, focused functions
- Extract common functionality into shared utilities

### Documentation Requirements
- All public APIs must have comprehensive documentation
- Complex business logic should include inline comments
- Configuration options must be documented
- Breaking changes require updated documentation

### Testing Standards
- Write tests for critical functionality
- Test both success and error scenarios
- Maintain test coverage above 80%
- Include integration tests for external dependencies

### Environment Awareness Standards

#### System Environment Tracking
- **Environment File**: Maintain `.tmp.md` file to track current development environment
- **OS Detection**: Always check current operating system before executing system commands
- **Command Adaptation**: Adapt commands based on the target environment (Windows/Linux/macOS)

#### Environment-Specific Command Guidelines
- **Windows Environment**: Use Windows-specific commands (`taskkill`, `netstat`, Windows paths)
- **Linux Environment**: Use Linux-specific commands (`kill`, `pkill`, `netstat`, Unix paths)
- **Cross-Platform**: When possible, use Go's built-in cross-platform capabilities

#### Command Execution Validation
Before executing any system command:
1. Check `.tmp.md` for current environment
2. Verify command compatibility with current OS
3. Use appropriate command syntax for the platform
4. Provide fallback options for cross-platform compatibility

#### Environment File Management
- **Location**: `.tmp.md` in project root (gitignored)
- **Content**: Current OS, architecture, and environment-specific notes
- **Updates**: Modify when switching development environments
- **Format**: Markdown format for easy reading and editing

#### Examples of Environment-Specific Commands
```bash
# Windows process termination
taskkill //PID 12345 //F

# Linux process termination  
kill -9 12345

# Windows network check
netstat -ano | findstr ":8080"

# Linux network check
netstat -tlnp | grep :8080
```

## Important Notes
- The application reads Excel files directly using `github.com/tealeg/xlsx/v3`
- Data is read from `words/IELTS.xlsx` worksheet "雅思真经词汇"
- Column 3 contains English words, Column 8 contains Chinese translations
- Web server runs on port 8082 by default (configurable via configs/config.yaml)
- Total vocabulary: 3673 words across 147 pages

## Development Guidelines
- **Code Standards Compliance**: All development must strictly follow the Development Standards section above
- **File Size Management**: Monitor file sizes and refactor when approaching 800 lines; mandatory refactor at 1000 lines
- **Environment Awareness**: Always check `.tmp.md` for current environment and adapt commands accordingly
- **Cross-Platform Compatibility**: Ensure all system commands work on the current development environment
- **After debugging/testing completion**: Always stop all running background services to free up ports
- **Configuration**: Server settings are loaded from `configs/config.yaml` with environment variable fallbacks
- **Background services**: Use KillBash tool to terminate any lingering background processes after development sessions
- **Code Reviews**: All code changes must be reviewed against the documented standards before committing
- **Documentation Updates**: Update documentation when adding new features or changing existing functionality
- **Environment Tracking**: Keep `.tmp.md` updated when switching development environments
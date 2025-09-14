# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Word Hero is a web-based Go application for learning IELTS vocabulary, offering a modern web interface. It uses the Gin web framework for efficient HTTP routing and middleware, reads vocabulary data from `configs/words/IELTS.xlsx` (3673 words) and displays it in a paginated format (25 words per page) with interactive navigation and search functionality.

## Project Architecture

### Project Structure
```
word-hero/
├── cmd/                    # Application entry points
│   └── main.go            # Main application entry point
├── internal/              # Private application code
│   ├── conf/              # Configuration management
│   │   └── config.go     # Configuration loading and management
│   ├── dao/               # Data Access Objects (data layer)
│   │   ├── models.go      # Data models and structures
│   │   └── excel_reader.go # Excel file operations
│   ├── service/           # Business logic layer
│   │   ├── vocabulary_service.go # Vocabulary operations
│   │   └── pager_service.go      # Pagination logic
│   └── router/            # HTTP routing layer
│       └── web_server.go  # Web server and API handlers
├── configs/               # Configuration files
│   └── config.yaml        # Application configuration
├── web/                   # Web assets
│   ├── templates/         # HTML templates
│   └── static/            # CSS, JavaScript, images
├── configs/words/         # Data files
│   └── IELTS.xlsx        # Vocabulary data source
├── build/                 # Build artifacts
└── go.mod, go.sum        # Go module files
```

### Architecture Layers

#### 1. Presentation Layer (internal/router)
- **web_server.go**: HTTP server using Gin framework, route handling, middleware
- Manages web requests and API endpoints with Gin's routing system
- Handles template rendering and static file serving
- Implements request logging and error handling using Gin middleware
- Uses Gin's context for efficient parameter handling and JSON responses

#### 2. Service Layer (internal/service)
- **vocabulary_service.go**: Business logic for vocabulary operations
- **pager_service.go**: Business logic for pagination and data management
- Orchestrates data flow between presentation and data layers
- Implements search, validation, and business rules

#### 3. Data Access Layer (internal/dao)
- **models.go**: Data structures and domain models
- **excel_reader.go**: Data persistence and retrieval operations
- Handles all database/file operations
- Manages data transformation and validation

#### 4. Configuration Layer (internal/conf)
- **config.go**: Application configuration management
- Handles environment variables and YAML configuration
- Provides configuration validation and defaults

### Data Flow
1. **Application Start**: `cmd/main.go` initializes configuration and services
2. **Configuration Loading**: `internal/conf/config.go` loads settings from YAML and environment
3. **Data Initialization**: `internal/dao/excel_reader.go` reads vocabulary from Excel
4. **Service Setup**: `internal/service/` layer initializes business logic services
5. **Web Server Start**: `internal/router/web_server.go` starts HTTP server with routes
6. **Request Processing**: HTTP requests flow through layers:
   - Router → Service → DAO → Data Source
   - Response flows back: Data Source → DAO → Service → Router → Client

### Design Patterns Used
- **Layered Architecture**: Clear separation of concerns across layers
- **Service Layer Pattern**: Business logic isolated in service layer
- **Data Access Object (DAO)**: Data operations abstracted from business logic
- **Dependency Injection**: Services are injected and dependencies managed explicitly
- **Configuration Management**: Centralized configuration with environment variable support

## Development Environment

### Prerequisites
- Go installed on the system
- Excel file: `configs/words/IELTS.xlsx` (3673 IELTS vocabulary words)
- External dependencies:
  - `github.com/gin-gonic/gin` for web framework
  - `github.com/tealeg/xlsx/v3` for Excel file reading

### Build and Run Commands

#### Manual Go Commands
- `go build -o build/word-hero ./cmd` - Build for current platform
- `go run ./cmd` - Run directly from source
- `./build/word-hero` - Run built application (Windows: `build/word-hero.exe`)

### Data Requirements
- Vocabulary data: `configs/words/IELTS.xlsx`
- Worksheet: "雅思真经词汇"
- Data structure: Column 3 contains English words, Column 8 contains Chinese translations
- Header row: Automatically skipped
- Total vocabulary: 3673 words

### Testing
- Web server: Test at `http://localhost:8080` in browser

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
- `GET /api/page/:pageNumber` - Specific page
- `GET /api/search?q={query}` - Search functionality
- `GET /api/stats` - Application statistics

#### Gin Framework Benefits
- **Efficient Routing**: Parameter-based routing with `:pageNumber` parameters
- **Middleware Support**: Built-in logging, recovery, and custom middleware
- **Context Handling**: Simplified parameter access with `c.Query()` and `c.Param()`
- **JSON Responses**: Automatic JSON serialization with `c.JSON()`
- **HTML Templates**: Built-in template rendering with `c.HTML()`

## Current Implementation Status

### ✅ Completed Features
- Web-based application using Gin framework
- Direct Excel file reading (3673 vocabulary words)
- Modular architecture with separate concerns
- Pagination system (25 words per page, 147 total pages)
- Modern web interface with search functionality
- Complete API endpoints with parameter-based routing
- Error handling and user feedback
- Responsive design for mobile devices
- Keyboard shortcuts and accessibility
- Built-in middleware for logging and recovery
- Efficient HTTP routing and context handling

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

### Database Standards

#### Table Creation Standards
1. **Primary Key Standards**: All primary keys cannot use auto-increment IDs. Use UUID or snowflake algorithm to generate unique IDs:
   ```go
   type User struct {
       ID string `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
   }
   ```

2. **Timestamp Standards**: All time fields should be saved in timestamp format with millisecond precision:
   ```go
   type User struct {
       CreatedAt int64 `gorm:"autoCreateTime:milli" json:"createdAt"`
       UpdatedAt int64 `gorm:"autoUpdateTime:milli" json:"updatedAt"`
   }
   ```

3. **ID Generation**: Use the provided utilities for consistent ID generation:
   - UUID: `utils.GenerateUUID()`
   - Snowflake: `utils.GenerateSnowflakeID()`

4. **Database Migration**: All database migrations must include PostgreSQL UUID extension support:
   ```go
   // Enable UUID extension for PostgreSQL
   if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
       return fmt.Errorf("failed to create UUID extension: %w", err)
   }
   ```

5. **Model Organization**: All data models must be placed in the `internal/models` package with appropriate GORM tags and JSON serialization support.

#### PostgreSQL Configuration
- Default UUID extension: `uuid-ossp`
- UUID generation: `uuid_generate_v4()`
- Automatic timestamp management with GORM hooks
- Proper indexing for UUID primary keys

#### Foreign Key Policy
**STRICT PROHIBITION**: The use of foreign key constraints is strictly prohibited in this project. All database relationships must be managed at the application level through business logic rather than database-level foreign key constraints.

**Reasoning**:
- Foreign keys create coupling between tables that complicates database schema evolution
- They reduce development flexibility and make schema migrations more complex
- Application-level relationship management provides better control over data integrity
- This approach aligns with modern microservices and distributed systems architecture patterns

**Implementation Requirements**:
- All database schemas must avoid `FOREIGN KEY` constraints
- Data relationships must be maintained through application logic in the service layer
- Referential integrity must be enforced programmatically during CRUD operations
- Data consistency checks should be implemented as part of business logic validation

### Go Project Structure Standards

#### Directory Organization
- **cmd/**: Application entry points and main functions
- **internal/**: Private application code (not importable by other projects)
- **internal/conf/**: Configuration management
- **internal/dao/**: Data Access Objects and data models
- **internal/service/**: Business logic and service layer
- **internal/router/**: HTTP routing and web server logic
- **configs/**: Configuration files (YAML, JSON, etc.)
- **web/**: Web assets (templates, static files)
- **build/**: Compiled binaries and build artifacts

#### Go Module Standards
- **Module Naming**: Use repository URL for module name (e.g., `github.com/user/project`)
- **Package Naming**: Use lowercase, concise names that describe package purpose
- **Internal Packages**: Use `internal/` for packages that should not be imported by external projects
- **Cyclic Dependencies**: Avoid cyclic dependencies between packages

#### Code Structure
- **File Size Limit**: Individual Go files must not exceed 1000 lines. If a file approaches this limit, consider refactoring into smaller, focused modules.
- **Package Organization**: Each `.go` file should focus on a single responsibility
- **Function Length**: Keep functions concise and focused on a single task
- **Exported Functions**: All exported functions must have Go-style comments explaining purpose, parameters, and return values

#### Struct Organization Standards
- **Data Access Layer (internal/dao/)**: Contains only data access operations and database interactions
- **Business Logic Layer (internal/models/)**: Contains business logic methods that wrap data structures with domain-specific operations
- **Data Transfer Objects (internal/dto/)**: Contains all API request and response structures, including input validation and output formatting
- **Database Table Structures (internal/table/)**: Contains database entity definitions with GORM tags and persistence-related methods
- **API Interface Standards**: All future API interface input and return parameter struct definitions must be placed in the `internal/dto/` package
- **Separation of Concerns**: Database table structs, business logic, and API communication structures must be kept in separate packages to maintain clear boundaries between layers

#### Architecture Guidelines
- **Layered Architecture**: Follow the established layer pattern (router → service → dao)
- **Dependency Direction**: Dependencies should flow downward (router depends on service, service depends on dao)
- **Interface Segregation**: Define interfaces at the layer boundaries
- **Single Responsibility**: Each layer should have a single, well-defined responsibility

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
- Data is read from `configs/words/IELTS.xlsx` worksheet "雅思真经词汇"
- Column 3 contains English words, Column 8 contains Chinese translations
- Web server runs on port 8080 by default (configurable via configs/config.yaml)
- Total vocabulary: 3673 words across 147 pages
- Web framework: Gin provides efficient HTTP routing, middleware, and context handling

## Development Guidelines

### Architecture Compliance
- **Project Structure**: All new code must follow the established directory structure (cmd/, internal/, configs/, etc.)
- **Layered Architecture**: Maintain clear separation between router, service, and dao layers
- **Dependency Direction**: Ensure dependencies flow downward (router → service → dao)
- **Internal Packages**: Use `internal/` for application-specific code that should not be imported externally

### Code Standards Compliance
- **File Size Management**: Monitor file sizes and refactor when approaching 800 lines; mandatory refactor at 1000 lines
- **Code Quality**: All development must strictly follow the Development Standards section above
- **Testing**: Write comprehensive tests for all service layer functionality
- **Error Handling**: Implement proper error handling at all layers with user-friendly messages

### Development Workflow
- **Build Commands**: Use manual Go commands for building and running
- **Environment Awareness**: Always check `.tmp.md` for current environment and adapt commands accordingly
- **Cross-Platform Compatibility**: Ensure all system commands work on the current development environment
- **After debugging/testing completion**: Always stop all running background services to free up ports

### Configuration Management
- **Configuration**: Server settings are loaded from `configs/config.yaml` with environment variable fallbacks
- **Environment Variables**: Use environment variables for sensitive configuration and deployment-specific settings
- **Configuration Validation**: Implement validation for all configuration values

### Process Management
- **Background services**: Use KillBash tool to terminate any lingering background processes after development sessions
- **Port Management**: Ensure no conflicts with default port (8082) or configured ports
- **Service Cleanup**: Always clean up temporary files and build artifacts

### Code Quality and Reviews
- **Code Reviews**: All code changes must be reviewed against the documented standards before committing
- **Documentation Updates**: Update documentation when adding new features or changing existing functionality
- **Environment Tracking**: Keep `.tmp.md` updated when switching development environments
- **Git Workflow**: Follow established git practices for branching, committing, and merging

### Build and Deployment
- **Build Artifacts**: All build artifacts should be placed in the `build/` directory
- **Manual Building**: Use Go commands directly for building and testing
- **Version Management**: Maintain version information in configuration and build process
- **Dependency Management**: Keep Go modules updated and secure

### API Development Standards
- **API Specifications**: Follow the interface specifications documented in `api/api.md` for all API development
- **Response Format**: Use the standardized `{code, data, msg}` response format for all API endpoints
- **Error Handling**: Implement consistent error codes (0 for success, 150321309 for general errors)
- **Pagination Format**: For paginated endpoints, use `{items: [], total: 1234}` format and let frontend calculate pagination
- **Documentation**: Keep API documentation updated when making changes to endpoints or response formats
- **Testing**: Verify all API endpoints comply with the documented specifications before deployment
package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sanmu2018/word-hero/log"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// AppConfig represents application configuration
type AppConfig struct {
	ExcelFile string `yaml:"excel_file"`
	PageSize  int    `yaml:"page_size"`
	StaticDir string `yaml:"static_dir"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	SSLMode         string `yaml:"ssl_mode"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret    string `yaml:"secret"`
	ExpiresIn string `yaml:"expires_in"`
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	// Default configuration
	config := &Config{
		Server: ServerConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
		App: AppConfig{
			ExcelFile: "configs/words/IELTS.xlsx",
			PageSize:  12,
			StaticDir: "web/static",
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			Password:        "postgres",
			DBName:          "word_hero",
			SSLMode:         "disable",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 3600,
		},
		JWT: JWTConfig{
			Secret:    "your-secret-key-change-in-production",
			ExpiresIn: "24h",
		},
	}

	// Try to load from config file
	configPath := "configs/config.yaml"
	if _, err := os.Stat(configPath); err == nil {
		log.Info().Str("path", configPath).Msg("Loading configuration from file")

		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Error(err).Str("path", configPath).Msg("Failed to read config file")
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			log.Error(err).Str("path", configPath).Msg("Failed to parse config file")
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}

		log.Info().Msg("Configuration loaded successfully")
	} else {
		log.Info().Msg("No config file found, using default configuration")
	}

	// Override with environment variables if provided
	if port := os.Getenv("WORD_HERO_PORT"); port != "" {
		if parsedPort, err := parseEnvInt(port, "port"); err == nil {
			config.Server.Port = parsedPort
			log.Info().Int("port", parsedPort).Msg("Port overridden by environment variable")
		}
	}

	if excelFile := os.Getenv("WORD_HERO_EXCEL_FILE"); excelFile != "" {
		config.App.ExcelFile = excelFile
		log.Info().Str("file", excelFile).Msg("Excel file overridden by environment variable")
	}

	if pageSize := os.Getenv("WORD_HERO_PAGE_SIZE"); pageSize != "" {
		if parsedPageSize, err := parseEnvInt(pageSize, "page size"); err == nil {
			config.App.PageSize = parsedPageSize
			log.Info().Int("pageSize", parsedPageSize).Msg("Page size overridden by environment variable")
		}
	}

	// Ensure Excel file path is absolute
	if !filepath.IsAbs(config.App.ExcelFile) {
		absPath, err := filepath.Abs(config.App.ExcelFile)
		if err != nil {
			log.Error(err).Str("path", config.App.ExcelFile).Msg("Failed to get absolute path for Excel file")
			return nil, fmt.Errorf("failed to get absolute path for Excel file: %w", err)
		}
		config.App.ExcelFile = absPath
	}

	return config, nil
}

// parseEnvInt parses an environment variable as integer
func parseEnvInt(value string, fieldName string) (int, error) {
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		log.Error(err).Str("value", value).Str("field", fieldName).Msg("Failed to parse environment variable as integer")
		return 0, fmt.Errorf("invalid %s value: %s", fieldName, value)
	}
	return result, nil
}
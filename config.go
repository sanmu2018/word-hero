package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	App struct {
		Name      string `yaml:"name"`
		PageSize  int    `yaml:"page_size"`
		ExcelFile string `yaml:"excel_file"`
	} `yaml:"app"`
	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`
}

// LoadConfig loads configuration from environment variables with fallback to config file
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: struct {
			Port int    `yaml:"port"`
			Host string `yaml:"host"`
		}{
			Port: 8082,
			Host: "localhost",
		},
		App: struct {
			Name      string `yaml:"name"`
			PageSize  int    `yaml:"page_size"`
			ExcelFile string `yaml:"excel_file"`
		}{
			Name:      "Word Hero",
			PageSize:  25,
			ExcelFile: "words/IELTS.xlsx",
		},
		Logging: struct {
			Level string `yaml:"level"`
			File  string `yaml:"file"`
		}{
			Level: "info",
			File:  "",
		},
	}

	// Load from environment variables if they exist
	if port := os.Getenv("WORD_HERO_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}

	if host := os.Getenv("WORD_HERO_HOST"); host != "" {
		config.Server.Host = host
	}

	if pageSize := os.Getenv("WORD_HERO_PAGE_SIZE"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			config.App.PageSize = ps
		}
	}

	if excelFile := os.Getenv("WORD_HERO_EXCEL_FILE"); excelFile != "" {
		config.App.ExcelFile = excelFile
	}

	return config, nil
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetURL returns the full URL for the server
func (c *Config) GetURL() string {
	return fmt.Sprintf("http://%s:%d", c.Server.Host, c.Server.Port)
}
package dao

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sanmu2018/word-hero/internal/conf"
	"github.com/sanmu2018/word-hero/log"
)

// DB holds the database connection
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase(config *conf.DatabaseConfig) error {
	// First connect to default 'postgres' database to create our database if it doesn't exist
	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.Port,
		config.SSLMode,
	)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to postgres database first
	defaultDB, err := gorm.Open(postgres.Open(defaultDSN), gormConfig)
	if err != nil {
		log.Error(err).Msg("Failed to connect to default postgres database")
		return fmt.Errorf("failed to connect to default postgres database: %w", err)
	}

	// Check if our database exists
	var exists bool
	err = defaultDB.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", config.DBName).Scan(&exists).Error
	if err != nil {
		log.Error(err).Msg("Failed to check if database exists")
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Create database if it doesn't exist
	if !exists {
		log.Info().Str("database", config.DBName).Msg("Database does not exist, creating it...")
		err = defaultDB.Exec("CREATE DATABASE " + config.DBName).Error
		if err != nil {
			log.Error(err).Str("database", config.DBName).Msg("Failed to create database")
			return fmt.Errorf("failed to create database %s: %w", config.DBName, err)
		}
		log.Info().Str("database", config.DBName).Msg("Database created successfully")
	}

	// Close the default connection
	sqlDB, err := defaultDB.DB()
	if err != nil {
		log.Error(err).Msg("Failed to get underlying sql.DB from default connection")
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.Close()

	// Now connect to our specific database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	// Connect to the target database
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Error(err).Msg("Failed to connect to target database")
		return fmt.Errorf("failed to connect to target database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err = db.DB()
	if err != nil {
		log.Error(err).Msg("Failed to get underlying sql.DB")
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		log.Error(err).Msg("Failed to ping database")
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Info().Str("host", config.Host).Int("port", config.Port).Str("dbname", config.DBName).Msg("Database connection established successfully")

	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDatabase returns the database connection
func GetDatabase() *gorm.DB {
	return DB
}
func GetDb(ctx context.Context) *gorm.DB {
	return DB.WithContext(ctx)
}

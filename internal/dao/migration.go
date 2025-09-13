package dao

import (
	"fmt"

	"github.com/sanmu2018/word-hero/internal/models"
	"github.com/sanmu2018/word-hero/internal/table"
	"github.com/sanmu2018/word-hero/log"
)

// AutoMigrate performs automatic database migration
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Info().Msg("Starting database migration...")

	// Enable UUID extension for PostgreSQL
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Error(err).Msg("Failed to create UUID extension")
		return fmt.Errorf("failed to create UUID extension: %w", err)
	}

	// Migrate all models
	err := DB.AutoMigrate(
		&table.User{},
		&table.Word{},
	)
	if err != nil {
		log.Error(err).Msg("Database migration failed")
		return fmt.Errorf("database migration failed: %w", err)
	}

	log.Info().Msg("Database migration completed successfully")
	return nil
}

// RunMigrations runs all database migrations and setup
func RunMigrations() error {
	if err := AutoMigrate(); err != nil {
		return err
	}
	if err := CreateDefaultUser(); err != nil {
		return err
	}
	return nil
}

// CreateDefaultUser creates a default admin user if no users exist
func CreateDefaultUser() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	var userCount int64
	if err := DB.Model(&table.User{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if userCount == 0 {
		log.Info().Msg("No users found, creating default admin user...")

	// In a real application, you would want to generate a secure random password
	// and display it to the user or send it via email
	defaultUser := &table.User{
			Username: "admin",
			Email:    "admin@wordhero.com",
			FullName: "System Administrator",
			Role:     "admin",
		}

	// Hash the password (you should use a secure password in production)
	// For now, we'll use a placeholder - this should be changed
	businessUser := models.NewUserBusiness(defaultUser)
	if err := businessUser.SetPassword("admin123"); err != nil {
			return fmt.Errorf("failed to set default password: %w", err)
		}

		if err := DB.Create(defaultUser).Error; err != nil {
			log.Error(err).Msg("Failed to create default user")
			return fmt.Errorf("failed to create default user: %w", err)
		}

		log.Info().
			Str("user_id", defaultUser.ID).
			Str("username", defaultUser.Username).
			Msg("Default admin user created successfully")
		log.Warn().
			Msg("IMPORTANT: Please change the default admin password immediately!")
	}

	return nil
}
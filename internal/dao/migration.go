package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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
		&table.WordTag{},
	)
	if err != nil {
		log.Error(err).Msg("Database migration failed")
		return fmt.Errorf("database migration failed: %w", err)
	}

	log.Info().Msg("Database migration completed successfully")
	return nil
}

// MigrateWordTagsTable migrates the word_tags table from the old JSONB structure to the new simple timestamp structure
func MigrateWordTagsTable() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	log.Info().Msg("Starting word_tags table migration...")

	// Check if the old UserTags column exists
	var columnExists bool
	err := DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'word_tags' AND column_name = 'user_tags')").Scan(&columnExists).Error
	if err != nil {
		return fmt.Errorf("failed to check for UserTags column: %w", err)
	}

	// If UserTags column doesn't exist, the migration has already been done
	if !columnExists {
		log.Info().Msg("UserTags column not found, migration already completed")
		return nil
	}

	// Check if the new Known column already exists
	var knownColumnExists bool
	err = DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'word_tags' AND column_name = 'known')").Scan(&knownColumnExists).Error
	if err != nil {
		return fmt.Errorf("failed to check for Known column: %w", err)
	}

	// If Known column already exists, drop it and recreate
	if knownColumnExists {
		log.Info().Msg("Known column already exists, dropping it...")
		if err := DB.Exec("ALTER TABLE word_tags DROP COLUMN known").Error; err != nil {
			return fmt.Errorf("failed to drop existing Known column: %w", err)
		}
	}

	// Check if the new UserID column already exists
	var userIDColumnExists bool
	err = DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'word_tags' AND column_name = 'user_id')").Scan(&userIDColumnExists).Error
	if err != nil {
		return fmt.Errorf("failed to check for UserID column: %w", err)
	}

	// If UserID column already exists, drop it and recreate
	if userIDColumnExists {
		log.Info().Msg("UserID column already exists, dropping it...")
		if err := DB.Exec("ALTER TABLE word_tags DROP COLUMN user_id").Error; err != nil {
			return fmt.Errorf("failed to drop existing UserID column: %w", err)
		}
	}

	// Add the new Known column
	log.Info().Msg("Adding Known column...")
	if err := DB.Exec("ALTER TABLE word_tags ADD COLUMN known BIGINT").Error; err != nil {
		return fmt.Errorf("failed to add Known column: %w", err)
	}

	// Add the new UserID column
	log.Info().Msg("Adding UserID column...")
	if err := DB.Exec("ALTER TABLE word_tags ADD COLUMN user_id UUID NOT NULL DEFAULT uuid_generate_v4()").Error; err != nil {
		return fmt.Errorf("failed to add UserID column: %w", err)
	}

	// Create index for UserID
	log.Info().Msg("Creating index for UserID...")
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_word_tags_user_id ON word_tags (user_id)").Error; err != nil {
		return fmt.Errorf("failed to create UserID index: %w", err)
	}

	// Migrate data from UserTags JSONB to Known timestamp
	log.Info().Msg("Migrating data from UserTags to Known...")

	type OldWordTag struct {
		ID       string
		WordID   string
		UserTags []byte `gorm:"type:jsonb"`
	}

	var oldTags []OldWordTag
	if err := DB.Table("word_tags").Select("id, word_id, user_tags").Find(&oldTags).Error; err != nil {
		return fmt.Errorf("failed to fetch old word tags: %w", err)
	}

	for _, oldTag := range oldTags {
		var userTags map[string]interface{}
		if len(oldTag.UserTags) > 0 {
			if err := json.Unmarshal(oldTag.UserTags, &userTags); err != nil {
				log.Error(err).Str("word_tag_id", oldTag.ID).Msg("Failed to unmarshal UserTags")
				continue
			}

			// Extract user data and create new records
			for userID, tagData := range userTags {
				if tagMap, ok := tagData.(map[string]interface{}); ok {
					var knownTime sql.NullInt64

					// Check if the user has marked this word as known
					if isKnown, exists := tagMap["is_known"]; exists {
						if known, ok := isKnown.(bool); ok && known {
							// Use current time as the known timestamp
							knownTime = sql.NullInt64{
								Int64: time.Now().UnixMilli(),
								Valid: true,
							}
						}
					} else if markedAt, exists := tagMap["marked_at"]; exists {
						// If marked_at timestamp exists, use it
						if timestamp, ok := markedAt.(float64); ok {
							knownTime = sql.NullInt64{
								Int64: int64(timestamp),
								Valid: true,
							}
						}
					}

					// Create new record for this user-word pair
					newTag := table.WordTag{
						WordID: oldTag.WordID,
						UserID: userID,
						Known:  func() *int64 {
							if knownTime.Valid {
								return &knownTime.Int64
							}
							return nil
						}(),
					}

					if err := DB.Create(&newTag).Error; err != nil {
						log.Error(err).Str("word_tag_id", oldTag.ID).Str("user_id", userID).Msg("Failed to create new word tag")
						continue
					}

					log.Info().
						Str("old_id", oldTag.ID).
						Str("new_id", newTag.ID).
						Str("user_id", userID).
						Str("word_id", oldTag.WordID).
						Msg("Migrated word tag")
				}
			}
		}
	}

	// Drop the old UserTags column
	log.Info().Msg("Dropping old UserTags column...")
	if err := DB.Exec("ALTER TABLE word_tags DROP COLUMN user_tags").Error; err != nil {
		return fmt.Errorf("failed to drop UserTags column: %w", err)
	}

	// The old records will be cleaned up by AutoMigrate since they don't match the new schema
	log.Info().Msg("word_tags table migration completed successfully")
	return nil
}

// RunMigrations runs all database migrations and setup
func RunMigrations() error {
	if err := AutoMigrate(); err != nil {
		return err
	}
	if err := MigrateWordTagsTable(); err != nil {
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
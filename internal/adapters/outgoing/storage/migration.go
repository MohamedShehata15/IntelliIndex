package storage

import (
	"fmt"
	"log"

	"gorm.io/gorm"

	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/storage/models"
)

// MigrationHandler handles database migrations
type MigrationHandler struct {
	db *gorm.DB
}

// NewMigrationHandler creates a new migration handler
func NewMigrationHandler(client *Client) *MigrationHandler {
	return &MigrationHandler{
		db: client.DB,
	}
}

// RunMigrations runs all migrations for the database
func (m *MigrationHandler) RunMigrations() error {
	log.Println("Running database migrations...")
	for _, dbModel := range m.models() {
		if err := m.db.AutoMigrate(dbModel); err != nil {
			if err := m.db.AutoMigrate(dbModel); err != nil {
				return fmt.Errorf("failed to migrate %T: %w", dbModel, err)
			}
		}
	}
	log.Println("Database migrations completed successfully")
	return nil
}

// ResetDatabase drops all tables and reruns migrations
func (m *MigrationHandler) ResetDatabase() error {
	log.Println("WARNING: Resetting database - all data will be lost")
	for _, dbModel := range m.models() {
		if err := m.db.Migrator().DropTable(dbModel); err != nil {
			return fmt.Errorf("failed to drop table for %T: %w", dbModel, err)
		}
	}
	return m.RunMigrations()
}

func (m *MigrationHandler) models() []interface{} {
	return []interface{}{
		&models.Document{},
		&models.DocumentMetadata{},
		&models.DocumentKeyword{},
		&models.DocumentLink{},
		&models.DocumentTag{},
		&models.Index{},
	}
}

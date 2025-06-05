package storage

import "gorm.io/gorm"

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

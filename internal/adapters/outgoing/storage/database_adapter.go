package storage

import (
	"fmt"
	"github.com/mohamedshehata15/intelli-index/pkg/config"
)

// SQLAdapter is a unified adapter for SQL database operations
type SQLAdapter struct {
	client           *Client
	documentRepo     *DocumentRepository
	indexRepo        *IndexRepository
	migrationHandler *MigrationHandler
}

// NewSQLAdapter creates a new SQL Adapter
func NewSQLAdapter(cfg *config.DBConfig) (*SQLAdapter, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}

	documentRepo := NewDocumentRepository(client)
	indexRepo := NewIndexRepository(client)
	migrationHandler := NewMigrationHandler(client)

	adapter := &SQLAdapter{
		client,
		documentRepo,
		indexRepo,
		migrationHandler,
	}
	return adapter, nil
}

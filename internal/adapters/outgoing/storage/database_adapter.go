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

// Initialize initializes the database, running migrations if configured
func (s *SQLAdapter) Initialize() error {
	if err := s.client.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	if s.client.Config.AutoMigrate {
		if err := s.migrationHandler.RunMigrations(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}
	return nil
}

package storage

import (
	"fmt"
	"github.com/mohamedshehata15/intelli-index/internal/pkg/di"
	"github.com/mohamedshehata15/intelli-index/pkg/config"
)

// StorageAdapterFactory encapsulates the configuration and registration logic
type StorageAdapterFactory struct {
	config *config.DBConfig
}

var _ di.AdapterRegistrar = (*StorageAdapterFactory)(nil)

func NewStorageAdapterFactory(cfg *config.DBConfig) *StorageAdapterFactory {
	return &StorageAdapterFactory{
		config: cfg,
	}
}

func (s *StorageAdapterFactory) Register(container *di.Container) error {
	return RegisterStorageAdapters(container, s.config)
}

// RegisterStorageAdapters registers all database-related adapters with the DI container
func RegisterStorageAdapters(container *di.Container, cfg *config.DBConfig) error {

	// Register SQLAdapter
	container.Register("sqlAdapter", func() (interface{}, error) {
		adapter, err := NewSQLAdapter(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create SQL adapter: %w", err)
		}

		if err := adapter.Initialize(); err != nil {
			return nil, fmt.Errorf("failed to initialize SQL adapter: %w", err)
		}

		return adapter, nil
	})

	// Register the database client
	container.Register("dbClient", func() (interface{}, error) {
		adapter := GetSQLAdapter(container)
		return adapter.client, nil
	})

	// Register document repository implementation
	container.Register("documentRepositoryDB", func() (interface{}, error) {
		adapter := GetSQLAdapter(container)
		return adapter.DocumentRepository(), nil
	})

	// Register index repository implementation
	container.Register("indexRepositoryDB", func() (interface{}, error) {
		adapter := GetSQLAdapter(container)
		return adapter.IndexRepository(), nil
	})

	// Register migration handler
	container.Register("migrationHandler", func() (interface{}, error) {
		adapter := GetSQLAdapter(container)
		return adapter.MigrationHandler(), nil
	})

	return nil
}

// GetSQLAdapter retrieves the SQL adapter from the container
func GetSQLAdapter(container *di.Container) *SQLAdapter {
	return container.MustResolve("sqlAdapter").(*SQLAdapter)
}

// GetDocumentRepository retrieves the document repository from the container
func GetDocumentRepository(container *di.Container) *DocumentRepository {
	return container.MustResolve("documentRepositoryDB").(*DocumentRepository)
}

// GetIndexRepository retrieves the index repository from the container
func GetIndexRepository(container *di.Container) *IndexRepository {
	return container.MustResolve("indexRepositoryDB").(*IndexRepository)
}

// GetMigrationHandler retrieves the migration handler from the container
func GetMigrationHandler(container *di.Container) *MigrationHandler {
	return container.MustResolve("migrationHandler").(*MigrationHandler)
}

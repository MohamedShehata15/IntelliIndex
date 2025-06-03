package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/storage/models"
	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// IndexRepository implements the outgoing.IndexRepository interface using GORM
type IndexRepository struct {
	db *gorm.DB
}

// NewIndexRepository creates a new gorm index repository
func NewIndexRepository(db *gorm.DB) *IndexRepository {
	return &IndexRepository{
		db: db,
	}
}

// Ensure IndexRepository implements the outgoing.IndexRepository interface
var _ outgoing.IndexRepository = (*IndexRepository)(nil)

func (i IndexRepository) Create(ctx context.Context, index *domain.Index) error {
	if index == nil {
		return errors.New("index cannot bi nil")
	}
	dbIndex := &models.Index{}
	if err := dbIndex.FromDomain(index); err != nil {
		return fmt.Errorf("failed to convert domain model to database model: %w", err)
	}
	if err := i.db.WithContext(ctx).Create(dbIndex).Error; err != nil {
		if i.isUniqueConstraintViolation(err) {
			return fmt.Errorf("index with name %s already exists", index.Name)
		}
		return fmt.Errorf("failed to create index: %w", err)
	}
	index.ID = dbIndex.ID

	return nil
}

func (i IndexRepository) GetByID(ctx context.Context, id string) (*domain.Index, error) {
	if id == "" {
		return nil, errors.New("index ID cannot be empty")
	}

	var dbIndex models.Index
	if err := i.db.WithContext(ctx).First(&dbIndex, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get index by ID: %w", err)
	}

	index, err := dbIndex.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert database model to domain model: %w", err)
	}

	return index, nil
}

func (i IndexRepository) GetByName(ctx context.Context, name string) (*domain.Index, error) {
	if name == "" {
		return nil, errors.New("index name cannot be empty")
	}

	var dbIndex models.Index
	if err := i.db.WithContext(ctx).First(&dbIndex, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get index by name: %w", err)
	}

	index, err := dbIndex.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert database model to domain model: %w", err)
	}

	return index, nil
}

func (i IndexRepository) List(ctx context.Context) ([]*domain.Index, error) {
	var dbIndices []models.Index
	if err := i.db.WithContext(ctx).Find(&dbIndices).Error; err != nil {
		return nil, fmt.Errorf("failed to list indices: %w", err)
	}

	indices := make([]*domain.Index, 0, len(dbIndices))
	for _, dbIndex := range dbIndices {
		index, err := dbIndex.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert database model to domain model: %w", err)
		}
		indices = append(indices, index)
	}

	return indices, nil
}

func (i IndexRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("index ID cannot be empty")
	}

	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		documentsSubQuery := tx.Model(&models.Document{}).Select("id").Where("index_id = ?", id)

		// Delete related document data
		if err := deleteDocumentRelations(tx, documentsSubQuery); err != nil {
			return err
		}

		// Delete all documents with index_id
		if err := tx.Where("index_id = ?", id).Delete(&models.Document{}).Error; err != nil {
			return fmt.Errorf("failed to delete documents for index: %w", err)
		}

		// Delete index
		result := tx.Delete(&models.Index{}, "id = ?", id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete index: %w", result.Error)
		}

		return nil
	})
}

func deleteDocumentRelations(tx *gorm.DB, documentsSubQuery *gorm.DB) error {
	// Delete document keywords
	if err := tx.Where("document_id IN (?)", documentsSubQuery).
		Delete(&models.DocumentKeyword{}).Error; err != nil {
		return fmt.Errorf("failed to delete document keywords: %w", err)
	}

	// Delete document links
	if err := tx.Where("source_id IN (?)", documentsSubQuery).
		Delete(&models.DocumentLink{}).Error; err != nil {
		return fmt.Errorf("failed to delete document links: %w", err)
	}

	// Delete document metadata
	if err := tx.Where("document_id IN (?)", documentsSubQuery).
		Delete(&models.DocumentMetadata{}).Error; err != nil {
		return fmt.Errorf("failed to delete document metadata: %w", err)
	}

	return nil
}

func (i IndexRepository) Update(ctx context.Context, index *domain.Index) error {
	if index == nil {
		return errors.New("index cannot be nil")
	}

	if index.ID == "" {
		return errors.New("index ID cannot be empty")
	}

	var existingIndex models.Index
	if err := i.db.WithContext(ctx).First(&existingIndex, "id = ?", index.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("index with ID %s does not exist", index.ID)
		}
		return fmt.Errorf("failed to check if index exists: %w", err)
	}

	if err := existingIndex.FromDomain(index); err != nil {
		return fmt.Errorf("failed to convert domain model to database model: %w", err)
	}

	if err := i.db.WithContext(ctx).Model(&existingIndex).Updates(models.Index{
		Name:         existingIndex.Name,
		Description:  existingIndex.Description,
		SettingsJSON: existingIndex.SettingsJSON,
		MappingsJSON: existingIndex.MappingsJSON,
	}).Error; err != nil {
		if i.isUniqueConstraintViolation(err) {
			return fmt.Errorf("index with name %s already exists", index.Name)
		}
		return fmt.Errorf("failed to update index: %w", err)
	}

	return nil
}

func (i IndexRepository) UpdateSettings(ctx context.Context, id string, settings domain.IndexSettings) error {
	if id == "" {
		return errors.New("index ID cannot be empty")
	}

	var existingIndex models.Index
	if err := i.db.WithContext(ctx).First(&existingIndex, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("index with ID %s does not exist", id)
		}
		return fmt.Errorf("failed to get index: %w", err)
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings to JSON: %w", err)
	}

	if err := i.db.WithContext(ctx).Model(&existingIndex).Update("settings", string(settingsJSON)).Error; err != nil {
		return fmt.Errorf("failed to update index settings: %w", err)
	}

	return nil
}

func (i IndexRepository) GetStats(ctx context.Context, id string) (map[string]interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexRepository) RefreshIndex(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (r *IndexRepository) isUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "UNIQUE constraint failed")
}

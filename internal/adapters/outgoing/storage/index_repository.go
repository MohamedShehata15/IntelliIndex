package storage

import (
	"context"
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
	//TODO implement me
	panic("implement me")
}

func (i IndexRepository) List(ctx context.Context) ([]*domain.Index, error) {
	//TODO implement me
	panic("implement me")
}

func (i IndexRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexRepository) Update(ctx context.Context, index *domain.Index) error {
	//TODO implement me
	panic("implement me")
}

func (i IndexRepository) UpdateSettings(ctx context.Context, id string, settings domain.IndexSettings) error {
	//TODO implement me
	panic("implement me")
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

package storage

import (
	"context"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"

	"gorm.io/gorm"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
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
	//TODO implement me
	panic("implement me")
}

func (i IndexRepository) GetByID(ctx context.Context, id string) (*domain.Index, error) {
	//TODO implement me
	panic("implement me")
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

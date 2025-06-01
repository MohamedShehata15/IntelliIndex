package outgoing

import (
	"context"
	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
)

// IndexRepository defines the interface for index storage operations
type IndexRepository interface {
	Create(ctx context.Context, index *domain.Index) error
	GetByID(ctx context.Context, id string) (*domain.Index, error)
	GetByName(ctx context.Context, name string) (*domain.Index, error)
	List(ctx context.Context) ([]*domain.Index, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, index *domain.Index) error
	UpdateSettings(ctx context.Context, id string, settings domain.IndexSettings) error
	GetStats(ctx context.Context, id string) (map[string]interface{}, error)
	RefreshIndex(ctx context.Context, id string) error
}

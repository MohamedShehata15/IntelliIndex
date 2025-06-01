package outgoing

import (
	"context"
	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
)

// DocumentRepository defines the interface for document storage operations
type DocumentRepository interface {
	Save(ctx context.Context, document *domain.Document) error
	GetByID(ctx context.Context, id string) (*domain.Document, error)
	GetByURL(ctx context.Context, url string) (*domain.Document, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, document *domain.Document) error
	List(ctx context.Context, page, pageSize int) ([]*domain.Document, int, error)
	Search(ctx context.Context, query *domain.SearchQuery) ([]*domain.Document, int, error)
	CountByIndexID(ctx context.Context, indexID string) (int, error)
}

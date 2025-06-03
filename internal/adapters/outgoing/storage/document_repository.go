package storage

import (
	"context"

	"gorm.io/gorm"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// DocumentRepository implements the outgoing.DocumentRepository interface using GORM
type DocumentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository creates a new document repository
type NewDocumentRepository struct {
	db *gorm.DB
}

// Ensure DocumentRepository implements the outgoing.DocumentRepository interface
var _ outgoing.DocumentRepository = (*DocumentRepository)(nil)

func (d DocumentRepository) Save(ctx context.Context, document *domain.Document) error {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) GetByID(ctx context.Context, id string) (*domain.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) GetByURL(ctx context.Context, url string) (*domain.Document, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) Update(ctx context.Context, document *domain.Document) error {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Document, int, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) Search(ctx context.Context, query *domain.SearchQuery) ([]*domain.Document, int, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) CountByIndexID(ctx context.Context, indexID string) (int, error) {
	//TODO implement me
	panic("implement me")
}

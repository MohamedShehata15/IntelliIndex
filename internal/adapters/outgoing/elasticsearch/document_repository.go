package elasticsearch

import (
	"context"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

const (
	DocumentIndex = "documents"
)

// DocumentRepository implements the outgoing.DocumentRepository interface using Elasticsearch
type DocumentRepository struct {
	client *Client
}

var _ outgoing.DocumentRepository = (*DocumentRepository)(nil)

func NewDocumentRepository(client *Client) *DocumentRepository {
	return &DocumentRepository{
		client,
	}
}

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

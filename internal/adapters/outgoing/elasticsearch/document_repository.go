package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/elasticsearch/models"

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
	if document == nil {
		return errors.New("document cannot be nil")
	}
	if err := document.Validate(); err != nil {
		return fmt.Errorf("invalid document: %w", err)
	}
	if document.ID == "" {
		document.ID = uuid.NewString()
	}
	modelDocument := models.FromDomain(document)
	if err := d.client.IndexDocument(ctx, DocumentIndex, document.ID, modelDocument); err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	return nil
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

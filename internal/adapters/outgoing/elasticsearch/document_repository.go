package elasticsearch

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
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
	if id == "" {
		return nil, errors.New("document ID cannot be empty")
	}
	modelDocument, err := d.client.GetDocument(ctx, DocumentIndex, id)
	if err != nil {
		if strings.Contains(err.Error(), "document not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting document by ID: %w", err)
	}
	return modelDocument.ToDomain(), nil
}

func (d DocumentRepository) GetByURL(ctx context.Context, url string) (*domain.Document, error) {
	if url == "" {
		return nil, errors.New("document URL cannot be empty")
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"url.keyword": url,
			},
		},
	}
	res, err := d.client.PerformRequest(ctx, &esapi.SearchRequest{
		Index: []string{d.client.IndexNameWithPrefix(DocumentIndex)},
		Body:  bytes.NewReader(mustMarshalJSON(query)),
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing for document by URL: %w", err)
	}
	var searchResult map[string]interface{}
	if err := parseResponse(res.Body, &searchResult); err != nil {
		return nil, fmt.Errorf("error parsing search response: %w", err)
	}
	hitsObj := searchResult["hits"].(map[string]interface{})
	hitsTotal := int(hitsObj["total"].(map[string]interface{})["value"].(float64))
	if hitsTotal == 0 {
		return nil, nil
	}
	hits := hitsObj["hits"].([]interface{})
	hit := hits[0].(map[string]interface{})
	doc := hit["_source"].(*models.Document)
	return doc.ToDomain(), nil
}

func (d DocumentRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("document ID cannot be empty")
	}
	if err := d.client.DeleteDocument(ctx, DocumentIndex, id); err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	return nil
}

func (d DocumentRepository) Update(ctx context.Context, document *domain.Document) error {
	if document == nil {
		return errors.New("document cannot be nil")
	}
	if document.ID == "" {
		return errors.New("document ID cannot be empty")
	}
	if err := document.Validate(); err != nil {
		return fmt.Errorf("invalid document: %w", err)
	}
	exists, err := d.client.DocumentExists(ctx, DocumentIndex, document.ID)
	if err != nil {
		return fmt.Errorf("error checking if document exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("document with ID %s does not exist", document.ID)
	}
	modelDocument := models.FromDomain(document)
	if err := d.client.UpdateDocument(ctx, DocumentIndex, document.ID, modelDocument); err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	return nil
}

func (d DocumentRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Document, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	from := (page - 1) * pageSize
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"from": from,
		"size": pageSize,
		"sort": []map[string]interface{}{
			{"last_crawled": map[string]interface{}{"order": "desc"}},
		},
	}
	res, err := d.client.PerformRequest(ctx, &esapi.SearchRequest{
		Index: []string{d.client.IndexNameWithPrefix(DocumentIndex)},
		Body:  bytes.NewReader(mustMarshalJSON(query)),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing for documents: %w", err)
	}
	var searchResult map[string]interface{}
	if err := parseResponse(res.Body, &searchResult); err != nil {
		return nil, 0, fmt.Errorf("error parsing search response: %w", err)
	}
	hitsObj := searchResult["hits"].(map[string]interface{})
	total := int(hitsObj["total"].(map[string]interface{})["value"].(float64))
	hits := hitsObj["hits"].([]interface{})

	documents := make([]*domain.Document, 0, len(hits))
	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		doc := hitMap["_source"].(*models.Document)
		doc.ID = hitMap["_id"].(string)
		if score, ok := hitMap["_score"].(float64); ok && score > 0 {
			doc.Score = score
		}
		documents = append(documents, doc.ToDomain())
	}
	return documents, total, nil
}

func (d DocumentRepository) Search(ctx context.Context, query *domain.SearchQuery) ([]*domain.Document, int, error) {
	//TODO implement me
	panic("implement me")
}

func (d DocumentRepository) CountByIndexID(ctx context.Context, indexID string) (int, error) {
	if indexID == "" {
		return 0, errors.New("index ID cannot be empty")
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"index_id.keyword": indexID,
			},
		},
	}
	res, err := d.client.PerformRequest(ctx, &esapi.SearchRequest{
		Index: []string{d.client.IndexNameWithPrefix(DocumentIndex)},
		Body:  bytes.NewReader(mustMarshalJSON(query)),
	})
	if err != nil {
		return 0, fmt.Errorf("error parsing for documents: %w", err)
	}
	var countResult map[string]interface{}
	if err := parseResponse(res.Body, &countResult); err != nil {
		return 0, fmt.Errorf("error parsing count response: %w", err)
	}
	count := int(countResult["count"].(float64))
	return count, nil
}

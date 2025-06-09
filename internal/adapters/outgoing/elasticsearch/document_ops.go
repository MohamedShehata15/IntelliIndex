package elasticsearch

import (
	"bytes"
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/elasticsearch/models"
)

func (c *Client) IndexDocument(ctx context.Context, indexName, docID string, doc interface{}) error {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	docJSON := mustMarshalJSON(doc)
	res, err := c.PerformRequest(ctx, &esapi.IndexRequest{
		Index:      fullIndexName,
		DocumentID: docID,
		Body:       bytes.NewReader(docJSON),
	})
	if err != nil {
		return err
	}
	defer closeBody(res.Body)
	return nil
}

func (c *Client) GetDocument(ctx context.Context, indexName, docID string) (models.Document, error) {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	res, err := c.PerformRequest(ctx, &esapi.GetRequest{
		Index:      fullIndexName,
		DocumentID: docID,
	})
	if err != nil {
		return models.Document{}, err
	}
	var response map[string]interface{}
	if err := parseResponse(res.Body, &response); err != nil {
		return models.Document{}, fmt.Errorf("error parsing get response: %w", err)
	}
	if found, ok := response["found"].(bool); !ok || !found {
		return models.Document{}, fmt.Errorf("document not found")
	}
	return response["_source"].(models.Document), err
}

func (c *Client) DeleteDocument(ctx context.Context, indexName, docID string) error {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	res, err := c.PerformRequest(ctx, &esapi.DeleteRequest{
		Index:      fullIndexName,
		DocumentID: docID,
	})
	if err != nil {
		return err
	}
	defer closeBody(res.Body)
	return nil
}

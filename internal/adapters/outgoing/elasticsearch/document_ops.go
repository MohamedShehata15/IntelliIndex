package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (c *Client) UpdateDocument(ctx context.Context, indexName, docID string, doc interface{}) error {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	updateBody := map[string]interface{}{
		"doc": doc,
	}
	updateJSON, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("error marshaling update body: %w", err)
	}
	res, err := c.PerformRequest(ctx, &esapi.UpdateRequest{
		Index:      fullIndexName,
		DocumentID: docID,
		Body:       bytes.NewReader(updateJSON),
	})
	if err != nil {
		return err
	}
	defer closeBody(res.Body)
	return nil
}

func (c *Client) DocumentExists(ctx context.Context, indexName, docID string) (bool, error) {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	res, err := c.PerformRequest(ctx, &esapi.ExistsRequest{
		Index:      fullIndexName,
		DocumentID: docID,
	})
	if err != nil {
		return false, err
	}
	defer closeBody(res.Body)
	return res.StatusCode == http.StatusOK, nil
}

func (c *Client) CountDocument(ctx context.Context, indexName string, query map[string]interface{}) (int64, error) {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	var body io.Reader
	if query != nil {
		bodyJSON, err := json.Marshal(query)
		if err != nil {
			return 0, fmt.Errorf("error marshaling query: %w", err)
		}
		body = bytes.NewReader(bodyJSON)
	}
	res, err := c.PerformRequest(ctx, &esapi.CountRequest{
		Index: []string{fullIndexName},
		Body:  body,
	})
	if err != nil {
		return 0, err
	}
	var countResponse map[string]interface{}
	if err := parseResponse(res.Body, &countResponse); err != nil {
		return 0, fmt.Errorf("error parsing count response: %w", err)
	}
	count, ok := countResponse["count"].(float64)
	if !ok {
		return 0, fmt.Errorf("error parsing count response format")
	}
	return int64(count), nil
}

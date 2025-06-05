package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// IndexMapping represents the structure of an Elasticsearch index mapping
type IndexMapping struct {
	Settings IndexSettings          `json:"settings"`
	Mappings map[string]interface{} `json:"mappings"`
}

// IndexSettings represents the settings section of an Elasticsearch index
type IndexSettings struct {
	NumberOfShards   int                    `json:"number_of_shards"`
	NumberOfReplicas int                    `json:"number_of_replicas"`
	Analysis         map[string]interface{} `json:"analysis,omitempty"`
}

// DefaultDocumentMapping returns the default mapping structure for document indices
func DefaultDocumentMapping() IndexMapping {
	return IndexMapping{
		Settings: IndexSettings{
			NumberOfShards:   1,
			NumberOfReplicas: 1,
			Analysis: map[string]interface{}{
				"analyzer": map[string]interface{}{
					"html_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"char_filter": []string{
							"html_strip",
						},
						"filter": []string{
							"lowercase",
							"stop",
							"snowball",
							"unique",
						},
					},
					"url_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "uax_url_email",
						"filter": []string{
							"lowercase",
							"unique",
						},
					},
					"keyword_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"asciifolding",
							"trim",
						},
					},
				},
				"filter": map[string]interface{}{
					"snowball": map[string]interface{}{
						"type":     "snowball",
						"language": "english",
					},
				},
			},
		},
		Mappings: map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "keyword",
				},
				"url": map[string]interface{}{
					"type": "keyword",
					"fields": map[string]interface{}{
						"text": map[string]interface{}{
							"type":     "text",
							"analyzer": "url_analyzer",
						},
					},
				},
				"title": map[string]interface{}{
					"type":     "text",
					"analyzer": "html_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":     "text",
							"analyzer": "keyword_analyzer",
						},
					},
				},
				"content": map[string]interface{}{
					"type":     "text",
					"analyzer": "html_analyzer",
				},
				"summary": map[string]interface{}{
					"type":     "text",
					"analyzer": "html_analyzer",
				},
				"language": map[string]interface{}{
					"type": "keyword",
				},
				"contentType": map[string]interface{}{
					"type": "keyword",
				},
				"metadata": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
				},
				"pageRank": map[string]interface{}{
					"type": "float",
				},
				"incomingLinks": map[string]interface{}{
					"type": "integer",
				},
				"outgoingLinks": map[string]interface{}{
					"type": "integer",
				},
				"lastCrawled": map[string]interface{}{
					"type": "date",
				},
				"lastModified": map[string]interface{}{
					"type": "date",
				},
				"createdAt": map[string]interface{}{
					"type": "date",
				},
				"keywords": map[string]interface{}{
					"type": "keyword",
					"fields": map[string]interface{}{
						"text": map[string]interface{}{
							"type":     "text",
							"analyzer": "html_analyzer",
						},
					},
				},
				"entities": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "keyword",
							"fields": map[string]interface{}{
								"text": map[string]interface{}{
									"type":     "text",
									"analyzer": "html_analyzer",
								},
							},
						},
						"type": map[string]interface{}{
							"type": "keyword",
						},
						"relevance": map[string]interface{}{
							"type": "float",
						},
					},
				},
				"EnhancedKeywords": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type": "keyword",
							"fields": map[string]interface{}{
								"analyzed": map[string]interface{}{
									"type":     "text",
									"analyzer": "keyword_analyzer",
								},
							},
						},
						"score": map[string]interface{}{
							"type": "float",
						},
						"isDomainSpecific": map[string]interface{}{
							"type": "boolean",
						},
						"category": map[string]interface{}{
							"type": "keyword",
						},
						"position": map[string]interface{}{
							"type": "integer",
						},
					},
				},
				"KeywordTexts": map[string]interface{}{
					"type": "keyword",
					"fields": map[string]interface{}{
						"text": map[string]interface{}{
							"type":     "text",
							"analyzer": "keyword_analyzer",
						},
					},
				},
				"ParsedContent": map[string]interface{}{
					"type":    "object",
					"dynamic": true,
					"properties": map[string]interface{}{
						"keywords": map[string]interface{}{
							"type":    "object",
							"dynamic": true,
							"properties": map[string]interface{}{
								"texts": map[string]interface{}{
									"type": "keyword",
									"fields": map[string]interface{}{
										"analyzed": map[string]interface{}{
											"type":     "text",
											"analyzer": "keyword_analyzer",
										},
									},
								},
								"high_impact": map[string]interface{}{
									"type": "keyword",
									"fields": map[string]interface{}{
										"analyzed": map[string]interface{}{
											"type":     "text",
											"analyzer": "keyword_analyzer",
										},
									},
								},
								"by_domain": map[string]interface{}{
									"type":    "object",
									"dynamic": true,
								},
							},
						},
						"entities_simple": map[string]interface{}{
							"type":    "nested",
							"dynamic": true,
						},
					},
				},
			},
		},
	}
}

// CreateIndexTemplate creates an index template in Elasticsearch
func (c *Client) CreateIndexTemplate(ctx context.Context, name string, mapping IndexMapping, patterns []string) error {
	if !strings.HasPrefix(name, c.indexPrefix) && c.indexPrefix != "" {
		name = c.indexPrefix + "-" + name
	}

	templateBody := map[string]interface{}{
		"index_patterns": patterns,
		"template":       mapping,
	}

	body, err := json.Marshal(templateBody)
	if err != nil {
		return fmt.Errorf("error marshaling index template: %w", err)
	}

	res, err := c.PerformRequest(ctx, &esapi.IndicesPutTemplateRequest{
		Name: name,
		Body: strings.NewReader(string(body)),
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// CreateIndex creates a new index with the given mapping
func (c *Client) CreateIndex(ctx context.Context, indexName string, mapping IndexMapping) error {
	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("error marshaling index mapping: %w", err)
	}
	fullIndexName := c.IndexNameWithPrefix(indexName)
	res, err := c.PerformRequest(ctx, &esapi.IndicesCreateRequest{
		Index: fullIndexName,
		Body:  strings.NewReader(string(body)),
	})
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	return nil
}

// IndexExists checks if an index exists
func (c *Client) IndexExists(ctx context.Context, indexName string) (bool, error) {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	res, err := c.PerformRequest(ctx, &esapi.IndicesExistsRequest{
		Index: []string{fullIndexName},
	})
	if err != nil {
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing response body:", err)
		}
	}(res.Body)

	return res.StatusCode == http.StatusOK, nil
}

// DeleteIndex deletes an index
func (c *Client) DeleteIndex(ctx context.Context, indexName string) error {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	res, err := c.PerformRequest(ctx, &esapi.IndicesDeleteRequest{
		Index: []string{fullIndexName},
	})
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	return nil
}

// RefreshIndex refreshes an index to make recently added documents available for search
func (c *Client) RefreshIndex(ctx context.Context, indexName string) error {
	fullIndexName := c.IndexNameWithPrefix(indexName)
	res, err := c.PerformRequest(ctx, &esapi.IndicesRefreshRequest{
		Index: []string{fullIndexName},
	})
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	return nil
}

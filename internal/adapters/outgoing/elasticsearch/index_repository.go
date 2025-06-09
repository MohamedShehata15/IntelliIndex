package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// IndexRepository implements the outgoing.IndexRepository interface
type IndexRepository struct {
	client         *Client
	refreshTickers map[string]*time.Ticker
	refreshMutex   sync.RWMutex
}

// Ensure IndexRepository implements the outgoing.IndexRepository interface
var _ outgoing.IndexRepository = (*IndexRepository)(nil)

func (i IndexRepository) Create(ctx context.Context, index *domain.Index) error {
	if index == nil {
		return errors.New("index cannot be nil")
	}

	if err := index.Validate(); err != nil {
		return fmt.Errorf("invalid index; %w", err)
	}

	if index.ID == "" {
		index.ID = uuid.NewString()
	}

	settings, mappings, err := i.buildIndexSettings(index)
	if err != nil {
		return fmt.Errorf("error building index settings: %w", err)
	}

	body := map[string]interface{}{
		"settings": settings,
		"mappings": mappings,
	}

	indexName := i.client.IndexNameWithPrefix(index.ID)
	res, err := i.client.PerformRequest(ctx, &esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader(mustMarshalJSON(body)),
	})
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)
	index.UpdateStatus(domain.IndexStatusActive)

	if err := i.saveIndexMetadata(ctx, index); err != nil {
		return fmt.Errorf("error saving index metadata: %w", err)
	}
	backgroundCtx := context.Background()
	if err := i.SetupAutoRefresh(backgroundCtx, index.ID); err != nil {
		fmt.Printf("Warning: Failed to setup automatic refresh for index %s: %v\n", index.ID, err)
	}
	return nil
}

func (i IndexRepository) GetByID(ctx context.Context, id string) (*domain.Index, error) {
	if id == "" {
		return nil, errors.New("index ID cannot be empty")
	}
	index, err := i.getIndexMetadata(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting index metadata: %w", err)
	}
	if index == nil {
		return nil, nil
	}
	count, err := i.getDocumentCount(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting document count: %w", err)
	}
	index.DocumentCount = count
	return index, nil
}

func (i IndexRepository) GetByName(ctx context.Context, name string) (*domain.Index, error) {
	if name == "" {
		return nil, fmt.Errorf("index name cannot be empty")
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"Name.keyword": name,
			},
		},
		"size": 1,
	}
	res, err := i.client.PerformRequest(ctx, &esapi.SearchRequest{
		Index: []string{i.client.IndexNameWithPrefix("indices-metadata")},
		Body:  bytes.NewReader(mustMarshalJSON(query)),
	})
	if err != nil {
		return nil, fmt.Errorf("error searching for index: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	var searchResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("error parsing search response: %w", err)
	}

	hitsObj := searchResult["hits"].(map[string]interface{})
	hitsTotal := int(hitsObj["total"].(map[string]interface{})["value"].(float64))
	if hitsTotal == 0 {
		return nil, nil
	}
	hits := hitsObj["hits"].([]interface{})
	hit := hits[0].(map[string]interface{})
	id := hit["_id"].(string)
	return i.GetByID(ctx, id)
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

func (i *IndexRepository) RefreshIndex(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("index ID cannot be empty")
	}
	exists, err := i.indexExists(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("index with ID %s does not exist", id)
	}

	indexName := i.client.IndexNameWithPrefix(id)
	res, err := i.client.PerformRequest(ctx, &esapi.IndicesRefreshRequest{
		Index: []string{indexName},
	})
	if err != nil {
		return fmt.Errorf("error refreshing index: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	if res.StatusCode >= 400 {
		return fmt.Errorf("error refreshing index: unexpected status code %d", res.StatusCode)
	}

	return nil
}

// NewIndexRepository creates a new Elasticsearch index repository
func NewIndexRepository(client *Client) *IndexRepository {
	return &IndexRepository{
		client:         client,
		refreshTickers: make(map[string]*time.Ticker),
	}
}

// indexExists checks if an index exists in Elasticsearch
func (i *IndexRepository) indexExists(ctx context.Context, id string) (bool, error) {
	indexName := i.client.IndexNameWithPrefix(id)
	res, err := i.client.PerformRequest(ctx, &esapi.IndicesExistsRequest{
		Index: []string{indexName},
	})
	if err != nil {
		return false, fmt.Errorf("error checking if index exists: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("error closing response body:", err)
		}
	}(res.Body)

	return res.StatusCode == 200, nil
}

// StopAutoRefresh stops automatic refreshing for an index
func (i *IndexRepository) StopAutoRefresh(id string) {
	i.refreshMutex.Lock()
	defer i.refreshMutex.Unlock()

	if ticker, exists := i.refreshTickers[id]; exists {
		ticker.Stop()
		delete(i.refreshTickers, id)
	}
}

// SetupAutoRefresh creates a periodic refresh for an index based on its settings
func (i *IndexRepository) SetupAutoRefresh(ctx context.Context, id string) error {
	index, err := i.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting index for auto-refresh: %w", err)
	}

	if index.Settings.RefreshInterval == "" || index.Settings.RefreshInterval == "-1" {
		return nil
	}

	duration, err := time.ParseDuration(index.Settings.RefreshInterval)
	if err != nil {
		return fmt.Errorf("error parsing refresh interval: %w", err)
	}

	i.StopAutoRefresh(id)

	ticker := time.NewTicker(duration)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				refreshCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				if err := i.RefreshIndex(refreshCtx, id); err != nil {
					fmt.Printf("Error auto-refreshing index %s: %v\n", id, err)
				}
				cancel()
			}
		}
	}()

	i.refreshMutex.Lock()
	i.refreshTickers[id] = ticker
	i.refreshMutex.Unlock()

	return nil
}

// ensureMetadataIndexExists creates the indices metadata index if it doesn't exist
func (i *IndexRepository) ensureMetadataIndexExists(ctx context.Context) error {
	indexName := i.client.IndexNameWithPrefix("indices-metadata")
	exists, err := i.client.IndexExists(ctx, "indices-metadata")
	if err != nil {
		return fmt.Errorf("error checking if metadata index exists: %w", err)
	}

	if exists {
		return nil // Index already exists
	}

	body := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"Name": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type":         "keyword",
							"ignore_above": 256,
						},
					},
				},
				"Description":     map[string]interface{}{"type": "text"},
				"Status":          map[string]interface{}{"type": "keyword"},
				"Created":         map[string]interface{}{"type": "date"},
				"LastUpdated":     map[string]interface{}{"type": "date"},
				"Settings":        map[string]interface{}{"type": "object"},
				"DocumentMapping": map[string]interface{}{"type": "object"},
			},
		},
	}

	res, err := i.client.PerformRequest(ctx, &esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader(mustMarshalJSON(body)),
	})
	if err != nil {
		return fmt.Errorf("error creating metadata index: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	return nil
}

// mapToIndex converts an Elasticsearch source map to a domain Index
func (i *IndexRepository) mapToIndex(source map[string]interface{}) (*domain.Index, error) {
	index := &domain.Index{
		Name:            getString(source, "Name"),
		Description:     getString(source, "Description"),
		Status:          domain.IndexStatus(getString(source, "Status")),
		DocumentMapping: make(map[string]string),
	}

	if createdAt, ok := source["CreatedAt"].(string); ok {
		if ts, err := time.Parse(time.RFC3339, createdAt); err == nil {
			index.CreatedAt = ts
		}
	}

	if lastUpdated, ok := source["LastUpdated"].(string); ok {
		if ts, err := time.Parse(time.RFC3339, lastUpdated); err == nil {
			index.LastUpdated = ts
		}
	}

	if settingsMap, ok := source["Settings"].(map[string]interface{}); ok {
		index.Settings.Shards = getIntFromMap(settingsMap, "Shards")
		index.Settings.Replicas = getIntFromMap(settingsMap, "Replicas")
		index.Settings.RefreshInterval = getStringFromMap(settingsMap, "RefreshInterval")

		// Parse stopwords
		if stopwords, ok := settingsMap["Stopwords"].([]interface{}); ok {
			for _, sw := range stopwords {
				if stopword, ok := sw.(string); ok {
					index.Settings.Stopwords = append(index.Settings.Stopwords, stopword)
				}
			}
		}

		// Parse languages
		if languages, ok := settingsMap["Languages"].([]interface{}); ok {
			for _, lang := range languages {
				if language, ok := lang.(string); ok {
					index.Settings.Languages = append(index.Settings.Languages, language)
				}
			}
		}
	}

	if mappingMap, ok := source["DocumentMapping"].(map[string]interface{}); ok {
		for k, v := range mappingMap {
			if valStr, ok := v.(string); ok {
				index.DocumentMapping[k] = valStr
			}
		}
	}

	return index, nil
}

// indexToMap converts a domain Index to a map for Elasticsearch
func (i *IndexRepository) indexToMap(index *domain.Index) (map[string]interface{}, error) {
	indexMap := map[string]interface{}{
		"Name":        index.Name,
		"Description": index.Description,
		"CreatedAt":   index.CreatedAt.Format(time.RFC3339),
		"LastUpdated": index.LastUpdated.Format(time.RFC3339),
		"Status":      string(index.Status),
	}

	settings := map[string]interface{}{
		"Shards":          index.Settings.Shards,
		"Replicas":        index.Settings.Replicas,
		"RefreshInterval": index.Settings.RefreshInterval,
	}

	if len(index.Settings.Stopwords) > 0 {
		settings["Stopwords"] = index.Settings.Stopwords
	}

	if len(index.Settings.Languages) > 0 {
		settings["Languages"] = index.Settings.Languages
	}

	indexMap["Settings"] = settings

	if len(index.DocumentMapping) > 0 {
		indexMap["DocumentMapping"] = index.DocumentMapping
	}

	return indexMap, nil
}

// saveIndexMetadata saves the index metadata to a special indices metadata index
func (i *IndexRepository) saveIndexMetadata(ctx context.Context, index *domain.Index) error {
	if err := i.ensureMetadataIndexExists(ctx); err != nil {
		return fmt.Errorf("error ensuring metadata index exists: %w", err)
	}

	indexMap, err := i.indexToMap(index)
	if err != nil {
		return fmt.Errorf("error converting index to map: %w", err)
	}

	res, err := i.client.PerformRequest(ctx, &esapi.IndexRequest{
		Index:      i.client.IndexNameWithPrefix("indices-metadata"),
		DocumentID: index.ID,
		Body:       bytes.NewReader(mustMarshalJSON(indexMap)),
	})
	if err != nil {
		return fmt.Errorf("error saving index metadata: %w", err)
	}
	defer res.Body.Close()

	return nil
}

// getIndexMetadata retrieves an index's metadata from the metadata index
func (i *IndexRepository) getIndexMetadata(ctx context.Context, id string) (*domain.Index, error) {
	if err := i.ensureMetadataIndexExists(ctx); err != nil {
		return nil, fmt.Errorf("error ensuring metadata index exists: %w", err)
	}

	res, err := i.client.PerformRequest(ctx, &esapi.GetRequest{
		Index:      i.client.IndexNameWithPrefix("indices-metadata"),
		DocumentID: id,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting index metadata: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, nil // Return nil without error for not found
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing metadata response: %w", err)
	}

	source, ok := response["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected metadata response format")
	}

	index, err := i.mapToIndex(source)
	if err != nil {
		return nil, fmt.Errorf("error converting map to index: %w", err)
	}

	index.ID = id

	return index, nil
}

// getDocumentCount retrieves the current document count for an index
func (i *IndexRepository) getDocumentCount(ctx context.Context, id string) (int, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"IndexID.keyword": id,
			},
		},
	}

	res, err := i.client.PerformRequest(ctx, &esapi.CountRequest{
		Index: []string{i.client.IndexNameWithPrefix(DocumentIndex)},
		Body:  bytes.NewReader(mustMarshalJSON(query)),
	})
	if err != nil {
		return 0, fmt.Errorf("error counting documents: %w", err)
	}
	defer res.Body.Close()

	var countResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0, fmt.Errorf("error parsing count response: %w", err)
	}

	count := int(countResult["count"].(float64))

	return count, nil
}

// buildIndexSettings converts domain model settings to Elasticsearch settings
func (i *IndexRepository) buildIndexSettings(index *domain.Index) (map[string]interface{}, map[string]interface{}, error) {
	settings := createBasicSettings(index)
	configureAnalysisSettings(index, settings)
	mappings := createDefaultMappings()
	return settings, mappings, nil
}

// createBasicSettings creates the basic index settings
func createBasicSettings(index *domain.Index) map[string]interface{} {
	settings := map[string]interface{}{
		"number_of_shards":   index.Settings.Shards,
		"number_of_replicas": index.Settings.Replicas,
	}
	if index.Settings.RefreshInterval != "" {
		settings["refresh_interval"] = index.Settings.RefreshInterval
	}
	return settings
}

// configureAnalysisSettings adds language and stopword configurations
func configureAnalysisSettings(index *domain.Index, settings map[string]interface{}) {
	needsAnalysis := len(index.Settings.Stopwords) > 0 || len(index.Settings.Languages) > 0

	if needsAnalysis {
		analysis := map[string]interface{}{}

		if len(index.Settings.Languages) > 0 {
			configureLanguageAnalyzers(index.Settings.Languages, analysis)
		}

		if len(index.Settings.Stopwords) > 0 {
			configureStopwords(index.Settings.Stopwords, analysis)
		}

		settings["analysis"] = analysis
	}
	mergeCustomAnalyzerSettings(index.Settings.AnalyzerSettings, settings)
}

// configureLanguageAnalyzers adds language-specific analyzers
func configureLanguageAnalyzers(languages []string, analysis map[string]interface{}) {
	analyzers := map[string]interface{}{}
	for _, lang := range languages {
		analyzers[lang+"_analyzer"] = map[string]interface{}{
			"type":      "standard",
			"stopwords": "_" + lang + "_",
		}
	}
	analysis["analyzer"] = analyzers
}

// configureStopwords adds custom stopword configuration
func configureStopwords(stopwords []string, analysis map[string]interface{}) {
	analysis["filter"] = map[string]interface{}{
		"custom_stop": map[string]interface{}{
			"type":      "stop",
			"stopwords": stopwords,
		},
	}
}

// mergeCustomAnalyzerSettings adds user-defined analyzer settings
func mergeCustomAnalyzerSettings(customSettings map[string]interface{}, settings map[string]interface{}) {
	if customSettings == nil {
		return
	}
	if analysisSettings, ok := settings["analysis"].(map[string]interface{}); ok {
		for k, v := range customSettings {
			analysisSettings[k] = v
		}
	} else {
		settings["analysis"] = customSettings
	}
}

// createDefaultMappings creates the default index mappings
func createDefaultMappings() map[string]interface{} {
	return map[string]interface{}{
		"properties": map[string]interface{}{
			"Name": map[string]interface{}{
				"type": "text",
				"fields": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type":         "keyword",
						"ignore_above": 256,
					},
				},
			},
			"Description": map[string]interface{}{
				"type": "text",
			},
		},
	}
}

// getString safely extracts a string value from a map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

// getIntFromMap safely extracts an int value from a map
func getIntFromMap(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

// getStringFromMap safely extracts a string value from a map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

package domain

import (
	"errors"
	"time"
)

// Index represents a searchable collection of documents
type Index struct {
	ID              string
	Name            string
	Description     string
	DocumentCount   int
	CreatedAt       time.Time
	LastUpdated     time.Time
	Settings        IndexSettings
	Status          IndexStatus
	DocumentMapping map[string]string
}

// IndexStatus represents the current state of an index
type IndexStatus string

const (
	IndexStatusCreating IndexStatus = "creating"
	IndexStatusActive   IndexStatus = "active"
	IndexStatusUpdating IndexStatus = "updating"
	IndexStatusDeleting IndexStatus = "deleting"
)

// IndexSettings contains configuration for the index
type IndexSettings struct {
	Shards           int
	Replicas         int
	RefreshInterval  string
	AnalyzerSettings map[string]interface{}
	Stopwords        []string
	Languages        []string
}

// NewIndex creates a new index with default settings
func NewIndex(name, description string) (*Index, error) {
	if name == "" {
		return nil, errors.New("index name cannot be empty")
	}
	now := time.Now()
	return &Index{
		Name:        name,
		Description: description,
		CreatedAt:   now,
		LastUpdated: now,
		Status:      IndexStatusCreating,
		Settings: IndexSettings{
			Shards:    1,
			Replicas:  0,
			Languages: []string{"en"},
		},
		DocumentMapping: make(map[string]string),
	}, nil
}

// Validate ensures the index is valid
func (i *Index) Validate() error {
	if i.Name == "" {
		return errors.New("index name cannot be empty")
	}
	return nil
}

// AddDocumentField adds a document field mapping to the index
func (i *Index) AddDocumentField(documentField, indexField string) {
	i.DocumentMapping[documentField] = indexField
	i.LastUpdated = time.Now()
}

// IsActive returns true if the index is in the active state
func (i *Index) IsActive() bool {
	return i.Status == IndexStatusActive
}

// UpdateStatus sets a new status for the index
func (i *Index) UpdateStatus(status IndexStatus) {
	i.Status = status
	i.LastUpdated = time.Now()
}

// IncrementDocumentCount increases the document count by one
// and updates the last updated timestamp
func (i *Index) IncrementDocumentCount() {
	i.DocumentCount++
	i.LastUpdated = time.Now()
}

// DecrementDocumentCount decreases the document count by one
// and updates the last updated timestamp
func (i *Index) DecrementDocumentCount() {
	if i.DocumentCount <= 0 {
		return
	}
	i.DocumentCount--
	i.LastUpdated = time.Now()
}

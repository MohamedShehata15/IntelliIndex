package elasticsearch

import (
	"sync"
	"time"
)

// IndexRepository implements the outgoing.IndexRepository interface
type IndexRepository struct {
	client         *Client
	refreshTickers map[string]*time.Ticker
	refreshMutex   sync.RWMutex
}

// NewIndexRepository creates a new Elasticsearch index repository
func NewIndexRepository(client *Client) *IndexRepository {
	return &IndexRepository{
		client:         client,
		refreshTickers: make(map[string]*time.Ticker),
	}
}

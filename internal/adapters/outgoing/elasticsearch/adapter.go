package elasticsearch

import (
	"github.com/mohamedshehata15/intelli-index/internal/pkg/di"
	"github.com/mohamedshehata15/intelli-index/pkg/config"
)

// ElasticsearchAdapterFactory encapsulates the configuration and registration logic
type ElasticsearchAdapterFactory struct {
	config *config.ElasticConfig
}

var _ di.AdapterRegistrar = (*ElasticsearchAdapterFactory)(nil)

// NewElasticsearchAdapterFactory creates a new factory with the given configuration
func NewElasticsearchAdapterFactory(cfg *config.ElasticConfig) *ElasticsearchAdapterFactory {
	return &ElasticsearchAdapterFactory{
		config: cfg,
	}
}

// Register implements the AdapterRegistrar interface
func (eaf *ElasticsearchAdapterFactory) Register(container *di.Container) error {
	return RegisterElasticsearchAdapters(container, eaf.config)
}

// RegisterElasticsearchAdapters registers all Elasticsearch implementations with the DI container
func RegisterElasticsearchAdapters(container *di.Container, cfg *config.ElasticConfig) error {
	// Create the Elasticsearch client
	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	// Register document repository
	container.Register("documentRepository", func() (interface{}, error) {
		return NewDocumentRepository(client), nil
	})

	// Register index repository
	container.Register("indexRepository", func() (interface{}, error) {
		return NewIndexRepository(client), nil
	})

	return nil
}

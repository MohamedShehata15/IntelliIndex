package webcrawler

import (
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
	"github.com/mohamedshehata15/intelli-index/internal/pkg/di"
	"github.com/mohamedshehata15/intelli-index/pkg/config"
)

// WebCrawlerAdapterFactory encapsulates the configuration and registration logic
type WebCrawlerAdapterFactory struct {
	config *config.CrawlerConfig
}

var _ di.AdapterRegistrar = (*WebCrawlerAdapterFactory)(nil)

func NewWebCrawlerAdapterFactory(cfg *config.CrawlerConfig) *WebCrawlerAdapterFactory {
	return &WebCrawlerAdapterFactory{
		config: cfg,
	}
}

// Register implements the AdapterRegistrar interface
func (w *WebCrawlerAdapterFactory) Register(container *di.Container) error {
	return RegisterWebCrawlerAdapter(container, w.config)
}

// RegisterWebCrawlerAdapter registers the web crawler adapter with the dependency injection container
func RegisterWebCrawlerAdapter(container *di.Container, cfg *config.CrawlerConfig) error {
	// Register web crawler implementation
	container.Register("webCrawler", func() (interface{}, error) {
		return NewWebCrawler(cfg), nil
	})
	return nil
}

// GetWebCrawler retrieves the web crawler implementation from the container
func GetWebCrawler(container *di.Container) outgoing.WebCrawler {
	return container.MustResolve("webCrawler").(outgoing.WebCrawler)
}

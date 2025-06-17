package webcrawler

import (
	"context"

	"github.com/mohamedshehata15/intelli-index/pkg/config"

	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
	"github.com/mohamedshehata15/intelli-index/internal/core/ports/outgoing"
)

// Crawler implements the outgoing.WebCrawler interface
type Crawler struct {
	config *config.CrawlerConfig
}

var _ outgoing.WebCrawler = (*Crawler)(nil)

// NewWebCrawler creates a new instance of Crawler
func NewWebCrawler(config *config.CrawlerConfig) *Crawler {
	return &Crawler{
		config,
	}
}

func (c *Crawler) Crawl(ctx context.Context, url string) (*domain.CrawlResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Crawler) StartCrawlingJob(ctx context.Context, seedURLs []string, options *domain.CrawlOptions) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Crawler) StopCrawlingJob(ctx context.Context, jobID string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Crawler) GetCrawlingJobStatus(ctx context.Context, jobID string) (*domain.CrawlProgress, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Crawler) AddURLToCrawl(ctx context.Context, jobID, url string) error {
	//TODO implement me
	panic("implement me")
}

func (c *Crawler) SetCrawlResultHandler(handler func(ctx context.Context, result *domain.CrawlResult) error) {
	//TODO implement me
	panic("implement me")
}

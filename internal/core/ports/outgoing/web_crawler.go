package outgoing

import (
	"context"
	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
)

// WebCrawler defines the interface for web crawling operations
type WebCrawler interface {
	Crawl(ctx context.Context, url string) (*domain.CrawlResult, error)
	StartCrawlingJob(ctx context.Context, seedURLs []string, options *domain.CrawlOptions) (string, error)
	StopCrawlingJob(ctx context.Context, jobID string) error
	GetCrawlingJobStatus(ctx context.Context, jobID string) (*domain.CrawlProgress, error)
	AddURLToCrawl(ctx context.Context, jobID, url string) error
	SetCrawlResultHandler(handler func(ctx context.Context, result *domain.CrawlResult) error)
}

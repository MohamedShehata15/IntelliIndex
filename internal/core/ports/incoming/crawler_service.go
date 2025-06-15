package incoming

import (
	"context"
	"github.com/mohamedshehata15/intelli-index/internal/core/domain"
)

// CrawlerService defines the primary port for web crawling functionality
type CrawlerService interface {
	StartCrawl(ctx context.Context, seedURLs []string, indexID string, options map[string]interface{}) (*domain.CrawlJob, error)
	StopCrawl(ctx context.Context, jobID string) error
	GetCrawlStatus(ctx context.Context, jobID string) (*domain.CrawlProgress, error)
	ListCrawlJobs(ctx context.Context) ([]*domain.CrawlJob, error)
	AddURLToCrawl(ctx context.Context, jobID string, url string) error
	ReindexDocument(ctx context.Context, url string, indexID string)
}

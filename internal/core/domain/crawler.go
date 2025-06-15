package domain

// CrawlStatus represents the status of a crawl job
type CrawlStatus string

const (
	CrawlStatusPending   CrawlStatus = "pending"
	CrawlStatusRunning   CrawlStatus = "running"
	CrawlStatusCompleted CrawlStatus = "completed"
	CrawlStatusFailed    CrawlStatus = "failed"
	CrawlStatusCancelled CrawlStatus = "cancelled"
)

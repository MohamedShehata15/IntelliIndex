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

// CrawlJob represents a web crawling job
type CrawlJob struct {
	ID             string
	SeedURLs       []string
	MaxDepth       int
	MaxURLs        int
	AllowedDomains []string
	StartedAt      int64
	CompletedAt    int64
	Status         CrawlStatus
	DocumentCount  int
	ErrorCount     int
	IndexID        string
	CreatedAt      int64
}

// CrawlProgress represents the current progress of a crawling job
type CrawlProgress struct {
	JobID          string
	Status         CrawlStatus
	ProcessedURLs  int
	DiscoveredURLs int
	CurrentDepth   int
	ErrorCount     int
	Errors         []string
	LastUpdated    int64
}

// MonitoringOptions contains configuration for crawler monitoring
type MonitoringOptions struct {
	LogLevel           string
	EnableFileLogging  bool
	MonitoringInterval int
	EnableTracing      bool
	AlertThresholds    map[string]interface{}
	EnableAlerts       bool
	ExternalLogging    bool
	ExternalEndpoint   string
	NotificationEmail  string
	ExportMetrics      bool
}

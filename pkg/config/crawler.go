package config

import "time"

// CrawlerConfig contains web crawler configuration
type CrawlerConfig struct {
	MaxDepth        int
	MaxURLs         int
	RequestDelay    time.Duration
	Timeout         time.Duration
	RespectRobotsTx bool
	UserAgent       string
	AllowedDomains  []string
	ExcludedPaths   []string
}

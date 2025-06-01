package elasticsearch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"github.com/mohamedshehata15/intelli-index/pkg/config"
)

// Client wraps an Elasticsearch client with additional functionality
type Client struct {
	es           *elasticsearch.Client
	indexPrefix  string
	retryBackoff time.Duration
	maxRetries   int
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// NewClient creates a new Elasticsearch client with the provided configuration
func NewClient(cfg *config.ElasticConfig, options ...ClientOption) (*Client, error) {
	client := initializeClient(cfg)
	applyClientOptions(client, options)

	es, err := createElasticsearchClient(cfg, client)
	if err != nil {
		return nil, err
	}

	client.es = es
	return client, nil
}

// initializeClient creates a new client with default settings
func initializeClient(cfg *config.ElasticConfig) *Client {
	return &Client{
		indexPrefix:  cfg.IndexPrefix,
		retryBackoff: 200 * time.Millisecond,
		maxRetries:   3,
	}
}

// applyClientOptions applies the provided options to the client
func applyClientOptions(client *Client, options []ClientOption) {
	for _, option := range options {
		option(client)
	}
}

// createElasticsearchClient configures and creates the Elasticsearch client
func createElasticsearchClient(cfg *config.ElasticConfig, client *Client) (*elasticsearch.Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.URL},
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: cfg.Timeout,
		},
		RetryOnStatus: []int{http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusTooManyRequests},
		MaxRetries:    client.maxRetries,
		RetryBackoff:  func(i int) time.Duration { return client.retryBackoff },
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	return es, nil
}

// WithMaxRetries sets the maximum number of retries for Elasticsearch operations
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithRetryBackoff sets the retry backoff duration for Elasticsearch operations
func WithRetryBackoff(backoff time.Duration) ClientOption {
	return func(c *Client) {
		c.retryBackoff = backoff
	}
}

// WithIndexPrefix sets the index prefix for Elasticsearch operations
func WithIndexPrefix(prefix string) ClientOption {
	return func(c *Client) {
		c.indexPrefix = prefix
	}
}

// Ping checks if the Elasticsearch cluster is available
func (c *Client) Ping(ctx context.Context) (bool, error) {
	res, err := c.es.Ping(
		c.es.Ping.WithContext(ctx),
		c.es.Ping.WithHuman(),
		c.es.Ping.WithPretty(),
	)
	if err != nil {
		return false, fmt.Errorf("error pinging Elasticsearch: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	return res.IsError() == false, nil
}

// GetClient returns the underlying Elasticsearch client
func (c *Client) GetClient() *elasticsearch.Client {
	return c.es
}

// IndexNameWithPrefix returns the index name with the configured prefix
func (c *Client) IndexNameWithPrefix(indexName string) string {
	return fmt.Sprintf("%s-%s", c.indexPrefix, indexName)
}

// PerformRequest performs an Elasticsearch API request and handles errors
func (c *Client) PerformRequest(ctx context.Context, req esapi.Request) (*esapi.Response, error) {
	res, err := req.Do(ctx, c.es)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}

	if res.IsError() {
		return res, fmt.Errorf("elasticsearch responded with error: %s", res.String())
	}

	return res, nil
}

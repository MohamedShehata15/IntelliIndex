package config

import "time"

// ElasticConfig represents Elasticsearch configuration settings
type ElasticConfig struct {
	URL              string        `mapstructure:"url" yaml:"url" default:"http://localhost:9200"`
	Username         string        `mapstructure:"username" yaml:"username" default:""`
	Password         string        `mapstructure:"password" yaml:"password" default:""`
	IndexPrefix      string        `mapstructure:"index_prefix" yaml:"index_prefix" default:"search-engine"`
	Timeout          time.Duration `mapstructure:"timeout" yaml:"timeout" default:"30s"`
	MaxRetries       int           `mapstructure:"max_retries" yaml:"max_retries" default:"3"`
	RetryBackoff     time.Duration `mapstructure:"retry_backoff" yaml:"retry_backoff" default:"200ms"`
	SnippetSize      int           `mapstructure:"snippet_size" yaml:"snippet_size" default:"160"`
	BulkSize         int           `mapstructure:"bulk_size" yaml:"bulk_size" default:"100"`
	RefreshInterval  string        `mapstructure:"refresh_interval" yaml:"refresh_interval" default:"1s"`
	NumberOfShards   int           `mapstructure:"number_of_shards" yaml:"number_of_shards" default:"1"`
	NumberOfReplicas int           `mapstructure:"number_of_replicas" yaml:"number_of_replicas" default:"1"`
	DefaultLanguage  string        `mapstructure:"default_language" yaml:"default_language" default:"english"`
}

// Validate checks if the Elasticsearch configuration is valid
func (e *ElasticConfig) Validate() error {
	if e.URL == "" {
		return &ValidationError{"Elasticsearch URL cannot be empty"}
	}
	return nil
}

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	Elastic  ElasticConfig
	Database DBConfig
}

// ValidationError is returned when configuration validation fails
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("config validation error: %s", e.Message)
}

// Validate validates that the configuration is correct
func (c *Config) Validate() error {
	if err := c.Elastic.Validate(); err != nil {
		return err
	}
	if err := c.Database.Validate(); err != nil {
		return err
	}
	return nil
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	applyEnvironmentOverrides(&config)

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadEnvFile attempts to load environment variables from the specified file or fallback locations
// Returns true if any env file was successfully loaded, false otherwise
func LoadEnvFile(envFile string) (bool, string, error) {
	if envFile != "" {
		if err := godotenv.Load(envFile); err == nil {
			return true, envFile, nil
		}
	}

	possibleLocations := []string{
		".env",
		"../.env",
		"../../.env",
		filepath.Join("cmd", "api", ".env"),
	}

	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		possibleLocations = append(possibleLocations,
			filepath.Join(execDir, ".env"),
			filepath.Join(execDir, "..", ".env"),
		)
	}

	if workDir, err := os.Getwd(); err == nil {
		possibleLocations = append(possibleLocations,
			filepath.Join(workDir, ".env"),
			filepath.Join(workDir, "..", ".env"),
			filepath.Join(workDir, "..", "..", ".env"),
		)
	}

	for _, path := range possibleLocations {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err == nil {
				return true, path, nil
			}
		}
	}

	return false, "", nil
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:            getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
			LogLevel:        getEnvStr("SERVER_LOG_LEVEL", "info"),
			AllowOrigins:    getEnvStringSlice("SERVER_ALLOW_ORIGINS", []string{"*"}),
		},
		Elastic: ElasticConfig{
			URL:              getEnvStr("ELASTICSEARCH_URL", "http://elasticsearch:9200"),
			IndexPrefix:      getEnvStr("ELASTICSEARCH_INDEX_PREFIX", "search-engine"),
			Username:         getEnvStr("ELASTICSEARCH_USERNAME", ""),
			Password:         getEnvStr("ELASTICSEARCH_PASSWORD", ""),
			Timeout:          getEnvDuration("ELASTICSEARCH_TIMEOUT", 5*time.Second),
			MaxRetries:       getEnvInt("ELASTICSEARCH_MAX_RETRIES", 3),
			RetryBackoff:     getEnvDuration("ELASTICSEARCH_RETRY_BACKOFF", 200*time.Millisecond),
			SnippetSize:      getEnvInt("ELASTICSEARCH_SNIPPET_SIZE", 160),
			BulkSize:         getEnvInt("ELASTICSEARCH_BULK_SIZE", 100),
			RefreshInterval:  getEnvStr("ELASTICSEARCH_REFRESH_INTERVAL", "1s"),
			NumberOfShards:   getEnvInt("ELASTICSEARCH_NUMBER_OF_SHARDS", 1),
			NumberOfReplicas: getEnvInt("ELASTICSEARCH_NUMBER_OF_REPLICAS", 1),
			DefaultLanguage:  getEnvStr("ELASTICSEARCH_DEFAULT_LANGUAGE", "english"),
		},
		Database: DBConfig{
			Type:          getEnvStr("DB_TYPE", "sqlite"),
			Host:          getEnvStr("DB_HOST", "localhost"),
			Port:          getEnvInt("DB_PORT", 5432),
			Username:      getEnvStr("DB_USERNAME", "postgres"),
			Password:      getEnvStr("DB_PASSWORD", "postgres"),
			Database:      getEnvStr("DB_DATABASE", "search_engine"),
			SSLMode:       getEnvStr("DB_SSLMODE", "disable"),
			MaxOpenConns:  getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:  getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLife:   getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			AutoMigrate:   getEnvBool("DB_AUTO_MIGRATE", true),
			LogLevel:      getEnvStr("DB_LOG_LEVEL", "error"),
			MigrationPath: getEnvStr("DB_MIGRATION_PATH", "migrations"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// applyEnvironmentOverrides applies environment variable values over the loaded config
func applyEnvironmentOverrides(config *Config) {
	// Elasticsearch overrides
	if val := os.Getenv("ELASTICSEARCH_URL"); val != "" {
		config.Elastic.URL = val
	}
	if val := os.Getenv("ELASTICSEARCH_INDEX_PREFIX"); val != "" {
		config.Elastic.IndexPrefix = val
	}
	if val := os.Getenv("ELASTICSEARCH_USERNAME"); val != "" {
		config.Elastic.Username = val
	}
	if val := os.Getenv("ELASTICSEARCH_PASSWORD"); val != "" {
		config.Elastic.Password = val
	}

	// Database overrides
	if val := os.Getenv("DB_TYPE"); val != "" {
		config.Database.Type = val
	}
	if val := os.Getenv("DB_HOST"); val != "" {
		config.Database.Host = val
	}
}

func getEnvStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func getEnvStringSlice(key string, fallback []string) []string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return strings.Split(value, ",")
	}
	return fallback
}

// FindConfigFile looks for config files in common locations
func FindConfigFile(configPath string) string {
	if configPath != "" {
		return configPath
	}

	locations := []string{
		"./config.yaml",
		"./config.yml",
		"./configs/config.yaml",
		"./configs/config.yml",
		"../configs/config.yaml",
		"../configs/config.yml",
	}

	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		locations = append(locations,
			filepath.Join(execDir, "config.yaml"),
			filepath.Join(execDir, "config.yml"),
			filepath.Join(execDir, "configs", "config.yaml"),
			filepath.Join(execDir, "configs", "config.yml"),
		)
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	return ""
}

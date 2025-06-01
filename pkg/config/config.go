package config

import "fmt"

// Config represents the application configuration
type Config struct {
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

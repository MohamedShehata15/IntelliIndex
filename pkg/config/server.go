package config

import "time"

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port            int           `mapstructure:"port" yaml:"port" default:"8080"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" yaml:"read_timeout" default:"15s"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" yaml:"write_timeout" default:"15s"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout" default:"5s"`
	LogLevel        string        `mapstructure:"log_level" yaml:"log_level" default:"info"`
	AllowOrigins    []string      `mapstructure:"allow_origins" yaml:"allow_origins" default:"*"`
}

// Validate checks if the server configuration is valid
func (s *ServerConfig) Validate() error {
	if s.Port < 0 || s.Port > 65535 {
		return &ValidationError{Message: "server port must be between 0 and 65535"}
	}
	return nil
}

// DefaultServerConfig returns the default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:            8080,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}
}

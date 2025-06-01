package config

import (
	"fmt"
	"time"
)

// DBConfig contains database connection configuration
type DBConfig struct {
	Type          string        `mapstructure:"type" yaml:"type" default:"postgres"`
	Host          string        `mapstructure:"host" yaml:"host" default:"localhost"`
	Port          int           `mapstructure:"port" yaml:"port" default:"5432"`
	Username      string        `mapstructure:"username" yaml:"username" default:""`
	Password      string        `mapstructure:"password" yaml:"password" default:""`
	Database      string        `mapstructure:"database" yaml:"database" default:""`
	SSLMode       string        `mapstructure:"ssl_mode" yaml:"ssl_mode" default:"disable"`
	MaxOpenConns  int           `mapstructure:"max_open_conns" yaml:"max_open_conns" default:"10"`
	MaxIdleConns  int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns" default:"5"`
	ConnMaxLife   time.Duration `mapstructure:"conn_max_life" yaml:"conn_max_life" default:"1h"`
	AutoMigrate   bool          `mapstructure:"auto_migrate" yaml:"auto_migrate" default:"false"`
	LogLevel      string        `mapstructure:"log_level" yaml:"log_level" default:"error"`
	MigrationPath string        `mapstructure:"migration_path" yaml:"migration_path" default:"migrations"`
}

// DSN returns the database connection string for the specified database type
func (c *DBConfig) DSN() string {
	switch c.Type {
	case "postgresql", "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
			c.Username, c.Password, c.Host, c.Port, c.Database)
	case "sqlite":
		return c.Database
	default:
		return ""
	}
}

// Validate checks if the database configuration is valid
func (c *DBConfig) Validate() error {
	if c.Type == "" {
		return &ValidationError{"database type cannot be empty"}
	}

	supportedTypes := map[string]bool{
		"postgresql": true,
		"postgres":   true,
		"mysql":      true,
		"sqlite":     true,
	}

	if !supportedTypes[c.Type] {
		return &ValidationError{fmt.Sprintf("unsupported database type: %s", c.Type)}
	}

	if c.Type != "sqlite" {
		if c.Host == "" {
			return &ValidationError{"database host cannot be empty"}
		}

		if c.Port <= 0 || c.Port > 65535 {
			return &ValidationError{"database port must be between 1 and 65535"}
		}

		if c.Database == "" {
			return &ValidationError{"database name cannot be empty"}
		}
	} else {
		if c.Database == "" {
			return &ValidationError{"database file path cannot be empty for SQLite"}
		}
	}

	return nil
}

package storage

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mohamedshehata15/intelli-index/pkg/config"
)

// Client Wraps a GORM database client
type Client struct {
	DB     *gorm.DB
	Config *config.DBConfig
}

// NewClient creates a new GORM database client based on the provided configuration
func NewClient(cfg *config.DBConfig) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	gormConfig := setupGormLogger(cfg)

	db, err := connectToDatabase(cfg, gormConfig)
	if err != nil {
		return nil, err
	}

	if err := configureConnectionPool(db, cfg); err != nil {
		return nil, err
	}

	return &Client{
		DB:     db,
		Config: cfg,
	}, nil
}

// setupGormLogger configures the GORM logger based on the provided configuration
func setupGormLogger(cfg *config.DBConfig) *gorm.Config {
	logLevel := logger.Error // Default to error only
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	return &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}
}

// connectToDatabase establishes a connection to the specified database type
func connectToDatabase(cfg *config.DBConfig, gormConfig *gorm.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.Type {
	case "postgresql", "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.DSN()), gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// configureConnectionPool sets up the database connection pool parameters
func configureConnectionPool(db *gorm.DB, cfg *config.DBConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLife)

	return nil
}

// Ping checks if the database connection is healthy
func (c *Client) Ping() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Transaction executes the given function within a database transaction
func (c *Client) Transaction(fn func(tx *gorm.DB) error) error {
	return c.DB.Transaction(fn)
}

// Close closes the database connection
func (c *Client) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}

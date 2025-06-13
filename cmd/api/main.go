package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpAdapter "github.com/mohamedshehata15/intelli-index/internal/adapters/incoming/http"
	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/elasticsearch"
	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/storage"
	"github.com/mohamedshehata15/intelli-index/internal/pkg/di"
	"github.com/mohamedshehata15/intelli-index/pkg/config"
)

func main() {
	fmt.Println("Starting IntelliIndex API...")

	// Command line flags
	configPath := flag.String("config", "", "Path to configuration file (if not specified, environment variables will be used)")
	envFile := flag.String("env", ".env", "Path to .env file for environment variables")
	runMigrations := flag.Bool("migrate", false, "Run database migrations")
	resetDB := flag.Bool("reset-db", false, "Reset database (drop all tables and recreate)")
	flag.Parse()

	// Load configuration
	cfg, err := loadConfig(*configPath, *envFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully (Elasticsearch: %s, Database: %s)",
		cfg.Elastic.URL, cfg.Database.Type)

	container := di.Bootstrap()
	err = di.BatchRegister(container,
		elasticsearch.NewElasticsearchAdapterFactory(&cfg.Elastic),
		storage.NewStorageAdapterFactory(&cfg.Database),
	)
	if err != nil {
		log.Fatalf("Failed to register adapters: %v", err)
	}

	// Database Migration via flags or config
	if *runMigrations || cfg.Database.AutoMigrate {
		log.Println("Running database migrations based on configuration or command line flag")
		migrationHandler := storage.GetMigrationHandler(container)

		if *resetDB {
			if err := migrationHandler.ResetDatabase(); err != nil {
				log.Fatalf("Failed to reset database: %v", err)
			}
			log.Println("Database reset successfully")
		} else {
			if err := migrationHandler.RunMigrations(); err != nil {
				log.Fatalf("Failed to run database migrations: %v", err)
			}
			log.Println("Database migrations completed successfully")
		}
	}

	// Create HTTP server
	server := httpAdapter.NewServer()
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: server.Handler(),
	}
	go func() {
		log.Printf("Starting server on %s\n", serverAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

}

// loadConfig loads configuration from file or environment variables
func loadConfig(configPath, envFile string) (*config.Config, error) {
	loaded, loadedPath, _ := config.LoadEnvFile(envFile)
	if loaded {
		log.Printf("Loaded environment variables from %s", loadedPath)
	} else {
		log.Println("No .env file found or loaded, using system environment variables only")
	}

	if configPath != "" {
		log.Printf("Loading configuration from file: %s", configPath)
		return config.LoadFromFile(configPath)
	}

	if foundConfigPath := config.FindConfigFile(""); foundConfigPath != "" {
		log.Printf("Found configuration file at: %s", foundConfigPath)
		return config.LoadFromFile(foundConfigPath)
	}

	log.Println("No configuration file found, loading from environment variables")
	return config.Load()
}

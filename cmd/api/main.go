package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/mohamedshehata15/intelli-index/pkg/config"

	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/elasticsearch"
	"github.com/mohamedshehata15/intelli-index/internal/adapters/outgoing/storage"
	"github.com/mohamedshehata15/intelli-index/internal/pkg/di"
)

func main() {
	fmt.Println("Starting IntelliIndex API...")

	// Command line flags
	configPath := flag.String("config", "", "Path to configuration file (if not specified, environment variables will be used)")
	envFile := flag.String("env", ".env", "Path to .env file for environment variables")
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

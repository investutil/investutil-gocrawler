package main

import (
    "context"
    "flag"
    "log"
    "time"

    "github.com/yourusername/investutil-gocrawler/internal/common/config"
    "github.com/yourusername/investutil-gocrawler/internal/crawler/crypto"
    "github.com/yourusername/investutil-gocrawler/internal/storage"
)

type Config struct {
    Storage struct {
        MongoDB storage.MongoDBConfig `yaml:"mongodb"`
    } `yaml:"storage"`
    Crawler crypto.Config `yaml:"crawler"`
}

func main() {
    commonConfig := flag.String("common-config", "configs/common.yaml", "path to common config file")
    specificConfig := flag.String("config", "configs/btchistory.yaml", "path to specific config file")
    flag.Parse()

    // Load configs
    var cfg Config
    if err := config.LoadCommonConfig(*commonConfig, *specificConfig, &cfg); err != nil {
        log.Fatalf("Failed to load configs: %v", err)
    }

    // Initialize MongoDB storage
    mongoStorage, err := storage.NewMongoDBStorage(cfg.Storage.MongoDB)
    if err != nil {
        log.Fatalf("Failed to initialize MongoDB storage: %v", err)
    }
    defer func() {
        if err := mongoStorage.Close(context.Background()); err != nil {
            log.Printf("Failed to close MongoDB connection: %v", err)
        }
    }()

    // Initialize crawler
    crawler := crypto.NewBitcoinCrawler(mongoStorage, &cfg.Crawler)

    // Run crawler
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    if err := crawler.Crawl(ctx); err != nil {
        log.Fatalf("Crawler failed: %v", err)
    }

    log.Printf("Crawler %s completed successfully", crawler.Name())
} 
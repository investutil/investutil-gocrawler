package main

import (
    "context"
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/yourusername/investutil-gocrawler/internal/collector"
    "github.com/yourusername/investutil-gocrawler/internal/config"
    "github.com/yourusername/investutil-gocrawler/internal/database"
    "github.com/yourusername/investutil-gocrawler/internal/queue"
)

func main() {
    configPath := flag.String("config", "config.yaml", "path to config file")
    mode := flag.String("mode", "collect", "operation mode: collect or process")
    flag.Parse()

    // Load configuration
    cfg, err := config.Load(*configPath)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize MongoDB
    db, err := database.NewMongoDB(cfg.Database.MongoDB)
    if err != nil {
        log.Fatalf("Failed to initialize MongoDB: %v", err)
    }
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := db.Close(ctx); err != nil {
            log.Printf("Failed to close MongoDB connection: %v", err)
        }
    }()

    // Initialize RabbitMQ
    rmq, err := queue.NewRabbitMQ(cfg.Queue.RabbitMQ)
    if err != nil {
        log.Fatalf("Failed to initialize RabbitMQ: %v", err)
    }
    defer rmq.Close()

    // Initialize Bitcoin collector
    btcCollector := collector.NewBitcoinCollector(db, rmq, cfg.Collector.Schedule)

    // Setup signal handling
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()

    // Run in specified mode
    switch *mode {
    case "collect":
        log.Printf("Starting collector: %s", btcCollector.Name())
        if err := btcCollector.Collect(ctx); err != nil {
            log.Fatalf("Collector failed: %v", err)
        }
        log.Printf("Collector %s completed successfully", btcCollector.Name())

    case "process":
        log.Printf("Starting processor: %s", btcCollector.Name())
        if err := btcCollector.Process(ctx); err != nil {
            log.Fatalf("Processor failed: %v", err)
        }
        log.Printf("Processor %s completed successfully", btcCollector.Name())

    default:
        log.Fatalf("Unknown mode: %s", *mode)
    }
} 
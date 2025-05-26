package database

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/yourusername/investutil-gocrawler/internal/models"
)

// Database defines the interface for database operations
type Database interface {
    // SaveBitcoinPrices saves bitcoin price data
    SaveBitcoinPrices(ctx context.Context, data models.BitcoinDailyData) error
    // Close closes the database connection
    Close(ctx context.Context) error
}

// MongoDB implements Database interface
type MongoDB struct {
    client   *mongo.Client
    database string
}

// Config holds MongoDB configuration
type Config struct {
    URI      string `yaml:"uri"`
    Database string `yaml:"database"`
}

// NewMongoDB creates a new MongoDB instance
func NewMongoDB(cfg Config) (*MongoDB, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    // Ping the database to verify connection
    if err := client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    return &MongoDB{
        client:   client,
        database: cfg.Database,
    }, nil
}

// SaveBitcoinPrices implements Database.SaveBitcoinPrices
func (m *MongoDB) SaveBitcoinPrices(ctx context.Context, data models.BitcoinDailyData) error {
    collection := m.client.Database(m.database).Collection("bitcoin_prices")
    
    _, err := collection.InsertOne(ctx, data)
    if err != nil {
        return fmt.Errorf("failed to insert bitcoin prices: %w", err)
    }

    return nil
}

// Close implements Database.Close
func (m *MongoDB) Close(ctx context.Context) error {
    if err := m.client.Disconnect(ctx); err != nil {
        return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
    }
    return nil
} 
package collector

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/yourusername/investutil-gocrawler/internal/database"
    "github.com/yourusername/investutil-gocrawler/internal/models"
    "github.com/yourusername/investutil-gocrawler/internal/queue"
)

// Collector defines the interface for data collectors
type Collector interface {
    // Name returns the collector's name
    Name() string
    // Collect performs the data collection
    Collect(ctx context.Context) error
    // Process processes the collected data
    Process(ctx context.Context) error
    // Schedule returns the collector's schedule in cron format
    Schedule() string
}

// BaseCollector provides common functionality for collectors
type BaseCollector struct {
    name     string
    schedule string
    db       database.Database
    queue    *queue.RabbitMQ
}

// NewBaseCollector creates a new BaseCollector
func NewBaseCollector(name, schedule string, db database.Database, queue *queue.RabbitMQ) *BaseCollector {
    return &BaseCollector{
        name:     name,
        schedule: schedule,
        db:       db,
        queue:    queue,
    }
}

// Name implements Collector.Name
func (b *BaseCollector) Name() string {
    return b.name
}

// Schedule implements Collector.Schedule
func (b *BaseCollector) Schedule() string {
    return b.schedule
}

// BitcoinCollector implements bitcoin price data collector
type BitcoinCollector struct {
    *BaseCollector
    client *http.Client
}

// NewBitcoinCollector creates a new BitcoinCollector
func NewBitcoinCollector(db database.Database, queue *queue.RabbitMQ, schedule string) *BitcoinCollector {
    return &BitcoinCollector{
        BaseCollector: NewBaseCollector("bitcoin-price", schedule, db, queue),
        client: &http.Client{
            Timeout: time.Second * 30,
        },
    }
}

// Collect implements Collector.Collect for BitcoinCollector
func (c *BitcoinCollector) Collect(ctx context.Context) error {
    url := "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart?vs_currency=usd&days=max&interval=daily"
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to fetch data: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var geckoResp models.CoinGeckoResponse
    if err := json.NewDecoder(resp.Body).Decode(&geckoResp); err != nil {
        return fmt.Errorf("failed to decode response: %w", err)
    }

    // Convert response to our data model
    var prices []models.BitcoinPrice
    for i := 0; i < len(geckoResp.Prices); i++ {
        timestamp := time.UnixMilli(int64(geckoResp.Prices[i][0]))
        price := models.BitcoinPrice{
            Timestamp: timestamp,
            Price:     geckoResp.Prices[i][1],
            MarketCap: geckoResp.MarketCaps[i][1],
            Volume24h: geckoResp.TotalVolumes[i][1],
        }
        prices = append(prices, price)
    }

    // Create data package
    data := models.BitcoinDailyData{
        LastUpdated: time.Now().UTC(),
        Data:        prices,
    }

    // Serialize data
    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal data: %w", err)
    }

    // Send to queue
    if err := c.queue.Publish(ctx, jsonData); err != nil {
        return fmt.Errorf("failed to publish data: %w", err)
    }

    return nil
}

// Process implements Collector.Process for BitcoinCollector
func (c *BitcoinCollector) Process(ctx context.Context) error {
    return c.queue.Consume(ctx, func(data []byte) error {
        var priceData models.BitcoinDailyData
        if err := json.Unmarshal(data, &priceData); err != nil {
            return fmt.Errorf("failed to unmarshal data: %w", err)
        }

        if err := c.db.SaveBitcoinPrices(ctx, priceData); err != nil {
            return fmt.Errorf("failed to save data: %w", err)
        }

        return nil
    })
} 
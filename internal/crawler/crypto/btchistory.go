package crypto

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/yourusername/investutil-gocrawler/internal/crawler"
    "github.com/yourusername/investutil-gocrawler/internal/models"
    "github.com/yourusername/investutil-gocrawler/internal/storage"
)

const (
    coinGeckoBaseURL = "https://api.coingecko.com/api/v3"
    bitcoinID        = "bitcoin"
)

// BitcoinCrawler implements bitcoin price data crawler
type BitcoinCrawler struct {
    *crawler.BaseCrawler
    storage storage.Storage
    client  *http.Client
    config  *Config
}

// Config holds configuration for BitcoinCrawler
type Config struct {
    DataPath string `yaml:"data_path"`
    Schedule string `yaml:"schedule"`
}

// NewBitcoinCrawler creates a new BitcoinCrawler instance
func NewBitcoinCrawler(storage storage.Storage, config *Config) *BitcoinCrawler {
    return &BitcoinCrawler{
        BaseCrawler: crawler.NewBaseCrawler("bitcoin-history", config.Schedule),
        storage:     storage,
        client: &http.Client{
            Timeout: time.Second * 30,
        },
        config: config,
    }
}

// Crawl implements the main crawling logic
func (c *BitcoinCrawler) Crawl(ctx context.Context) error {
    // Fetch data from CoinGecko
    url := fmt.Sprintf("%s/coins/%s/market_chart?vs_currency=usd&days=max&interval=daily",
        coinGeckoBaseURL, bitcoinID)
    
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

    // Prepare data for storage
    data := models.BitcoinDailyData{
        LastUpdated: time.Now().UTC(),
        Data:        prices,
    }

    // Save to storage
    key := fmt.Sprintf("%s/latest.json", c.config.DataPath)
    if err := c.storage.Save(ctx, key, data); err != nil {
        return fmt.Errorf("failed to save data: %w", err)
    }

    // Save yearly data
    year := time.Now().Format("2006")
    yearlyKey := fmt.Sprintf("%s/%s/btc-%s.json", c.config.DataPath, year, year)
    if err := c.storage.Save(ctx, yearlyKey, data); err != nil {
        return fmt.Errorf("failed to save yearly data: %w", err)
    }

    // Update last run time
    c.UpdateLastRun()
    return nil
}

// Name returns the crawler name
func (c *BitcoinCrawler) Name() string {
    return "bitcoin-history"
} 
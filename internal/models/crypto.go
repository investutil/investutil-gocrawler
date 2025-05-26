package models

import "time"

// BitcoinPrice represents a single bitcoin price data point
type BitcoinPrice struct {
    Timestamp time.Time `json:"timestamp"`
    Price     float64   `json:"price"`
    Volume24h float64   `json:"volume_24h"`
    MarketCap float64   `json:"market_cap"`
}

// BitcoinDailyData represents a collection of daily bitcoin price data
type BitcoinDailyData struct {
    LastUpdated time.Time      `json:"last_updated"`
    Data        []BitcoinPrice `json:"data"`
}

// CoinGeckoResponse represents the response from CoinGecko API
type CoinGeckoResponse struct {
    Prices       [][2]float64 `json:"prices"`
    MarketCaps   [][2]float64 `json:"market_caps"`
    TotalVolumes [][2]float64 `json:"total_volumes"`
} 
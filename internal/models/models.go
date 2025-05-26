package models

import (
    "time"
)

// BitcoinPrice represents a single bitcoin price data point
type BitcoinPrice struct {
    Timestamp time.Time `json:"timestamp" bson:"timestamp"`
    Price     float64   `json:"price" bson:"price"`
    MarketCap float64   `json:"market_cap" bson:"market_cap"`
    Volume24h float64   `json:"volume_24h" bson:"volume_24h"`
}

// BitcoinDailyData represents a collection of bitcoin price data
type BitcoinDailyData struct {
    LastUpdated time.Time      `json:"last_updated" bson:"last_updated"`
    Data        []BitcoinPrice `json:"data" bson:"data"`
}

// CoinGeckoResponse represents the response from CoinGecko API
type CoinGeckoResponse struct {
    Prices       [][2]float64 `json:"prices"`
    MarketCaps   [][2]float64 `json:"market_caps"`
    TotalVolumes [][2]float64 `json:"total_volumes"`
} 
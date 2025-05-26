package crawler

import (
    "context"
    "time"
)

// Crawler defines the interface that all crawlers must implement
type Crawler interface {
    // Name returns the crawler's name
    Name() string
    
    // Crawl performs the data crawling
    Crawl(ctx context.Context) error
    
    // Schedule returns the crawler's schedule in cron format
    Schedule() string
    
    // LastRun returns the last successful run time
    LastRun() time.Time
}

// BaseCrawler provides common functionality for crawlers
type BaseCrawler struct {
    name     string
    schedule string
    lastRun  time.Time
}

// NewBaseCrawler creates a new BaseCrawler
func NewBaseCrawler(name, schedule string) *BaseCrawler {
    return &BaseCrawler{
        name:     name,
        schedule: schedule,
    }
}

// Name implements Crawler.Name
func (b *BaseCrawler) Name() string {
    return b.name
}

// Schedule implements Crawler.Schedule
func (b *BaseCrawler) Schedule() string {
    return b.schedule
}

// LastRun implements Crawler.LastRun
func (b *BaseCrawler) LastRun() time.Time {
    return b.lastRun
}

// UpdateLastRun updates the last run time
func (b *BaseCrawler) UpdateLastRun() {
    b.lastRun = time.Now().UTC()
} 
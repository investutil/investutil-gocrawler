package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
    "github.com/yourusername/investutil-gocrawler/internal/database"
    "github.com/yourusername/investutil-gocrawler/internal/queue"
)

// Config represents the application configuration
type Config struct {
    Database struct {
        MongoDB database.Config `yaml:"mongodb"`
    } `yaml:"database"`
    Queue struct {
        RabbitMQ queue.Config `yaml:"rabbitmq"`
    } `yaml:"queue"`
    Collector struct {
        Schedule string `yaml:"schedule"`
    } `yaml:"collector"`
}

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
    absPath, err := filepath.Abs(path)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    return &cfg, nil
} 
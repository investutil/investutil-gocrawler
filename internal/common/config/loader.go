package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from a YAML file into the provided struct
func LoadConfig(configPath string, cfg interface{}) error {
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return fmt.Errorf("failed to get absolute path: %w", err)
    }

    data, err := os.ReadFile(absPath)
    if err != nil {
        return fmt.Errorf("failed to read config file: %w", err)
    }

    if err := yaml.Unmarshal(data, cfg); err != nil {
        return fmt.Errorf("failed to parse config: %w", err)
    }

    return nil
}

// LoadCommonConfig loads the common configuration and merges it with specific config
func LoadCommonConfig(commonPath, specificPath string, cfg interface{}) error {
    // Load common config first
    if err := LoadConfig(commonPath, cfg); err != nil {
        return fmt.Errorf("failed to load common config: %w", err)
    }

    // Then load and merge specific config
    if err := LoadConfig(specificPath, cfg); err != nil {
        return fmt.Errorf("failed to load specific config: %w", err)
    }

    return nil
} 
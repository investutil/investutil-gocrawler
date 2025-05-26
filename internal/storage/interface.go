package storage

import "context"

// Storage defines the interface for data storage operations
type Storage interface {
    // Save saves data to storage
    Save(ctx context.Context, key string, data interface{}) error
    
    // Load loads data from storage
    Load(ctx context.Context, key string, v interface{}) error
    
    // Delete deletes data from storage
    Delete(ctx context.Context, key string) error
    
    // List lists all keys with the given prefix
    List(ctx context.Context, prefix string) ([]string, error)
} 
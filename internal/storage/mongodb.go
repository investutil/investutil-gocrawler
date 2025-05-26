package storage

import (
    "context"
    "fmt"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStorage struct {
    client     *mongo.Client
    database   string
    collection string
}

type MongoDBConfig struct {
    URI        string `yaml:"uri"`
    Database   string `yaml:"database"`
    Collection string `yaml:"collection"`
}

func NewMongoDBStorage(cfg MongoDBConfig) (*MongoDBStorage, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    // Ping the database
    if err := client.Ping(ctx, nil); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    return &MongoDBStorage{
        client:     client,
        database:   cfg.Database,
        collection: cfg.Collection,
    }, nil
}

func (m *MongoDBStorage) Save(ctx context.Context, key string, data interface{}) error {
    coll := m.client.Database(m.database).Collection(m.collection)
    
    doc := bson.M{
        "_id":         key,
        "data":        data,
        "updated_at":  time.Now(),
    }

    opts := options.Update().SetUpsert(true)
    _, err := coll.UpdateOne(ctx, 
        bson.M{"_id": key}, 
        bson.M{"$set": doc}, 
        opts,
    )
    
    if err != nil {
        return fmt.Errorf("failed to save document: %w", err)
    }
    
    return nil
}

func (m *MongoDBStorage) Load(ctx context.Context, key string, v interface{}) error {
    coll := m.client.Database(m.database).Collection(m.collection)
    
    result := coll.FindOne(ctx, bson.M{"_id": key})
    if err := result.Err(); err != nil {
        return fmt.Errorf("failed to find document: %w", err)
    }

    var doc bson.M
    if err := result.Decode(&doc); err != nil {
        return fmt.Errorf("failed to decode document: %w", err)
    }

    data, ok := doc["data"]
    if !ok {
        return fmt.Errorf("data field not found in document")
    }

    bsonData, err := bson.Marshal(data)
    if err != nil {
        return fmt.Errorf("failed to marshal data: %w", err)
    }

    if err := bson.Unmarshal(bsonData, v); err != nil {
        return fmt.Errorf("failed to unmarshal data: %w", err)
    }

    return nil
}

func (m *MongoDBStorage) Delete(ctx context.Context, key string) error {
    coll := m.client.Database(m.database).Collection(m.collection)
    
    _, err := coll.DeleteOne(ctx, bson.M{"_id": key})
    if err != nil {
        return fmt.Errorf("failed to delete document: %w", err)
    }
    
    return nil
}

func (m *MongoDBStorage) List(ctx context.Context, prefix string) ([]string, error) {
    coll := m.client.Database(m.database).Collection(m.collection)
    
    filter := bson.M{"_id": bson.M{"$regex": fmt.Sprintf("^%s", prefix)}}
    
    cursor, err := coll.Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to list documents: %w", err)
    }
    defer cursor.Close(ctx)

    var keys []string
    for cursor.Next(ctx) {
        var doc bson.M
        if err := cursor.Decode(&doc); err != nil {
            return nil, fmt.Errorf("failed to decode document: %w", err)
        }
        if id, ok := doc["_id"].(string); ok {
            keys = append(keys, id)
        }
    }

    return keys, nil
}

func (m *MongoDBStorage) Close(ctx context.Context) error {
    return m.client.Disconnect(ctx)
} 
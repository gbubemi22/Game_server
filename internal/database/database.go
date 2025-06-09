package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var globalClient *mongo.Client


func init() {
    
	client, err := DBinstance()
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB client: %v", err)
	}
	globalClient = client
}
func DBinstance() (*mongo.Client, error) {
    MongoDb := os.Getenv("DB_DATABASE")
    if MongoDb == "" {
        return nil, fmt.Errorf("DB_DATABASE environment variable not set")
    }

    fmt.Println("GOT HEAR")

    // Use a longer timeout (e.g., 30 seconds)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    clientOptions := options.Client().
        ApplyURI(MongoDb).
        SetMaxPoolSize(100).               // Max connections
        SetMinPoolSize(10).                // Min connections
        SetMaxConnIdleTime(5 * time.Minute) // How long a connection can be idle

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    // Ping with same context or create a new one if needed
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }

    fmt.Println("Connected to MongoDB")
    return client, nil
}


func New() (*mongo.Database, error) {
	client, err := DBinstance()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize service: %w", err)
	}
	return client.Database("Gaming"), nil
}

type Service struct {
	db *mongo.Client
}

// OpenCollection opens a MongoDB collection
func (s *Service) OpenCollection(collectionName string) *mongo.Collection {
	collection := s.db.Database("Gaming").Collection(collectionName)
	return collection
}
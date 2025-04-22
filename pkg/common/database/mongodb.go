package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDBConfig struct {
	URI          string
	DatabaseName string
	Timeout      time.Duration
}

func NewMongoDBConnection(cfg MongoDBConfig) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(cfg.URI)
	
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	
	return client.Database(cfg.DatabaseName), nil
}

func Collection(db *mongo.Database, name string) *mongo.Collection {
	return db.Collection(name)
}
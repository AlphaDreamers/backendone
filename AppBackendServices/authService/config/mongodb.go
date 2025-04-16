package config

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDBConfig holds MongoDB connection configuration
type MongoDBConfig struct {
	URI            string
	DatabaseName   string
	ConnectTimeout time.Duration
	MaxPoolSize    uint64
	MinPoolSize    uint64
}

// NewMongoDBConfig creates a new MongoDB configuration with default values
func NewMongoDBConfig() *MongoDBConfig {
	return &MongoDBConfig{
		URI:            "mongodb://localhost:27017",
		DatabaseName:   "auth_service",
		ConnectTimeout: 10 * time.Second,
		MaxPoolSize:    100,
		MinPoolSize:    5,
	}
}

// Connect establishes a connection to MongoDB
func (c *MongoDBConfig) Connect(logger *logrus.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.ConnectTimeout)
	defer cancel()

	// Set client options
	clientOptions := options.Client().
		ApplyURI(c.URI).
		SetMaxPoolSize(c.MaxPoolSize).
		SetMinPoolSize(c.MinPoolSize).
		SetConnectTimeout(c.ConnectTimeout).
		SetServerSelectionTimeout(c.ConnectTimeout)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to MongoDB")
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.WithError(err).Error("Failed to ping MongoDB")
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("Successfully connected to MongoDB")
	return client.Database(c.DatabaseName), nil
}

// Close closes the MongoDB connection
func (c *MongoDBConfig) Close(client *mongo.Client, logger *logrus.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		logger.WithError(err).Error("Failed to disconnect from MongoDB")
	}
}

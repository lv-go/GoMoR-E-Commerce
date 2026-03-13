package config

import (
	"context"
	"log"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context) *mongo.Database {
	slog.Info("Connecting to MongoDB...")
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://root:password@localhost:27017"
	}
	slog.Info("MongoDB URI: ", "uri", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
		return nil
	}
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = "gomor-e-commerce"
	}
	slog.Info("MongoDB DB Name: ", "dbName", dbName)

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
		return nil
	}

	// close the connection when context is cancelled
	go func() {
		<-ctx.Done()
		slog.Info("Closing MongoDB connection...")
		client.Disconnect(ctx)
		slog.Info("MongoDB connection closed")
	}()

	slog.Info("Connected to MongoDB")
	return client.Database(dbName)
}

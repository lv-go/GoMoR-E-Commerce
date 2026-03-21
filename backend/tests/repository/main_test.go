package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"gomor-e-commerce/internal/app"

	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())

	DB = app.ConnectDB(ctx)
	if DB == nil {
		log.Fatal("Failed to connect to database in TestMain")
	}

	code := m.Run()

	cancel()                          // This will trigger the disconnect goroutine in ConnectDB
	time.Sleep(10 * time.Millisecond) // Allow logs to print
	os.Exit(code)
}

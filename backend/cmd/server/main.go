package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gomor-e-commerce/internal/config"
	"gomor-e-commerce/internal/handlers"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := config.ConnectDB(ctx)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	// Create repositories
	userRepo := repository.NewMongoCRUDRepository[models.User, primitive.ObjectID](db, "users")
	productRepo := repository.NewMongoCRUDRepository[models.Product, primitive.ObjectID](db, "products")
	categoryRepo := repository.NewMongoCRUDRepository[models.Category, primitive.ObjectID](db, "categories")
	orderRepo := repository.NewMongoCRUDRepository[models.Order, primitive.ObjectID](db, "orders")

	// Create router
	mux := http.NewServeMux()

	// API base path
	apiPath := "/api/v1"

	// Create handlers
	handlers.NewCRUDHandler(mux, userRepo, apiPath+"/users")
	handlers.NewCRUDHandler(mux, productRepo, apiPath+"/products")
	handlers.NewCRUDHandler(mux, categoryRepo, apiPath+"/categories")
	handlers.NewCRUDHandler(mux, orderRepo, apiPath+"/orders")

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	slog.Info("Server started on :" + port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}

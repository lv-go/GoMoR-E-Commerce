package app

import (
	"context"
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/handlers"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Setup(ctx context.Context) *http.ServeMux {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	db := ConnectDB(ctx)

	// Create repositories
	userRepo := repository.NewMongoCRUDRepository[models.User, primitive.ObjectID](db, "users")
	productRepo := repository.NewMongoCRUDRepository[models.Product, primitive.ObjectID](db, "products")
	categoryRepo := repository.NewMongoCRUDRepository[models.Category, primitive.ObjectID](db, "categories")
	orderRepo := repository.NewMongoCRUDRepository[models.Order, primitive.ObjectID](db, "orders")

	// Create router
	mux := http.NewServeMux()

	// Add Auth middleware
	authClient := auth.Setup(ctx)
	authMiddleware := auth.NewAuthMiddleware(authClient)

	// API base handler
	apiPath := "/api/v1"
	apiMux := http.NewServeMux()
	mux.Handle(apiPath+"/", http.StripPrefix(apiPath, apiMux))

	// Create handlers
	authMiddleware.AuthTokenMiddleware(apiMux)
	handlers.SetupCRUDHandler(apiMux, userRepo, "/users")
	handlers.SetupCRUDHandler(apiMux, productRepo, "/products")
	handlers.SetupCRUDHandler(apiMux, categoryRepo, "/categories")
	handlers.SetupCRUDHandler(apiMux, orderRepo, "/orders")

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return mux
}

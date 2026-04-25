package app

import (
	"context"
	"errors"
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/handlers"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repositories"
	"gomor-e-commerce/internal/repository"
	"gomor-e-commerce/internal/utils"
	"log/slog"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Setup(ctx context.Context) http.Handler {
	// Connect to database
	db := ConnectDB(ctx)

	// Create repositories
	userRepo := repository.NewMongoCRUDRepository[models.User, string](db, "users")
	productRepo := repository.NewMongoCRUDRepository[models.Product, primitive.ObjectID](db, "products")
	categoryRepo := repository.NewMongoCRUDRepository[models.Category, primitive.ObjectID](db, "categories")
	orderRepo := repositories.NewOrdersRepository(db)

	// Create router
	mux := http.NewServeMux()

	// Add Auth middleware
	authClient := auth.NewClient(ctx)
	authMiddleware := auth.NewAuthMiddleware(authClient, userRepo)

	// API base handler
	apiPath := "/api"
	apiMux := http.NewServeMux()
	mux.Handle(apiPath+"/", http.StripPrefix(apiPath, apiMux))

	// Create handlers
	handlers.SetupUsersHandlers(apiMux, authMiddleware, userRepo)
	handlers.SetupProductsHandlers(apiMux, authMiddleware, productRepo)
	handlers.SetupCategoriesHandlers(apiMux, authMiddleware, categoryRepo)
	handlers.SetupOrdersHandlers(apiMux, authMiddleware, orderRepo, productRepo)

	// PayPal config
	apiMux.HandleFunc("/config/paypal", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("PayPal config")
		clientId := os.Getenv("PAYPAL_CLIENT_ID")
		if clientId == "" {
			slog.Error("PAYPAL_CLIENT_ID not set")
			utils.InternalServerError(w, r, errors.New("PAYPAL_CLIENT_ID not set"))
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"clientId": clientId,
		})
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	return corsMiddleware(mux)
}

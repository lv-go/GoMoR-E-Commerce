package app

import (
	"context"
	"errors"
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/handlers"
	"gomor-e-commerce/internal/models"
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
	orderRepo := repository.NewMongoCRUDRepository[models.Order, primitive.ObjectID](db, "orders")

	// Create router
	mux := http.NewServeMux()

	// Add Auth middleware
	authClient := auth.Setup(ctx)
	authMiddleware := auth.NewAuthMiddleware(authClient)

	// API base handler
	apiPath := "/api"
	apiMux := http.NewServeMux()
	mux.Handle(apiPath+"/", http.StripPrefix(apiPath, apiMux))

	// Create handlers
	userHandler := handlers.NewCRUDHandler(userRepo)
	apiMux.HandleFunc(http.MethodPost+" /users", authMiddleware.IsAdmin(userHandler.Create))
	apiMux.HandleFunc(http.MethodPut+" /users/{id}", authMiddleware.IsAdmin(userHandler.Update))
	apiMux.HandleFunc(http.MethodDelete+" /users/{id}", authMiddleware.IsAdmin(userHandler.DeleteById))
	apiMux.HandleFunc(http.MethodGet+" /users/{id}", authMiddleware.IsAdmin(userHandler.FindById))
	apiMux.HandleFunc(http.MethodGet+" /users", authMiddleware.IsAdmin(userHandler.FindPage))

	productHandler := handlers.NewProductHandler(productRepo, db)
	apiMux.HandleFunc(http.MethodPost+" /products", authMiddleware.IsAdmin(productHandler.Create))
	apiMux.HandleFunc(http.MethodPut+" /products/{id}", authMiddleware.IsAdmin(productHandler.Update))
	apiMux.HandleFunc(http.MethodDelete+" /products/{id}", authMiddleware.IsAdmin(productHandler.DeleteById))
	apiMux.HandleFunc(http.MethodGet+" /products/{id}", productHandler.FindById)
	apiMux.HandleFunc(http.MethodGet+" /products", productHandler.FindPage)
	apiMux.HandleFunc(http.MethodGet+" /products/brands", productHandler.GetBrands)

	categoryHandler := handlers.NewCRUDHandler(categoryRepo)
	apiMux.HandleFunc(http.MethodPost+" /categories", authMiddleware.IsAdmin(categoryHandler.Create))
	apiMux.HandleFunc(http.MethodPut+" /categories/{id}", authMiddleware.IsAdmin(categoryHandler.Update))
	apiMux.HandleFunc(http.MethodDelete+" /categories/{id}", authMiddleware.IsAdmin(categoryHandler.DeleteById))
	apiMux.HandleFunc(http.MethodGet+" /categories/{id}", categoryHandler.FindById)
	apiMux.HandleFunc(http.MethodGet+" /categories", categoryHandler.FindPage)

	orderHandler := handlers.NewOrderHandler(orderRepo, productRepo)
	apiMux.HandleFunc(http.MethodPost+" /orders", authMiddleware.IsAuthenticated(orderHandler.Create))
	apiMux.HandleFunc(http.MethodPut+" /orders/{id}", authMiddleware.IsAdmin(orderHandler.Update))
	apiMux.HandleFunc(http.MethodDelete+" /orders/{id}", authMiddleware.IsAdmin(orderHandler.DeleteById))
	apiMux.HandleFunc(http.MethodGet+" /orders/{id}", authMiddleware.IsAuthenticated(orderHandler.FindById))
	apiMux.HandleFunc(http.MethodGet+" /orders", authMiddleware.IsAdmin(orderHandler.FindPage))

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

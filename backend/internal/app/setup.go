package app

import (
	"context"
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/handlers"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Setup(ctx context.Context) http.Handler {
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
	apiPath := "/api"
	apiMux := http.NewServeMux()
	mux.Handle(apiPath+"/", http.StripPrefix(apiPath, apiMux))

	// Create handlers
	userHandler := handlers.NewCRUDHandler(userRepo)
	apiMux.HandleFunc(http.MethodPost+" /users", authMiddleware.IsAdmin(userHandler.Create))
	apiMux.HandleFunc(http.MethodPut+" /users/{id}", authMiddleware.IsAdmin(userHandler.Update))
	apiMux.HandleFunc(http.MethodDelete+" /users/{id}", authMiddleware.IsAdmin(userHandler.Delete))
	apiMux.HandleFunc(http.MethodGet+" /users/{id}", authMiddleware.IsAdmin(userHandler.FindById))
	apiMux.HandleFunc(http.MethodGet+" /users", authMiddleware.IsAdmin(userHandler.FindPage))

	productHandler := handlers.NewCRUDHandler(productRepo)
	apiMux.HandleFunc(http.MethodPost+" /products", authMiddleware.IsAdmin(productHandler.Create))
	apiMux.HandleFunc(http.MethodPut+" /products/{id}", authMiddleware.IsAdmin(productHandler.Update))
	apiMux.HandleFunc(http.MethodDelete+" /products/{id}", authMiddleware.IsAdmin(productHandler.Delete))
	apiMux.HandleFunc(http.MethodGet+" /products/{id}", productHandler.FindById)
	apiMux.HandleFunc(http.MethodGet+" /products", productHandler.FindPage)

	categoryHandler := handlers.NewCRUDHandler(categoryRepo)
	apiMux.HandleFunc(http.MethodPost+" /categories", authMiddleware.IsAdmin(categoryHandler.Create))
	apiMux.HandleFunc(http.MethodPut+" /categories/{id}", authMiddleware.IsAdmin(categoryHandler.Update))
	apiMux.HandleFunc(http.MethodDelete+" /categories/{id}", authMiddleware.IsAdmin(categoryHandler.Delete))
	apiMux.HandleFunc(http.MethodGet+" /categories/{id}", categoryHandler.FindById)
	apiMux.HandleFunc(http.MethodGet+" /categories", categoryHandler.FindPage)

	orderHandler := handlers.NewCRUDHandler(orderRepo)
	apiMux.HandleFunc(http.MethodPost+" /orders", authMiddleware.IsAuthenticated(orderHandler.Create))
	apiMux.HandleFunc(http.MethodPut+" /orders/{id}", authMiddleware.IsAdmin(orderHandler.Update))
	apiMux.HandleFunc(http.MethodDelete+" /orders/{id}", authMiddleware.IsAdmin(orderHandler.Delete))
	apiMux.HandleFunc(http.MethodGet+" /orders/{id}", authMiddleware.IsAuthenticated(orderHandler.FindById))
	apiMux.HandleFunc(http.MethodGet+" /orders", authMiddleware.IsAdmin(orderHandler.FindPage))

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
			next.ServeHTTP(w, r)
		})
	}

	return corsMiddleware(mux)
}

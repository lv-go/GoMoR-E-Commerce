package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/handlers"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	args := m.Called(ctx, idToken)
	if args.Get(0) != nil {
		return args.Get(0).(*auth.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestOrderHandler(t *testing.T) {
	authClient := new(MockAuthClient)
	authMiddleware := auth.NewAuthMiddleware(authClient)
	adminToken := &auth.Token{
		UID: "user123",
		Claims: map[string]interface{}{
			"email":          "testadmin@example.com",
			"username":       "testadmin",
			"role":           "admin",
			"email_verified": true,
		},
	}
	userToken := &auth.Token{
		UID: "user123",
		Claims: map[string]interface{}{
			"email":          "testuser@example.com",
			"username":       "testuser",
			"role":           "user",
			"email_verified": true,
		},
	}

	authClient.On("VerifyIDToken", mock.Anything, "admin-token").Return(adminToken, nil).Maybe()
	authClient.On("VerifyIDToken", mock.Anything, "user-token").Return(userToken, nil).Maybe()

	// Create repository
	repo := repository.NewMongoCRUDRepository[models.Order, primitive.ObjectID](DB, "orders")
	productRepo := repository.NewMongoCRUDRepository[models.Product, primitive.ObjectID](DB, "products")

	// Create handler
	mux := http.NewServeMux()
	handler := handlers.NewOrderHandler(repo, productRepo)
	mux.HandleFunc("POST /orders", authMiddleware.IsAuthenticated(handler.Create))
	mux.HandleFunc("GET /orders/{id}", authMiddleware.IsAuthenticated(handler.FindById))
	mux.HandleFunc("PUT /orders/{id}", authMiddleware.IsAdmin(handler.Update))
	mux.HandleFunc("DELETE /orders/{id}", authMiddleware.IsAdmin(handler.DeleteById))
	mux.HandleFunc("GET /orders", authMiddleware.IsAuthenticated(handler.FindPage))

	product, err := productRepo.FindOne(t.Context(), bson.D{
		primitive.E{Key: "name", Value: "Test Product"},
	})
	if err != nil {
		product = &models.Product{
			Name:         "Test Product",
			Price:        100.0,
			CountInStock: 10,
			Category:     primitive.NewObjectID(),
			Description:  "Test Description",
		}
		productRepo.Create(context.Background(), product)
	}

	order := &models.Order{
		PaymentMethod: "PayPal",
		ItemsPrice:    100.0,
		TaxPrice:      10.0,
		ShippingPrice: 15.0,
		TotalPrice:    125.0,
		OrderItems: []models.OrderItem{
			{
				ProductID: *product.ID,
				Quantity:  1,
			},
		},
	}

	t.Run("Create", func(t *testing.T) {
		rr := httptest.NewRecorder()
		jsonBody, err := json.Marshal(order)
		assert.NoError(t, err)
		req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer user-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var createdOrder models.Order
		err = json.NewDecoder(rr.Body).Decode(&createdOrder)
		assert.NoError(t, err)
		assert.Equal(t, order.PaymentMethod, createdOrder.PaymentMethod)
		assert.Equal(t, order.TotalPrice, createdOrder.TotalPrice)
		assert.NotNil(t, createdOrder.ID)
		assert.NotEqual(t, primitive.NilObjectID, createdOrder.ID)
		assert.NotNil(t, createdOrder.Auditable.CreatedAt)
		assert.NotNil(t, createdOrder.Auditable.UpdatedAt)

		order.ID = createdOrder.ID
		order.Auditable.CreatedAt = createdOrder.Auditable.CreatedAt
		order.Auditable.UpdatedAt = createdOrder.Auditable.UpdatedAt
	})

	t.Run("FindById", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/orders/"+order.ID.Hex(), nil)
		req.Header.Set("Authorization", "Bearer user-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var foundOrder models.Order
		err := json.NewDecoder(rr.Body).Decode(&foundOrder)
		assert.NoError(t, err)
		assert.Equal(t, order.PaymentMethod, foundOrder.PaymentMethod)
		assert.Equal(t, order.ID, foundOrder.ID)
	})

	t.Run("Update", func(t *testing.T) {
		rr := httptest.NewRecorder()
		order.PaymentMethod = "Stripe"
		jsonBody, err := json.Marshal(order)
		assert.NoError(t, err)
		req := httptest.NewRequest("PUT", "/orders/"+order.ID.Hex(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Authorization", "Bearer admin-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var updatedOrder models.Order
		err = json.NewDecoder(rr.Body).Decode(&updatedOrder)
		assert.NoError(t, err)
		assert.Equal(t, order.PaymentMethod, updatedOrder.PaymentMethod)
		assert.Equal(t, order.ID, updatedOrder.ID)

		assert.GreaterOrEqual(t, order.Auditable.UpdatedAt.Unix(), updatedOrder.Auditable.UpdatedAt.Unix())
		order.Auditable.UpdatedAt = updatedOrder.Auditable.UpdatedAt
	})

	t.Run("FindPage", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/orders?limit=100&offset=0", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var page repository.Page[models.Order]
		err := json.NewDecoder(rr.Body).Decode(&page)
		assert.NoError(t, err)
		assert.LessOrEqual(t, int64(1), page.Total)
		assert.LessOrEqual(t, int32(1), page.Size)
		assert.LessOrEqual(t, int32(1), page.TotalPages)
		assert.LessOrEqual(t, 1, len(page.Items))
		found := false
		for _, item := range page.Items {
			if item.ID != nil && item.ID.Hex() == order.ID.Hex() {
				found = true
				assert.Equal(t, order.PaymentMethod, item.PaymentMethod)
				assert.Equal(t, order.TotalPrice, item.TotalPrice)
				assert.Equal(t, order.TaxPrice, item.TaxPrice)
				assert.Equal(t, order.ShippingPrice, item.ShippingPrice)
				assert.Equal(t, order.ItemsPrice, item.ItemsPrice)
				assert.Equal(t, order.IsPaid, item.IsPaid)
				assert.Equal(t, order.IsDelivered, item.IsDelivered)
				assert.Equal(t, order.Auditable.CreatedAt.Unix(), item.Auditable.CreatedAt.Unix())
				assert.Equal(t, order.Auditable.UpdatedAt.Unix(), item.Auditable.UpdatedAt.Unix())
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Delete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/orders/"+order.ID.Hex(), nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/orders/"+order.ID.Hex(), nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

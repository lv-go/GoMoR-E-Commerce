package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gomor-e-commerce/internal/handlers"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOrderHandler(t *testing.T) {
	// Create repository
	repo := repository.NewMongoCRUDRepository[models.Order, primitive.ObjectID](DB, "orders")

	// Create handler
	mux := http.NewServeMux()
	handler := handlers.NewCRUDHandler(repo)
	mux.HandleFunc("POST /orders", handler.Create)
	mux.HandleFunc("GET /orders/{id}", handler.FindById)
	mux.HandleFunc("PUT /orders/{id}", handler.Update)
	mux.HandleFunc("DELETE /orders/{id}", handler.DeleteById)
	mux.HandleFunc("GET /orders", handler.FindPage)

	order := &models.Order{
		PaymentMethod: "PayPal",
		ItemsPrice:    100.0,
		TaxPrice:      10.0,
		ShippingPrice: 5.0,
		TotalPrice:    115.0,
	}

	t.Run("Create", func(t *testing.T) {
		rr := httptest.NewRecorder()
		jsonBody, err := json.Marshal(order)
		assert.NoError(t, err)
		mux.ServeHTTP(rr, httptest.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusCreated, rr.Code)

		var createdOrder models.Order
		err = json.NewDecoder(rr.Body).Decode(&createdOrder)
		assert.NoError(t, err)
		assert.Equal(t, order.PaymentMethod, createdOrder.PaymentMethod)
		assert.Equal(t, order.TotalPrice, createdOrder.TotalPrice)
		assert.NotNil(t, createdOrder.ID)
		assert.NotEqual(t, primitive.NilObjectID, createdOrder.ID)
		assert.NotNil(t, createdOrder.CreatedAt)
		assert.NotNil(t, createdOrder.UpdatedAt)

		order.ID = createdOrder.ID
		order.CreatedAt = createdOrder.CreatedAt
		order.UpdatedAt = createdOrder.UpdatedAt
	})

	t.Run("FindById", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/orders/"+order.ID.Hex(), nil))
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
		mux.ServeHTTP(rr, httptest.NewRequest("PUT", "/orders/"+order.ID.Hex(), bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusOK, rr.Code)

		var updatedOrder models.Order
		err = json.NewDecoder(rr.Body).Decode(&updatedOrder)
		assert.NoError(t, err)
		assert.Equal(t, order.PaymentMethod, updatedOrder.PaymentMethod)
		assert.Equal(t, order.ID, updatedOrder.ID)
	})

	t.Run("FindPage", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/orders?limit=10&offset=0", nil))
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
				assert.Equal(t, order.CreatedAt.Unix(), item.CreatedAt.Unix())
				assert.Equal(t, order.UpdatedAt.Unix(), item.UpdatedAt.Unix())
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Delete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("DELETE", "/orders/"+order.ID.Hex(), nil))
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/orders/"+order.ID.Hex(), nil))
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

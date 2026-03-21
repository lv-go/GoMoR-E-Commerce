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
	handlers.SetupCRUDHandler(mux, repo, "/orders")

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
		order.ID = createdOrder.ID
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

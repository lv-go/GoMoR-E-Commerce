package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gomor-e-commerce/internal/handlers"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repositories"
	"gomor-e-commerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOrderHandler(t *testing.T) {
	// Create repository
	repo := repositories.NewOrdersRepository(DB)
	productRepo := repository.NewMongoCRUDRepository[models.Product, primitive.ObjectID](DB, "products")

	// Create handler
	mux := http.NewServeMux()
	handlers.SetupOrdersHandlers(mux, authMiddleware, repo, productRepo)

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
		productRepo.Create(t.Context(), product)
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
		assert.GreaterOrEqual(t, page.Total, int64(1))
		assert.GreaterOrEqual(t, page.Size, int32(1))
		assert.GreaterOrEqual(t, page.TotalPages, int32(1))
		assert.GreaterOrEqual(t, len(page.Items), 1)
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

	t.Run("GetTotal", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/orders/total", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var total int64
		err = json.NewDecoder(rr.Body).Decode(&total)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
	})

	t.Run("GetTotalSales", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/orders/total-sales", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var totalSales float64
		err = json.NewDecoder(rr.Body).Decode(&totalSales)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, totalSales, float64(1))
	})

	t.Run("GetTotalSalesByDate", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/orders/total-sales-by-date", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var totalSalesByDate []models.OrderSalesTotal
		err = json.NewDecoder(rr.Body).Decode(&totalSalesByDate)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(totalSalesByDate), 1)
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

package repository

import (
	"testing"
	"time"

	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repositories"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOrderRepository(t *testing.T) {
	var err error

	// Create repository
	repo := repositories.NewOrdersRepository(DB)

	// Test data
	userID := primitive.NewObjectID().Hex()
	productID := primitive.NewObjectID()
	order := &models.Order{
		UserID: userID,
		OrderItems: []models.OrderItem{
			{
				Name:      "Test Item",
				Quantity:  1,
				Image:     "https://picsum.photos/id/237/200/300",
				Price:     99.99,
				ProductID: productID,
			},
		},
		ShippingAddress: models.ShippingAddress{
			Address:    "123 Test St",
			City:       "Test City",
			PostalCode: "12345",
			Country:    "Test Country",
		},
		PaymentMethod: "PayPal",
		ItemsPrice:    99.99,
		TaxPrice:      10.00,
		ShippingPrice: 5.00,
		TotalPrice:    114.99,
		IsPaid:        false,
		IsDelivered:   false,
	}

	// Test Create
	t.Run("Create", func(t *testing.T) {
		err = repo.Create(t.Context(), order)
		assert.NoError(t, err)
		assert.NotNil(t, order.ID)
	})

	// Test FindById
	var foundOrder *models.Order
	t.Run("FindById", func(t *testing.T) {
		foundOrder, err = repo.FindById(t.Context(), *order.ID)
		assert.NoError(t, err)
		assert.Equal(t, order.TotalPrice, foundOrder.TotalPrice)
		assert.Equal(t, order.PaymentMethod, foundOrder.PaymentMethod)
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		foundOrder.IsPaid = true
		foundOrder.PaidAt = time.Now()
		err = repo.Update(t.Context(), foundOrder)
		assert.NoError(t, err)
	})

	// Test FindById after update
	t.Run("FindByIdAfterUpdate", func(t *testing.T) {
		foundOrder, err = repo.FindById(t.Context(), *order.ID)
		assert.NoError(t, err)
		assert.True(t, foundOrder.IsPaid)
	})

	// Test GetTotal Sales
	t.Run("GetTotalSales", func(t *testing.T) {
		totalSales, err := repo.GetTotalSales(t.Context())
		assert.NoError(t, err)
		assert.NotNil(t, totalSales)
		assert.GreaterOrEqual(t, totalSales, 0.00)
	})

	// Test GetTotal Sales By Date
	t.Run("GetTotalSalesByDate", func(t *testing.T) {
		totalSales, err := repo.GetTotalSalesByDate(t.Context())
		assert.NoError(t, err)
		assert.NotNil(t, totalSales)
		for _, sale := range totalSales {
			assert.NotNil(t, sale.ID)
			assert.GreaterOrEqual(t, sale.Total, 0.00)
		}
	})

	// Test GetTotal
	t.Run("GetTotal", func(t *testing.T) {
		totalSales, err := repo.GetTotal(t.Context())
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, totalSales, int64(1))
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err = repo.Delete(t.Context(), *order.ID)
		assert.NoError(t, err)
	})

	// Test FindById after delete
	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		foundOrder, err = repo.FindById(t.Context(), *order.ID)
		assert.Error(t, err)
		assert.Nil(t, foundOrder)
	})
}

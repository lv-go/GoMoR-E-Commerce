package repository

import (
	"testing"
	"time"

	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestOrderRepository(t *testing.T) {
	var err error

	// Create repository
	repo := repository.NewMongoCRUDRepository[models.Order, primitive.ObjectID](DB, "orders")

	// Test data
	userID := primitive.NewObjectID()
	productID := primitive.NewObjectID()
	order := &models.Order{
		User: userID,
		OrderItems: []models.OrderItem{
			{
				Name:    "Test Item",
				Qty:     1,
				Image:   "/images/test.jpg",
				Price:   99.99,
				Product: productID,
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
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
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
		foundOrder, err = repo.FindById(t.Context(), order.ID)
		assert.NoError(t, err)
		assert.Equal(t, order.TotalPrice, foundOrder.TotalPrice)
		assert.Equal(t, order.PaymentMethod, foundOrder.PaymentMethod)
	})

	// Test Update
	foundOrder.IsPaid = true
	foundOrder.PaidAt = time.Now()
	t.Run("Update", func(t *testing.T) {
		err = repo.Update(t.Context(), foundOrder)
		assert.NoError(t, err)
	})

	// Test FindById after update
	t.Run("FindByIdAfterUpdate", func(t *testing.T) {
		foundOrder, err = repo.FindById(t.Context(), order.ID)
		assert.NoError(t, err)
		assert.True(t, foundOrder.IsPaid)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err = repo.Delete(t.Context(), order.ID)
		assert.NoError(t, err)
	})

	// Test FindById after delete
	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		foundOrder, err = repo.FindById(t.Context(), order.ID)
		assert.Error(t, err)
		assert.Nil(t, foundOrder)
	})
}

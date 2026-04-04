package repository

import (
	"testing"

	"gomor-e-commerce/internal/models"

	"gomor-e-commerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProductRepository(t *testing.T) {
	var err error

	// Create repository
	repo := repository.NewMongoCRUDRepository[models.Product, primitive.ObjectID](DB, "products")

	// Test data
	categoryID := primitive.NewObjectID()
	product := &models.Product{
		Name:         "Test Product",
		Image:        "/images/test.jpg",
		Brand:        "Test Brand",
		Quantity:     10,
		Category:     categoryID,
		Description:  "Test Description",
		Price:        99.99,
		CountInStock: 5,
	}

	// Test Create
	t.Run("Create", func(t *testing.T) {
		err = repo.Create(t.Context(), product)
		assert.NoError(t, err)
		assert.NotNil(t, product.ID)
	})

	// Test FindById
	var foundProduct *models.Product
	t.Run("FindById", func(t *testing.T) {
		foundProduct, err = repo.FindById(t.Context(), *product.ID)
		assert.NoError(t, err)
		assert.Equal(t, product.Name, foundProduct.Name)
		assert.Equal(t, product.Price, foundProduct.Price)
	})

	// Test Update
	foundProduct.Name = "Updated Product"
	t.Run("Update", func(t *testing.T) {
		err = repo.Update(t.Context(), foundProduct)
		assert.NoError(t, err)
	})

	// Test FindById after update
	t.Run("FindByIdAfterUpdate", func(t *testing.T) {
		foundProduct, err = repo.FindById(t.Context(), *product.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Product", foundProduct.Name)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err = repo.Delete(t.Context(), *product.ID)
		assert.NoError(t, err)
	})

	// Test FindById after delete
	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		foundProduct, err = repo.FindById(t.Context(), *product.ID)
		assert.Error(t, err)
		assert.Nil(t, foundProduct)
	})
}

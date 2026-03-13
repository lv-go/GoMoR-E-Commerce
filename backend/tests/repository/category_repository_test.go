package repository

import (
	"testing"

	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCategoryRepository(t *testing.T) {
	// Create repository
	repo := repository.NewMongoCRUDRepository[models.Category, primitive.ObjectID](DB, "categories")

	// Test Create
	category := &models.Category{
		Name: "Electronics",
	}
	var err error
	t.Run("Create", func(t *testing.T) {
		err = repo.Create(t.Context(), category)
		assert.NoError(t, err)
		assert.NotNil(t, category.ID)
	})

	// Test FindById
	var foundCategory *models.Category
	t.Run("FindById", func(t *testing.T) {
		foundCategory, err = repo.FindById(t.Context(), category.ID)
		assert.NoError(t, err)
		assert.Equal(t, category.Name, foundCategory.Name)
	})

	// Test Update
	foundCategory.Name = "Electronics Updated"
	t.Run("Update", func(t *testing.T) {
		err = repo.Update(t.Context(), foundCategory)
		assert.NoError(t, err)
	})

	// Test FindById after update
	t.Run("FindByIdAfterUpdate", func(t *testing.T) {
		foundCategory, err = repo.FindById(t.Context(), category.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Electronics Updated", foundCategory.Name)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err = repo.Delete(t.Context(), category.ID)
		assert.NoError(t, err)
	})

	// Test FindById after delete
	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		foundCategory, err = repo.FindById(t.Context(), category.ID)
		assert.Error(t, err)
		assert.Nil(t, foundCategory)
	})
}

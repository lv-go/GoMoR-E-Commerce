package repository

import (
	"log/slog"
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
		slog.Info("Category created", "category", category)
		assert.NoError(t, err)
		assert.NotNil(t, category.ID)
		assert.NotEqual(t, primitive.NilObjectID, category.ID)
	})

	// Test FindById
	var foundCategory *models.Category
	t.Run("FindById", func(t *testing.T) {
		foundCategory, err = repo.FindById(t.Context(), *category.ID)
		slog.Info("Category found", "category", foundCategory)
		assert.NoError(t, err)
		assert.Equal(t, category.Name, foundCategory.Name)
		assert.Equal(t, category.ID, foundCategory.ID)
		assert.NotNil(t, foundCategory.CreatedAt)
		assert.NotNil(t, foundCategory.UpdatedAt)
		assert.Equal(t, foundCategory.CreatedAt.Unix(), category.CreatedAt.Unix())
		assert.Equal(t, foundCategory.UpdatedAt.Unix(), category.UpdatedAt.Unix())
	})

	// Test Update
	foundCategory.Name = "Electronics Updated"
	t.Run("Update", func(t *testing.T) {
		err = repo.Update(t.Context(), foundCategory)
		assert.NoError(t, err)
		assert.Equal(t, "Electronics Updated", foundCategory.Name)
		assert.NotNil(t, foundCategory.UpdatedAt)
		assert.Equal(t, foundCategory.UpdatedAt.Unix(), category.UpdatedAt.Unix())
	})

	// Test FindById after update
	t.Run("FindByIdAfterUpdate", func(t *testing.T) {
		foundCategory, err = repo.FindById(t.Context(), *category.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Electronics Updated", foundCategory.Name)
		assert.NotNil(t, foundCategory.UpdatedAt)
		assert.Equal(t, foundCategory.UpdatedAt.Unix(), category.UpdatedAt.Unix())
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err = repo.Delete(t.Context(), *category.ID)
		assert.NoError(t, err)
	})

	// Test FindById after delete
	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		foundCategory, err = repo.FindById(t.Context(), *category.ID)
		assert.Error(t, err)
		assert.Nil(t, foundCategory)
	})
}

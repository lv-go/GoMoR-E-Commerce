package repository

import (
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserHooks(t *testing.T) {
	model := &models.User{
		Name:     "test",
		Email:    "test",
		IsActive: true,
		Role:     "test",
	}

	repository.BeforeCreate[string](model)
	assert.NotNil(t, model.CreatedAt)
	assert.NotNil(t, model.UpdatedAt)

	repository.BeforeUpdate[string](model)
	assert.NotNil(t, model.UpdatedAt)
}

func TestCategoryHooks(t *testing.T) {
	model := &models.Category{
		Name: "test",
	}

	repository.BeforeCreate[string](model)
	assert.NotNil(t, model.CreatedAt)
	assert.NotNil(t, model.UpdatedAt)

	repository.BeforeUpdate[string](model)
	assert.NotNil(t, model.UpdatedAt)
}

func TestProductHooks(t *testing.T) {
	model := &models.Product{
		Name:        "test",
		Description: "test",
		Price:       10,
	}

	repository.BeforeCreate[string](model)
	assert.NotNil(t, model.CreatedAt)
	assert.NotNil(t, model.UpdatedAt)

	repository.BeforeUpdate[string](model)
	assert.NotNil(t, model.UpdatedAt)
}

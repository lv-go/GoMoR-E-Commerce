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

func TestCategoryHandler(t *testing.T) {
	// Create repository
	repo := repository.NewMongoCRUDRepository[models.Category, primitive.ObjectID](DB, "categories")

	// Create handler
	mux := http.NewServeMux()
	handlers.SetupCRUDHandler(mux, repo, "/categories")

	category := &models.Category{
		Name: "Electronics",
	}

	t.Run("Create", func(t *testing.T) {
		rr := httptest.NewRecorder()
		jsonBody, err := json.Marshal(category)
		assert.NoError(t, err)
		mux.ServeHTTP(rr, httptest.NewRequest("POST", "/categories", bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusCreated, rr.Code)

		var createdCategory models.Category
		err = json.NewDecoder(rr.Body).Decode(&createdCategory)
		assert.NoError(t, err)
		assert.Equal(t, category.Name, createdCategory.Name)
		assert.NotNil(t, createdCategory.ID)
		category.ID = createdCategory.ID
	})

	t.Run("FindById", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/categories/"+category.ID.Hex(), nil))
		assert.Equal(t, http.StatusOK, rr.Code)

		var foundCategory models.Category
		err := json.NewDecoder(rr.Body).Decode(&foundCategory)
		assert.NoError(t, err)
		assert.Equal(t, category.Name, foundCategory.Name)
		assert.Equal(t, category.ID, foundCategory.ID)
	})

	t.Run("Update", func(t *testing.T) {
		rr := httptest.NewRecorder()
		category.Name = "Electronics Updated"
		jsonBody, err := json.Marshal(category)
		assert.NoError(t, err)
		mux.ServeHTTP(rr, httptest.NewRequest("PUT", "/categories/"+category.ID.Hex(), bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusOK, rr.Code)

		var updatedCategory models.Category
		err = json.NewDecoder(rr.Body).Decode(&updatedCategory)
		assert.NoError(t, err)
		assert.Equal(t, category.Name, updatedCategory.Name)
		assert.Equal(t, category.ID, updatedCategory.ID)
	})

	t.Run("Delete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("DELETE", "/categories/"+category.ID.Hex(), nil))
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/categories/"+category.ID.Hex(), nil))
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

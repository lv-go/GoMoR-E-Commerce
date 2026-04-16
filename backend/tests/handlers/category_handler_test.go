package handlers

import (
	"bytes"
	"encoding/json"
	"log/slog"
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
	handler := handlers.NewCRUDHandler(repo)
	mux.HandleFunc("POST /categories", handler.Create)
	mux.HandleFunc("GET /categories/{id}", handler.FindById)
	mux.HandleFunc("PUT /categories/{id}", handler.Update)
	mux.HandleFunc("DELETE /categories/{id}", handler.DeleteById)
	mux.HandleFunc("GET /categories", handler.FindPage)

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
		assert.NotNil(t, createdCategory.Auditable.UpdatedAt)
		assert.NotNil(t, createdCategory.Auditable.CreatedAt)
		category.ID = createdCategory.ID
		category.Auditable.UpdatedAt = createdCategory.Auditable.UpdatedAt
		category.Auditable.CreatedAt = createdCategory.Auditable.CreatedAt
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

	t.Run("FindPage", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/categories?limit=10&offset=0", nil))
		assert.Equal(t, http.StatusOK, rr.Code)

		var page repository.Page[models.Category]
		err := json.NewDecoder(rr.Body).Decode(&page)
		assert.NoError(t, err)
		assert.LessOrEqual(t, int64(1), page.Total)
		assert.LessOrEqual(t, int32(1), page.Size)
		assert.LessOrEqual(t, int32(1), page.TotalPages)
		assert.LessOrEqual(t, 1, len(page.Items))
		found := false
		for _, item := range page.Items {
			slog.Info("FindPage:", "item", item)
			if item.ID != nil && category.ID != nil && item.ID.Hex() == category.ID.Hex() {
				found = true
				assert.Equal(t, category.Name, item.Name)
				assert.Equal(t, category.ID, item.ID)
				assert.Equal(t, category.Auditable.UpdatedAt.Unix(), item.Auditable.UpdatedAt.Unix())
				assert.Equal(t, category.Auditable.CreatedAt.Unix(), item.Auditable.CreatedAt.Unix())
				break
			}
		}
		assert.True(t, found)
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

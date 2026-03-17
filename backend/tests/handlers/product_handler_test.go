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

func TestProductHandler(t *testing.T) {
	// Create repository
	repo := repository.NewMongoCRUDRepository[models.Product, primitive.ObjectID](DB, "products")

	// Create handler
	mux := http.NewServeMux()
	handlers.NewCRUDHandler(mux, repo, "/products")

	product := &models.Product{
		Name:         "Test Product",
		Brand:        "Test Brand",
		Description:  "Test Description",
		Price:        99.99,
		CountInStock: 10,
	}

	t.Run("Create", func(t *testing.T) {
		rr := httptest.NewRecorder()
		jsonBody, err := json.Marshal(product)
		assert.NoError(t, err)
		mux.ServeHTTP(rr, httptest.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusCreated, rr.Code)

		var createdProduct models.Product
		err = json.NewDecoder(rr.Body).Decode(&createdProduct)
		assert.NoError(t, err)
		assert.Equal(t, product.Name, createdProduct.Name)
		assert.Equal(t, product.Price, createdProduct.Price)
		assert.NotNil(t, createdProduct.ID)
		product.ID = createdProduct.ID
	})

	t.Run("FindById", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/products/"+product.ID.Hex(), nil))
		assert.Equal(t, http.StatusOK, rr.Code)

		var foundProduct models.Product
		err := json.NewDecoder(rr.Body).Decode(&foundProduct)
		assert.NoError(t, err)
		assert.Equal(t, product.Name, foundProduct.Name)
		assert.Equal(t, product.ID, foundProduct.ID)
	})

	t.Run("Update", func(t *testing.T) {
		rr := httptest.NewRecorder()
		product.Name = "Test Product Updated"
		jsonBody, err := json.Marshal(product)
		assert.NoError(t, err)
		mux.ServeHTTP(rr, httptest.NewRequest("PUT", "/products/"+product.ID.Hex(), bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusOK, rr.Code)

		var updatedProduct models.Product
		err = json.NewDecoder(rr.Body).Decode(&updatedProduct)
		assert.NoError(t, err)
		assert.Equal(t, product.Name, updatedProduct.Name)
		assert.Equal(t, product.ID, updatedProduct.ID)
	})

	t.Run("Delete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("DELETE", "/products/"+product.ID.Hex(), nil))
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/products/"+product.ID.Hex(), nil))
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

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
	handler := handlers.NewCRUDHandler(repo)
	mux.HandleFunc("POST /products", handler.Create)
	mux.HandleFunc("GET /products/{id}", handler.FindById)
	mux.HandleFunc("PUT /products/{id}", handler.Update)
	mux.HandleFunc("DELETE /products/{id}", handler.DeleteById)
	mux.HandleFunc("GET /products", handler.FindPage)

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
		assert.NotEqual(t, primitive.NilObjectID, createdProduct.ID)
		assert.NotNil(t, createdProduct.CreatedAt)
		assert.NotNil(t, createdProduct.UpdatedAt)

		product.ID = createdProduct.ID
		product.CreatedAt = createdProduct.CreatedAt
		product.UpdatedAt = createdProduct.UpdatedAt
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

	t.Run("FindPage", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/products?limit=10&offset=0", nil))
		assert.Equal(t, http.StatusOK, rr.Code)

		var page repository.Page[models.Product]
		err := json.NewDecoder(rr.Body).Decode(&page)
		assert.NoError(t, err)
		assert.LessOrEqual(t, int64(1), page.Total)
		assert.LessOrEqual(t, int32(1), page.Size)
		assert.LessOrEqual(t, int32(1), page.TotalPages)
		assert.LessOrEqual(t, 1, len(page.Items))
		found := false
		for _, item := range page.Items {
			if item.ID != nil && item.ID.Hex() == product.ID.Hex() {
				found = true
				assert.Equal(t, product.Name, item.Name)
				assert.Equal(t, product.Brand, item.Brand)
				assert.Equal(t, product.Description, item.Description)
				assert.Equal(t, product.Price, item.Price)
				assert.Equal(t, product.CountInStock, item.CountInStock)
				assert.Equal(t, product.CreatedAt.Unix(), item.CreatedAt.Unix())
				assert.Equal(t, product.UpdatedAt.Unix(), item.UpdatedAt.Unix())
				break
			}
		}
		assert.True(t, found)
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

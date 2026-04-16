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
	"syreclabs.com/go/faker"
)

func TestUserHandler(t *testing.T) {
	// Create repository
	repo := repository.NewMongoCRUDRepository[models.User, string](DB, "users")

	// Create handler
	mux := http.NewServeMux()
	handler := handlers.NewCRUDHandler(repo)
	mux.HandleFunc("POST /users", handler.Create)
	mux.HandleFunc("GET /users/{id}", handler.FindById)
	mux.HandleFunc("PUT /users/{id}", handler.Update)
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteById)
	mux.HandleFunc("GET /users", handler.FindPage)

	userId := "user-" + faker.RandomString(10)
	user := &models.User{
		ID:       userId,
		Name:     faker.Name().FirstName() + faker.Name().LastName(),
		Email:    faker.Internet().Email(),
		IsActive: false,
		Role:     "user",
	}

	t.Run("Create", func(t *testing.T) {
		rr := httptest.NewRecorder()
		jsonBody, err := json.Marshal(user)
		assert.NoError(t, err)
		mux.ServeHTTP(rr, httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusCreated, rr.Code)

		var createdUser models.User
		err = json.NewDecoder(rr.Body).Decode(&createdUser)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, createdUser.Name)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.Equal(t, user.IsActive, createdUser.IsActive)
		assert.Equal(t, user.Role, createdUser.Role)
		assert.Equal(t, userId, createdUser.ID)
		assert.NotNil(t, createdUser.Auditable.CreatedAt)
		assert.NotNil(t, createdUser.Auditable.UpdatedAt)

		user.ID = createdUser.ID
		user.Auditable.CreatedAt = createdUser.Auditable.CreatedAt
		user.Auditable.UpdatedAt = createdUser.Auditable.UpdatedAt
	})

	t.Run("FindById", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/users/"+userId, nil))
		assert.Equal(t, http.StatusOK, rr.Code)

		var foundUser models.User
		err := json.NewDecoder(rr.Body).Decode(&foundUser)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.ID, foundUser.ID)
	})

	t.Run("Update", func(t *testing.T) {
		rr := httptest.NewRecorder()
		user.Name = "testuser_updated"
		jsonBody, err := json.Marshal(user)
		assert.NoError(t, err)
		mux.ServeHTTP(rr, httptest.NewRequest("PUT", "/users/"+userId, bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusOK, rr.Code)

		var updatedUser models.User
		err = json.NewDecoder(rr.Body).Decode(&updatedUser)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, updatedUser.Name)
		assert.Equal(t, user.ID, updatedUser.ID)
	})

	t.Run("FindPage", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/users?limit=10&offset=0", nil))
		assert.Equal(t, http.StatusOK, rr.Code)

		var page repository.Page[models.User]
		err := json.NewDecoder(rr.Body).Decode(&page)
		assert.NoError(t, err)
		assert.LessOrEqual(t, int64(1), page.Total)
		assert.LessOrEqual(t, int32(1), page.Size)
		assert.LessOrEqual(t, int32(1), page.TotalPages)
		assert.LessOrEqual(t, 1, len(page.Items))
		found := false
		for _, item := range page.Items {
			if item.ID == user.ID {
				found = true
				assert.Equal(t, user.Name, item.Name)
				assert.Equal(t, user.Email, item.Email)
				assert.Equal(t, user.IsActive, item.IsActive)
				assert.Equal(t, user.Role, item.Role)
				assert.Equal(t, user.Auditable.CreatedAt.Unix(), item.Auditable.CreatedAt.Unix())
				assert.Equal(t, user.Auditable.UpdatedAt.Unix(), item.Auditable.UpdatedAt.Unix())
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("Delete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("DELETE", "/users/"+userId, nil))
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/users/"+userId, nil))
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

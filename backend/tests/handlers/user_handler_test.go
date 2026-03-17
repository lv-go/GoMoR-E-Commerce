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

func TestUserHandler(t *testing.T) {
	// Create repository
	repo := repository.NewMongoCRUDRepository[models.User, primitive.ObjectID](DB, "users")

	// Create handler
	mux := http.NewServeMux()
	handlers.NewCRUDHandler(mux, repo, "/users")

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
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
		assert.Equal(t, user.Username, createdUser.Username)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.NotNil(t, createdUser.ID)
		user.ID = createdUser.ID
	})

	t.Run("FindById", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/users/"+user.ID.Hex(), nil))
		assert.Equal(t, http.StatusOK, rr.Code)

		var foundUser models.User
		err := json.NewDecoder(rr.Body).Decode(&foundUser)
		assert.NoError(t, err)
		assert.Equal(t, user.Username, foundUser.Username)
		assert.Equal(t, user.ID, foundUser.ID)
	})

	t.Run("Update", func(t *testing.T) {
		rr := httptest.NewRecorder()
		user.Username = "testuser_updated"
		jsonBody, err := json.Marshal(user)
		assert.NoError(t, err)
		mux.ServeHTTP(rr, httptest.NewRequest("PUT", "/users/"+user.ID.Hex(), bytes.NewBuffer(jsonBody)))
		assert.Equal(t, http.StatusOK, rr.Code)

		var updatedUser models.User
		err = json.NewDecoder(rr.Body).Decode(&updatedUser)
		assert.NoError(t, err)
		assert.Equal(t, user.Username, updatedUser.Username)
		assert.Equal(t, user.ID, updatedUser.ID)
	})

	t.Run("Delete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("DELETE", "/users/"+user.ID.Hex(), nil))
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest("GET", "/users/"+user.ID.Hex(), nil))
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

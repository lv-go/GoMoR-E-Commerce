package repository

import (
	"testing"

	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestUserRepository(t *testing.T) {
	var err error

	// Create repository
	repo := repository.NewMongoCRUDRepository[models.User, string](DB, "users")

	// Test data
	id := "user-" + faker.RandomString(10)
	user := &models.User{
		ID:       id,
		Name:     faker.Name().FirstName() + faker.Name().LastName(),
		Email:    faker.Internet().Email(),
		IsActive: false,
		Role:     "user",
	}

	// Test Create
	t.Run("Create", func(t *testing.T) {
		err = repo.Create(t.Context(), user)
		assert.NoError(t, err)
		assert.NotNil(t, user.ID)
	})

	// Test FindById
	var foundUser *models.User
	t.Run("FindById", func(t *testing.T) {
		foundUser, err = repo.FindById(t.Context(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
	})

	// Test Update
	foundUser.Name = "updateduser"
	t.Run("Update", func(t *testing.T) {
		err = repo.Update(t.Context(), foundUser)
		assert.NoError(t, err)
	})

	// Test FindById after update
	t.Run("FindByIdAfterUpdate", func(t *testing.T) {
		foundUser, err = repo.FindById(t.Context(), user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "updateduser", foundUser.Name)
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err = repo.Delete(t.Context(), user.ID)
		assert.NoError(t, err)
	})

	// Test FindById after delete
	t.Run("FindByIdAfterDelete", func(t *testing.T) {
		foundUser, err = repo.FindById(t.Context(), user.ID)
		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})
}

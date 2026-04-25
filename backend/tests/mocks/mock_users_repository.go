package mocks

import (
	"context"
	"gomor-e-commerce/internal/models"

	"gomor-e-commerce/internal/repository"

	"github.com/stretchr/testify/mock"
)

type MockUsersRepository struct {
	repository.CRUDRepository[models.User, string]
	mock.Mock
}

func (m *MockUsersRepository) Save(
	ctx context.Context,
	data *models.User,
	oneOpts ...repository.OneOpts,
) error {
	args := m.Called(ctx, data, oneOpts)
	return args.Error(0)
}

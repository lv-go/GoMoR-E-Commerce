package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsAuthenticatedMiddleware(t *testing.T) {
	mockClient := new(mocks.MockAuthClient)
	mockUserRepo := new(mocks.MockUsersRepository)
	mockUserRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	middleware := auth.NewAuthMiddleware(mockClient, mockUserRepo)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetUserFromContext(r)
		if user != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(user.ID))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	handler := middleware.IsAuthenticated(nextHandler)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "unauthorized")
	})

	t.Run("Malformed Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "BearerTokenWithoutSpace")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "unauthorized")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()

		mockClient.On("VerifyIDToken", mock.Anything, "invalid-token").Return((*auth.Token)(nil), errors.New("invalid token")).Once()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "unauthorized")
	})

	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()

		validToken := &auth.Token{
			UID: "user123",
			Claims: map[string]interface{}{
				"email":          "test@example.com",
				"username":       "testuser",
				"role":           "user",
				"email_verified": true,
			},
		}

		mockClient.On("VerifyIDToken", mock.Anything, "valid-token").Return(validToken, nil).Once()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "user123", rec.Body.String())
	})
}

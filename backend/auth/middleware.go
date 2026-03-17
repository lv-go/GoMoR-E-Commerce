package auth

import (
	"fmt"
	"net/http"
	"strings"

	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/utils"
)

type Config struct {
	Auth struct {
		Basic struct {
			User string
			Pass string
		}
	}
	RateLimiter struct {
		Enabled bool
	}
}

type Middlewares struct {
	config     Config
	authClient Client
}

func NewMiddlewares(
	authClient Client,
	config Config,
) Middlewares {
	return Middlewares{
		authClient: authClient,
		config:     config,
	}
}

func (m *Middlewares) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		idToken := parts[1]
		jwtToken, err := m.authClient.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			utils.UnauthorizedErrorResponse(w, r, err)
			return
		}

		// The UID from the token is not a standard UUID. We generate a deterministic UUIDv5 from it.
		r = SetUserInContext(r, &models.User{
			ID:       jwtToken.UID,
			Email:    jwtToken.Claims["email"].(string),
			Username: jwtToken.Claims["username"].(string),
			Role:     jwtToken.Claims["role"].(string),
			IsActive: jwtToken.Claims["email_verified"].(bool),
		})
		next.ServeHTTP(w, r)
	})
}

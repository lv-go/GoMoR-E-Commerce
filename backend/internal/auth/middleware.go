package auth

import (
	"fmt"
	"net/http"
	"strings"

	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/utils"
)

type AuthMiddleware struct {
	authClient Client
}

func NewAuthMiddleware(
	authClient Client,
) AuthMiddleware {
	return AuthMiddleware{
		authClient: authClient,
	}
}

func (m *AuthMiddleware) IsAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Helper function to safely get string claims
		getStrClaim := func(claims map[string]interface{}, key string) string {
			if val, ok := claims[key]; ok {
				if s, ok := val.(string); ok {
					return s
				}
			}
			return ""
		}

		// Helper function to safely get bool claims
		getBoolClaim := func(claims map[string]interface{}, key string) bool {
			if val, ok := claims[key]; ok {
				if b, ok := val.(bool); ok {
					return b
				}
			}
			return false
		}

		r = SetUserInContext(r, &models.User{
			ID:       jwtToken.UID,
			Email:    getStrClaim(jwtToken.Claims, "email"),
			Name:     getStrClaim(jwtToken.Claims, "name"),
			Role:     getStrClaim(jwtToken.Claims, "role"),
			IsActive: getBoolClaim(jwtToken.Claims, "email_verified"),
		})
		next(w, r)
	}
}

func (m *AuthMiddleware) IsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return m.IsAuthenticated(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		if user.Role != "admin" {
			utils.UnauthorizedErrorResponse(w, r, fmt.Errorf("user is not admin"))
			return
		}
		next(w, r)
	})
}

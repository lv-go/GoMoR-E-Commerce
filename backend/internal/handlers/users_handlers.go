package handlers

import (
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"net/http"
)

func SetupUsersHandlers(apiMux *http.ServeMux, authMiddleware *auth.AuthMiddleware, repo repository.CRUDRepository[models.User, string]) {
	h := NewCRUDHandler(repo)

	usersUri := "/users"

	apiMux.HandleFunc(http.MethodPost+" "+usersUri, authMiddleware.IsAdmin(h.Create))
	apiMux.HandleFunc(http.MethodPut+" "+usersUri+"/{id}", authMiddleware.IsAdmin(h.Update))
	apiMux.HandleFunc(http.MethodDelete+" "+usersUri+"/{id}", authMiddleware.IsAdmin(h.DeleteById))
	apiMux.HandleFunc(http.MethodGet+" "+usersUri+"/{id}", authMiddleware.IsAdmin(h.FindById))
	apiMux.HandleFunc(http.MethodGet+" "+usersUri, authMiddleware.IsAdmin(h.FindPage))
}

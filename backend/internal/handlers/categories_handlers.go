package handlers

import (
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupCategoriesHandlers(
	apiMux *http.ServeMux,
	authMiddleware *auth.AuthMiddleware,
	repo repository.CRUDRepository[models.Category, primitive.ObjectID],
) {
	handler := NewCRUDHandler(repo)

	categoriesUri := "/categories"

	apiMux.HandleFunc(http.MethodPost+" "+categoriesUri, authMiddleware.IsAdmin(handler.Create))
	apiMux.HandleFunc(http.MethodPut+" "+categoriesUri+"/{id}", authMiddleware.IsAdmin(handler.Update))
	apiMux.HandleFunc(http.MethodDelete+" "+categoriesUri+"/{id}", authMiddleware.IsAdmin(handler.DeleteById))
	apiMux.HandleFunc(http.MethodGet+" "+categoriesUri+"/{id}", handler.FindById)
	apiMux.HandleFunc(http.MethodGet+" "+categoriesUri, handler.FindPage)
}

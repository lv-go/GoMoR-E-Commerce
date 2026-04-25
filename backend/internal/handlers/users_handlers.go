package handlers

import (
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"gomor-e-commerce/internal/utils"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type usersHandler struct {
	*CRUDHandler[models.User, string]
}

func SetupUsersHandlers(apiMux *http.ServeMux, authMiddleware *auth.AuthMiddleware, repo repository.CRUDRepository[models.User, string]) {
	h := usersHandler{
		CRUDHandler: NewCRUDHandler(repo),
	}

	usersUri := "/users"

	apiMux.HandleFunc(http.MethodPost+" "+usersUri, authMiddleware.IsAdmin(h.Create))
	apiMux.HandleFunc(http.MethodPut+" "+usersUri+"/{id}", authMiddleware.IsAdmin(h.Update))
	apiMux.HandleFunc(http.MethodDelete+" "+usersUri+"/{id}", authMiddleware.IsAdmin(h.DeleteById))
	apiMux.HandleFunc(http.MethodGet+" "+usersUri+"/{id}", authMiddleware.IsAdmin(h.FindById))
	apiMux.HandleFunc(http.MethodGet+" "+usersUri, authMiddleware.IsAdmin(h.FindPage))
	apiMux.HandleFunc(http.MethodGet+" "+usersUri+"/total", authMiddleware.IsAdmin(h.GetTotal))
}

// GetTotal returns the number of users matching the given filter.
// Optional filter parameter: role
func (h *usersHandler) GetTotal(w http.ResponseWriter, r *http.Request) {
	slog.Info("GetTotal")

	role := r.URL.Query().Get("role")
	filter := bson.D{}
	if role != "" {
		filter = append(filter, bson.E{
			Key:   "role",
			Value: role,
		})
	}

	count, err := h.repo.Count(r.Context(), filter)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, count)
}

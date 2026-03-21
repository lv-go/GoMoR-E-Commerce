package handlers

import (
	"encoding/json"
	"errors"
	"gomor-e-commerce/internal/repository"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CRUDHandler[T any, ID comparable] struct {
	repo repository.CRUDRepository[T, ID]
}

func SetupCRUDHandler[T any, ID comparable](
	mux *http.ServeMux,
	repo repository.CRUDRepository[T, ID],
	path string,
) {
	h := &CRUDHandler[T, ID]{
		repo: repo,
	}
	mux.HandleFunc(http.MethodPost+" "+path, h.Create)
	mux.HandleFunc(http.MethodGet+" "+path+"/{id}", h.FindById)
	mux.HandleFunc(http.MethodPut+" "+path+"/{id}", h.Update)
	mux.HandleFunc(http.MethodDelete+" "+path+"/{id}", h.Delete)
	mux.HandleFunc(http.MethodGet+" "+path, h.FindPage)
}

func (h *CRUDHandler[T, ID]) parseID(r *http.Request) (*ID, error) {
	id := r.PathValue("id")
	idType := reflect.TypeFor[ID]()

	var idValue any
	var err error
	switch idType.Kind() {
	case reflect.String:
		idValue = id
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		idValue, err = strconv.Atoi(id)
	case reflect.TypeFor[primitive.ObjectID]().Kind():
		idValue, err = primitive.ObjectIDFromHex(id)
	default:
		idValue, err = nil, errors.New("invalid id type")
	}

	return new(idValue.(ID)), err
}

func (h *CRUDHandler[T, ID]) Create(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Create", "path", r.URL.Path)
	var entity T
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.repo.Create(r.Context(), &entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entity)
}

func (h *CRUDHandler[T, ID]) FindById(w http.ResponseWriter, r *http.Request) {
	slog.Debug("FindById", "id", r.PathValue("id"))
	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	entity, err := h.repo.FindById(r.Context(), *id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(entity)
}

func (h *CRUDHandler[T, ID]) Update(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Update", "id", r.PathValue("id"))
	var entity T
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.repo.Update(r.Context(), &entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity)
}

func (h *CRUDHandler[T, ID]) Delete(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Delete", "id", r.PathValue("id"))
	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.repo.Delete(r.Context(), *id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *CRUDHandler[T, ID]) FindPage(w http.ResponseWriter, r *http.Request) {
	slog.Debug("FindPage", "path", r.URL.Path)
	var err error
	var offset int
	var limit int

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	result, err := h.repo.FindPage(r.Context(), map[string]any{}, repository.ManyOpts{
		Limit:  new(int64(limit)),
		Offset: new(int64(offset)),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

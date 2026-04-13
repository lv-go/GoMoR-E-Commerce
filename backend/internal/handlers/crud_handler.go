package handlers

import (
	"encoding/json"
	"errors"
	"gomor-e-commerce/internal/repository"
	"gomor-e-commerce/internal/utils"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CRUDHandler[T any, ID comparable] struct {
	repo repository.CRUDRepository[T, ID]
}

func NewCRUDHandler[T any, ID comparable](repo repository.CRUDRepository[T, ID]) *CRUDHandler[T, ID] {
	return &CRUDHandler[T, ID]{
		repo: repo,
	}
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
	if err != nil {
		return nil, err
	}
	typedID := idValue.(ID)
	return &typedID, nil
}

func (h *CRUDHandler[T, ID]) Create(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Create", "path", r.URL.Path)
	var entity T
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	err = utils.ValidateStruct(&entity)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	err = h.repo.Create(r.Context(), &entity)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &entity)
}

func (h *CRUDHandler[T, ID]) FindById(w http.ResponseWriter, r *http.Request) {
	slog.Debug("FindById", "id", r.PathValue("id"))
	id, err := h.parseID(r)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	entity, err := h.repo.FindById(r.Context(), *id)
	if err != nil {
		utils.NotFoundResponse(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, entity)
}

func (h *CRUDHandler[T, ID]) Update(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Update", "id", r.PathValue("id"))
	var entity T
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	err = utils.ValidateStruct(&entity)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	err = h.repo.Update(r.Context(), &entity)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, entity)
}

func (h *CRUDHandler[T, ID]) DeleteById(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Delete", "id", r.PathValue("id"))
	id, err := h.parseID(r)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}
	err = h.repo.Delete(r.Context(), *id)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSONMessage(w, http.StatusOK, "Deleted successfully")
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
			utils.BadRequestResponse(w, r, err)
			return
		}
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			utils.BadRequestResponse(w, r, err)
			return
		}
	}

	result, err := h.repo.FindPage(r.Context(), map[string]any{}, repository.ManyOpts{
		Limit:  new(int64(limit)),
		Offset: new(int64(offset)),
	})
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

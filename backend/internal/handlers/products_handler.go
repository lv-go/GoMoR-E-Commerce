package handlers

import (
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"gomor-e-commerce/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductHandler struct {
	CRUDHandler[models.Product, primitive.ObjectID]
	db *mongo.Database
}

func NewProductHandler(
	repo repository.CRUDRepository[models.Product, primitive.ObjectID],
	db *mongo.Database,
) *ProductHandler {
	return &ProductHandler{
		CRUDHandler: *NewCRUDHandler(repo),
		db:          db,
	}
}

func (h *ProductHandler) GetBrands(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ProductHandler.GetBrands", "path", r.URL.Path)
	brands, err := h.db.Collection("products").Distinct(r.Context(), "brand", map[string]any{})
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, brands)
}

func (h *ProductHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ProductHandler.FindPage", "path", r.URL.Path)
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

	filter := map[string]any{}

	if brand := r.URL.Query().Get("brand"); brand != "" {
		filter["brand"] = brand
	}

	if category := r.URL.Query().Get("category"); category != "" {
		categories := strings.Split(category, ",")
		for i := range categories {
			categories[i] = strings.TrimSpace(categories[i])
		}
		objectIds := make([]primitive.ObjectID, len(categories))
		for i, category := range categories {
			objectIds[i], err = primitive.ObjectIDFromHex(category)
			if err != nil {
				utils.BadRequestResponse(w, r, err)
				return
			}
		}
		filter["category"] = map[string]any{"$in": objectIds}
	}

	// if search := r.URL.Query().Get("search"); search != "" {
	// 	filter["name"] = map[string]any{"$regex": search, "$options": "i"}
	// }

	slog.Debug("ProductHandler.FindPage", "filter", filter)

	result, err := h.repo.FindPage(r.Context(), filter, repository.ManyOpts{
		Limit:  new(int64(limit)),
		Offset: new(int64(offset)),
	})
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

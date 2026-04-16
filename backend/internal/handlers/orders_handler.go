package handlers

import (
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"gomor-e-commerce/internal/utils"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderHandler struct {
	CRUDHandler[models.Order, primitive.ObjectID]
	orderRepo   repository.CRUDRepository[models.Order, primitive.ObjectID]
	productRepo repository.CRUDRepository[models.Product, primitive.ObjectID]
}

func NewOrderHandler(orderRepo repository.CRUDRepository[models.Order, primitive.ObjectID], productRepo repository.CRUDRepository[models.Product, primitive.ObjectID]) *OrderHandler {
	return &OrderHandler{
		CRUDHandler: *NewCRUDHandler(orderRepo),
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.Create", "path", r.URL.Path)
	var order models.Order
	err := utils.ReadJSON(w, r, &order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.UserID = auth.GetUserFromContext(r).ID

	// Calculate prices
	order.ItemsPrice = 0
	for _, item := range order.OrderItems {
		product, err := h.productRepo.FindById(r.Context(), item.ProductID)
		if err != nil {
			slog.Error("Error finding product", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		order.ItemsPrice += product.Price * float64(item.Quantity)
	}

	order.ShippingPrice = 0
	if order.ItemsPrice > 100 {
		order.ShippingPrice = 0
	} else {
		order.ShippingPrice = 10
	}
	order.TaxPrice = order.ItemsPrice * 0.15
	order.TotalPrice = order.ItemsPrice + order.ShippingPrice + order.TaxPrice

	err = h.orderRepo.Create(r.Context(), &order)
	if err != nil {
		slog.Error("Error creating order", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &order)
}

func (h *OrderHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.FindPage", "path", r.URL.Path)
	user := auth.GetUserFromContext(r)

	filter := make(map[string]string)

	// if user is customer, find only their orders
	if user.Role != "admin" {
		slog.Debug("OrderHandler.FindPage", "userId", user.ID)
		filter["userId"] = user.ID
	}

	slog.Debug("OrderHandler.FindPage", "filter", filter)

	sortBy := []repository.SortBy{{
		Field:     "createdAt",
		Direction: repository.SortDirection_Descending,
	}}

	orders, err := h.orderRepo.FindPage(r.Context(), filter, repository.ManyOpts{
		SortBy: &sortBy,
	})
	if err != nil {
		slog.Error("Error finding orders", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}

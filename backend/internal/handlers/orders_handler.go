package handlers

import (
	"errors"
	"gomor-e-commerce/internal/auth"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repositories"
	"gomor-e-commerce/internal/repository"
	"gomor-e-commerce/internal/utils"
	"log/slog"
	"net/http"
	"time"

	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ordersHandlers struct {
	CRUDHandler[models.Order, primitive.ObjectID]
	orderRepo   repositories.OrdersRepository
	productRepo repository.CRUDRepository[models.Product, primitive.ObjectID]
}

func SetupOrdersHandlers(
	apiMux *http.ServeMux,
	authMiddleware *auth.AuthMiddleware,
	orderRepo repositories.OrdersRepository,
	productRepo repository.CRUDRepository[models.Product, primitive.ObjectID],
) {
	h := &ordersHandlers{
		CRUDHandler: *NewCRUDHandler(orderRepo),
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}

	ordersUri := "/orders"

	// Public routes
	apiMux.HandleFunc(http.MethodPost+" "+ordersUri, authMiddleware.IsAuthenticated(h.Create))

	// Authenticated routes
	apiMux.HandleFunc(http.MethodPut+" "+ordersUri+"/{id}", authMiddleware.IsAuthenticated(h.Update))
	apiMux.HandleFunc(http.MethodDelete+" "+ordersUri+"/{id}", authMiddleware.IsAuthenticated(h.DeleteById))
	apiMux.HandleFunc(http.MethodGet+" "+ordersUri+"/{id}", authMiddleware.IsAuthenticated(h.FindById))
	apiMux.HandleFunc(http.MethodGet+" "+ordersUri, authMiddleware.IsAuthenticated(h.FindPage))
	apiMux.HandleFunc(http.MethodPut+" "+ordersUri+"/{id}/pay", authMiddleware.IsAuthenticated(h.Pay))
	apiMux.HandleFunc(http.MethodPut+" "+ordersUri+"/{id}/deliver", authMiddleware.IsAuthenticated(h.Deliver))

	// Admin routes
	apiMux.HandleFunc(http.MethodGet+" "+ordersUri+"/total", authMiddleware.IsAdmin(h.GetTotal))
	apiMux.HandleFunc(http.MethodGet+" "+ordersUri+"/total-sales", authMiddleware.IsAdmin(h.GetTotalSales))
	apiMux.HandleFunc(http.MethodGet+" "+ordersUri+"/total-sales-by-date", authMiddleware.IsAdmin(h.GetTotalSalesByDate))
}

func (h *ordersHandlers) Create(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.Create", "path", r.URL.Path)
	var order models.Order
	err := utils.ReadJSON(w, r, &order)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	order.UserID = auth.GetUserFromContext(r).ID

	// Calculate prices
	order.ItemsPrice = 0
	for _, item := range order.OrderItems {
		product, err := h.productRepo.FindById(r.Context(), item.ProductID)
		if err != nil {
			utils.InternalServerError(w, r, err)
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
		utils.InternalServerError(w, r, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, &order)
}

func (h *ordersHandlers) FindPage(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.FindPage", "path", r.URL.Path)
	user := auth.GetUserFromContext(r)

	filter := make(map[string]string)

	// if user is customer, find only their orders
	if user.Role != "admin" {
		slog.Debug("OrderHandler.FindPage", "userId", user.ID)
		filter["userId"] = user.ID
	}

	slog.Debug("OrderHandler.FindPage", "filter", filter)

	var limit int64 = 10
	var offset int64 = 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.ParseInt(o, 10, 64); err == nil {
			offset = val
		}
	}

	sortBy := []repository.SortBy{{
		Field:     "createdAt",
		Direction: repository.SortDirection_Descending,
	}}

	orders, err := h.orderRepo.FindPage(r.Context(), filter, repository.ManyOpts{
		SortBy: &sortBy,
		Limit:  &limit,
		Offset: &offset,
	})
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, orders)
}

func (h *ordersHandlers) Pay(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.Pay", "path", r.URL.Path)
	orderIdStr := r.PathValue("id")
	orderId, err := primitive.ObjectIDFromHex(orderIdStr)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	order, err := h.orderRepo.FindById(r.Context(), orderId)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	var orderUpdate models.Order
	err = utils.ReadJSON(w, r, &orderUpdate)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
		return
	}

	order.IsPaid = true
	order.PaidAt = time.Now()
	order.PaymentResult = models.PaymentResult{
		ID:           orderUpdate.PaymentResult.ID,
		Status:       orderUpdate.PaymentResult.Status,
		UpdateTime:   orderUpdate.PaymentResult.UpdateTime,
		EmailAddress: orderUpdate.PaymentResult.EmailAddress,
	}

	err = h.orderRepo.Update(r.Context(), order)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, order)
}

func (h *ordersHandlers) Deliver(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.Deliver", "path", r.URL.Path)
	orderIdStr := r.PathValue("id")
	orderId, err := primitive.ObjectIDFromHex(orderIdStr)
	if err != nil {
		utils.BadRequestResponse(w, r, errors.New("Invalid order ID"))
		return
	}

	order, err := h.orderRepo.FindById(r.Context(), orderId)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	order.IsDelivered = true
	order.DeliveredAt = time.Now()

	err = h.orderRepo.Update(r.Context(), order)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, order)
}

// GetTotal gets the total number of orders
func (h *ordersHandlers) GetTotal(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.GetTotal", "path", r.URL.Path)

	total, err := h.orderRepo.GetTotal(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, total)
}

// GetTotalSales gets the total sales
func (h *ordersHandlers) GetTotalSales(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.GetTotalSales", "path", r.URL.Path)

	total, err := h.orderRepo.GetTotalSales(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, total)
}

// GetTotalSalesByDate gets the total sales by date
func (h *ordersHandlers) GetTotalSalesByDate(w http.ResponseWriter, r *http.Request) {
	slog.Debug("OrderHandler.GetTotalSalesByDate", "path", r.URL.Path)

	totalSalesByDate, err := h.orderRepo.GetTotalSalesByDate(r.Context())
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, totalSalesByDate)
}

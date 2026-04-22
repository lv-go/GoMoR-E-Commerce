package repositories

import (
	"context"
	"gomor-e-commerce/internal/models"
	"gomor-e-commerce/internal/repository"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrdersRepository interface {
	repository.CRUDRepository[models.Order, primitive.ObjectID]
	GetTotal(ctx context.Context) (int64, error)
	GetTotalSales(ctx context.Context) (float64, error)
	GetTotalSalesByDate(ctx context.Context) ([]models.OrderSalesTotal, error)
}

type ordersRepository struct {
	repository.CRUDRepository[models.Order, primitive.ObjectID]
	collection *mongo.Collection
}

func NewOrdersRepository(db *mongo.Database) OrdersRepository {
	return &ordersRepository{
		CRUDRepository: repository.NewMongoCRUDRepository[models.Order, primitive.ObjectID](db, "orders"),
		collection:     db.Collection("orders"),
	}
}

// GetTotal returns the total number of orders
func (r *ordersRepository) GetTotal(ctx context.Context) (int64, error) {
	slog.Debug("ordersRepository.GetTotal", "path", r.collection.Name())

	result, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return result, nil
}

// GetTotalSales calculates the total sales
func (r *ordersRepository) GetTotalSales(ctx context.Context) (float64, error) {
	slog.Debug("ordersRepository.GetTotalSales", "path", r.collection.Name())

	result, err := r.collection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"isPaid": true,
			},
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$totalPrice"},
			},
		},
	})
	if err != nil {
		return 0, err
	}

	defer result.Close(ctx)

	var orderSalesTotals []models.OrderSalesTotal
	if err := result.All(ctx, &orderSalesTotals); err != nil {
		slog.Error("ordersRepository.GetTotalSales.Decode", "error", err)
		return 0, err
	}

	return orderSalesTotals[0].Total, nil
}

// GetTotalSalesByDate gets the total sales by date
func (r *ordersRepository) GetTotalSalesByDate(ctx context.Context) ([]models.OrderSalesTotal, error) {
	slog.Debug("ordersRepository.GetTotalSalesByDate", "path", r.collection.Name())

	result, err := r.collection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"isPaid": true,
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$auditable.createdAt",
					},
				},
				"total": bson.M{"$sum": "$totalPrice"},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	defer result.Close(ctx)

	var orderSalesByDate []models.OrderSalesTotal
	if err := result.All(ctx, &orderSalesByDate); err != nil {
		return nil, err
	}

	return orderSalesByDate, nil
}

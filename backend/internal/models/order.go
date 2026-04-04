package models

import (
	"gomor-e-commerce/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	Product primitive.ObjectID `json:"product" bson:"product"`
	Name    string             `json:"name" bson:"name"`
	Qty     int                `json:"qty" bson:"qty"`
	Image   string             `json:"image" bson:"image"`
	Price   float64            `json:"price" bson:"price"`
}

type ShippingAddress struct {
	Address    string `json:"address" bson:"address"`
	City       string `json:"city" bson:"city"`
	PostalCode string `json:"postalCode" bson:"postalCode"`
	Country    string `json:"country" bson:"country"`
}

type PaymentResult struct {
	ID           string `json:"id" bson:"id"`
	Status       string `json:"status" bson:"status"`
	UpdateTime   string `json:"update_time" bson:"update_time"`
	EmailAddress string `json:"email_address" bson:"email_address"`
}

type Order struct {
	repository.Auditable
	ID              *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User            primitive.ObjectID  `json:"user" bson:"user"`
	OrderItems      []OrderItem         `json:"orderItems" bson:"orderItems"`
	ShippingAddress ShippingAddress     `json:"shippingAddress" bson:"shippingAddress"`
	PaymentMethod   string              `json:"paymentMethod" bson:"paymentMethod"`
	PaymentResult   PaymentResult       `json:"paymentResult" bson:"paymentResult"`
	ItemsPrice      float64             `json:"itemsPrice" bson:"itemsPrice"`
	TaxPrice        float64             `json:"taxPrice" bson:"taxPrice"`
	ShippingPrice   float64             `json:"shippingPrice" bson:"shippingPrice"`
	TotalPrice      float64             `json:"totalPrice" bson:"totalPrice"`
	IsPaid          bool                `json:"isPaid" bson:"isPaid"`
	PaidAt          time.Time           `json:"paidAt" bson:"paidAt"`
	IsDelivered     bool                `json:"isDelivered" bson:"isDelivered"`
	DeliveredAt     time.Time           `json:"deliveredAt" bson:"deliveredAt"`
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Rating    float64            `json:"rating" bson:"rating"`
	Comment   string             `json:"comment" bson:"comment"`
	User      primitive.ObjectID `json:"user" bson:"user"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Product struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Image        string             `json:"image" bson:"image"`
	Brand        string             `json:"brand" bson:"brand"`
	Quantity     int                `json:"quantity" bson:"quantity"`
	Category     primitive.ObjectID `json:"category" bson:"category"`
	Description  string             `json:"description" bson:"description"`
	Reviews      []Review           `json:"reviews" bson:"reviews"`
	Rating       float64            `json:"rating" bson:"rating"`
	NumReviews   int                `json:"numReviews" bson:"numReviews"`
	Price        float64            `json:"price" bson:"price"`
	CountInStock int                `json:"countInStock" bson:"countInStock"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

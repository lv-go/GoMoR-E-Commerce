package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"gomor-e-commerce/internal/repository"
)

type Category struct {
	repository.Auditable
	ID   *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string              `json:"name" bson:"name"`
}

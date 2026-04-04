package models

import (
	"gomor-e-commerce/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	repository.Auditable
	ID   *primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string              `json:"name" bson:"name"`
}

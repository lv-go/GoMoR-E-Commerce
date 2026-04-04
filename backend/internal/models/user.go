package models

import (
	repository "gomor-e-commerce/internal/repository"
)

type User struct {
	repository.Auditable
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	IsActive bool   `json:"isActive" bson:"isActive"`
	Role     string `json:"role" bson:"role"`
}

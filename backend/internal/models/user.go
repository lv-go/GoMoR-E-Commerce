package models

import (
	"time"
)

type User struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	IsActive  bool      `json:"isActive" bson:"isActive"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

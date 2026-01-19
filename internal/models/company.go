// Package models defines MongoDB models for the application
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Company represents a company in the system
type Company struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Address   string             `json:"address" bson:"address"`
	Currency  string             `json:"currency" bson:"currency"` // Default currency (USD, VND, etc.)
	Timezone  string             `json:"timezone" bson:"timezone"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewCompany creates a new company with defaults
func NewCompany(name string) *Company {
	now := time.Now()
	return &Company{
		Name:      name,
		Currency:  "USD",
		Timezone:  "UTC",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

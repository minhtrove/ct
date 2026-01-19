// Package models defines MongoDB models for the application
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryType represents whether category is for income or expense
type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

// Category represents a transaction category
type Category struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Type        CategoryType       `json:"type" bson:"type"` // income or expense
	Color       string             `json:"color" bson:"color"`
	Icon        string             `json:"icon" bson:"icon"`
	Description string             `json:"description" bson:"description"`
	CompanyID   primitive.ObjectID `json:"company_id" bson:"company_id"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewCategory creates a new category
func NewCategory(name string, categoryType CategoryType, companyID primitive.ObjectID) *Category {
	now := time.Now()
	return &Category{
		Name:      name,
		Type:      categoryType,
		Color:     "#6366f1", // Default indigo
		IsActive:  true,
		CompanyID: companyID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsValidCategoryType checks if the category type is valid
func IsValidCategoryType(t string) bool {
	switch CategoryType(t) {
	case CategoryTypeIncome, CategoryTypeExpense:
		return true
	}
	return false
}

// GetDefaultCategories returns default categories for a new company
func GetDefaultCategories(companyID primitive.ObjectID) []*Category {
	now := time.Now()
	return []*Category{
		// Income categories
		{Name: "Sales", Type: CategoryTypeIncome, Color: "#22c55e", Icon: "dollar-sign", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Services", Type: CategoryTypeIncome, Color: "#3b82f6", Icon: "briefcase", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Investment", Type: CategoryTypeIncome, Color: "#8b5cf6", Icon: "trending-up", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Other Income", Type: CategoryTypeIncome, Color: "#06b6d4", Icon: "plus-circle", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		// Expense categories
		{Name: "Payroll", Type: CategoryTypeExpense, Color: "#ef4444", Icon: "users", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Office Supplies", Type: CategoryTypeExpense, Color: "#f97316", Icon: "package", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Utilities", Type: CategoryTypeExpense, Color: "#eab308", Icon: "zap", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Rent", Type: CategoryTypeExpense, Color: "#84cc16", Icon: "home", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Marketing", Type: CategoryTypeExpense, Color: "#ec4899", Icon: "megaphone", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Travel", Type: CategoryTypeExpense, Color: "#14b8a6", Icon: "plane", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Software", Type: CategoryTypeExpense, Color: "#6366f1", Icon: "code", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Other Expense", Type: CategoryTypeExpense, Color: "#64748b", Icon: "minus-circle", CompanyID: companyID, IsActive: true, CreatedAt: now, UpdatedAt: now},
	}
}

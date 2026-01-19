// Package models defines MongoDB models for the application
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BudgetPeriod represents the budget time period
type BudgetPeriod string

const (
	BudgetPeriodMonthly   BudgetPeriod = "monthly"
	BudgetPeriodQuarterly BudgetPeriod = "quarterly"
	BudgetPeriodYearly    BudgetPeriod = "yearly"
)

// Budget represents a spending limit for a category
type Budget struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CategoryID primitive.ObjectID `json:"category_id" bson:"category_id"`
	Name       string             `json:"name" bson:"name"`
	Amount     float64            `json:"amount" bson:"amount"` // Budget limit
	Spent      float64            `json:"spent" bson:"spent"`   // Amount spent
	Currency   string             `json:"currency" bson:"currency"`
	Period     BudgetPeriod       `json:"period" bson:"period"`
	StartDate  time.Time          `json:"start_date" bson:"start_date"`
	EndDate    time.Time          `json:"end_date" bson:"end_date"`
	CompanyID  primitive.ObjectID `json:"company_id" bson:"company_id"`
	IsActive   bool               `json:"is_active" bson:"is_active"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	// Populated fields
	CategoryName string `json:"category_name,omitempty" bson:"-"`
}

// NewBudget creates a new budget
func NewBudget(name string, categoryID primitive.ObjectID, amount float64, currency string, period BudgetPeriod, companyID primitive.ObjectID) *Budget {
	now := time.Now()
	startDate, endDate := GetBudgetPeriodDates(period, now)
	return &Budget{
		CategoryID: categoryID,
		Name:       name,
		Amount:     amount,
		Spent:      0,
		Currency:   currency,
		Period:     period,
		StartDate:  startDate,
		EndDate:    endDate,
		CompanyID:  companyID,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Remaining returns the remaining budget amount
func (b *Budget) Remaining() float64 {
	remaining := b.Amount - b.Spent
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Utilization returns the budget utilization percentage
func (b *Budget) Utilization() float64 {
	if b.Amount == 0 {
		return 0
	}
	return (b.Spent / b.Amount) * 100
}

// IsOverBudget checks if spending exceeds the budget
func (b *Budget) IsOverBudget() bool {
	return b.Spent > b.Amount
}

// CanSpend checks if the amount can be spent within budget
func (b *Budget) CanSpend(amount float64) bool {
	return (b.Spent + amount) <= b.Amount
}

// AddSpending adds to the spent amount
func (b *Budget) AddSpending(amount float64) {
	b.Spent += amount
	b.UpdatedAt = time.Now()
}

// GetBudgetPeriodDates returns start and end dates for a budget period
func GetBudgetPeriodDates(period BudgetPeriod, referenceDate time.Time) (time.Time, time.Time) {
	year, month, _ := referenceDate.Date()
	loc := referenceDate.Location()

	switch period {
	case BudgetPeriodMonthly:
		startDate := time.Date(year, month, 1, 0, 0, 0, 0, loc)
		endDate := startDate.AddDate(0, 1, -1)
		return startDate, endDate
	case BudgetPeriodQuarterly:
		quarter := (int(month) - 1) / 3
		startMonth := time.Month(quarter*3 + 1)
		startDate := time.Date(year, startMonth, 1, 0, 0, 0, 0, loc)
		endDate := startDate.AddDate(0, 3, -1)
		return startDate, endDate
	case BudgetPeriodYearly:
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
		endDate := time.Date(year, 12, 31, 23, 59, 59, 0, loc)
		return startDate, endDate
	default:
		return referenceDate, referenceDate.AddDate(0, 1, 0)
	}
}

// IsValidBudgetPeriod checks if the period is valid
func IsValidBudgetPeriod(p string) bool {
	switch BudgetPeriod(p) {
	case BudgetPeriodMonthly, BudgetPeriodQuarterly, BudgetPeriodYearly:
		return true
	}
	return false
}

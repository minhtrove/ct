// Package models defines MongoDB models for the application
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccountType represents the type of financial account
type AccountType string

const (
	AccountTypeBank   AccountType = "bank"
	AccountTypeCash   AccountType = "cash"
	AccountTypeCredit AccountType = "credit"
	AccountTypeWallet AccountType = "wallet"
	AccountTypeOther  AccountType = "other"
)

// Account represents a financial account
type Account struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	Type          AccountType        `json:"type" bson:"type"`
	Balance       float64            `json:"balance" bson:"balance"`
	Currency      string             `json:"currency" bson:"currency"`
	Description   string             `json:"description" bson:"description"`
	AccountNumber string             `json:"account_number" bson:"account_number"`
	CompanyID     primitive.ObjectID `json:"company_id" bson:"company_id"`
	IsActive      bool               `json:"is_active" bson:"is_active"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

// NewAccount creates a new account with defaults
func NewAccount(name string, accountType AccountType, currency string, companyID primitive.ObjectID) *Account {
	now := time.Now()
	return &Account{
		Name:      name,
		Type:      accountType,
		Balance:   0,
		Currency:  currency,
		CompanyID: companyID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Credit adds money to the account
func (a *Account) Credit(amount float64) {
	a.Balance += amount
	a.UpdatedAt = time.Now()
}

// Debit removes money from the account
func (a *Account) Debit(amount float64) error {
	if a.Balance < amount {
		return ErrInsufficientBalance
	}
	a.Balance -= amount
	a.UpdatedAt = time.Now()
	return nil
}

// AccountTypeDisplayName returns human-readable name for account type
func AccountTypeDisplayName(t AccountType) string {
	switch t {
	case AccountTypeBank:
		return "Bank Account"
	case AccountTypeCash:
		return "Cash"
	case AccountTypeCredit:
		return "Credit Card"
	case AccountTypeWallet:
		return "Digital Wallet"
	case AccountTypeOther:
		return "Other"
	default:
		return string(t)
	}
}

// IsValidAccountType checks if the account type is valid
func IsValidAccountType(t string) bool {
	switch AccountType(t) {
	case AccountTypeBank, AccountTypeCash, AccountTypeCredit, AccountTypeWallet, AccountTypeOther:
		return true
	}
	return false
}

// GetAccountTypes returns all valid account types
func GetAccountTypes() []AccountType {
	return []AccountType{
		AccountTypeBank,
		AccountTypeCash,
		AccountTypeCredit,
		AccountTypeWallet,
		AccountTypeOther,
	}
}

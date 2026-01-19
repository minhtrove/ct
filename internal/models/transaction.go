// Package models defines MongoDB models for the application
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

// TransactionStatus represents the approval status
type TransactionStatus string

const (
	TransactionStatusPending  TransactionStatus = "pending"
	TransactionStatusApproved TransactionStatus = "approved"
	TransactionStatusRejected TransactionStatus = "rejected"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type            TransactionType    `json:"type" bson:"type"`
	Amount          float64            `json:"amount" bson:"amount"`
	Currency        string             `json:"currency" bson:"currency"`
	Description     string             `json:"description" bson:"description"`
	FromAccountID   primitive.ObjectID `json:"from_account_id,omitempty" bson:"from_account_id,omitempty"`
	ToAccountID     primitive.ObjectID `json:"to_account_id,omitempty" bson:"to_account_id,omitempty"`
	CategoryID      primitive.ObjectID `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Status          TransactionStatus  `json:"status" bson:"status"`
	CreatedByID     primitive.ObjectID `json:"created_by_id" bson:"created_by_id"`
	CreatedByName   string             `json:"created_by_name" bson:"created_by_name"`
	ApprovedByID    primitive.ObjectID `json:"approved_by_id,omitempty" bson:"approved_by_id,omitempty"`
	ApprovedByName  string             `json:"approved_by_name,omitempty" bson:"approved_by_name,omitempty"`
	RejectionReason string             `json:"rejection_reason,omitempty" bson:"rejection_reason,omitempty"`
	CompanyID       primitive.ObjectID `json:"company_id" bson:"company_id"`
	TransactionDate time.Time          `json:"transaction_date" bson:"transaction_date"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
	ApprovedAt      time.Time          `json:"approved_at,omitempty" bson:"approved_at,omitempty"`
	// Populated fields (not stored in DB)
	FromAccountName string `json:"from_account_name,omitempty" bson:"-"`
	ToAccountName   string `json:"to_account_name,omitempty" bson:"-"`
	CategoryName    string `json:"category_name,omitempty" bson:"-"`
}

// NewTransaction creates a new transaction
func NewTransaction(txType TransactionType, amount float64, currency string, companyID, createdByID primitive.ObjectID, createdByName string) *Transaction {
	now := time.Now()
	return &Transaction{
		Type:            txType,
		Amount:          amount,
		Currency:        currency,
		Status:          TransactionStatusPending,
		CreatedByID:     createdByID,
		CreatedByName:   createdByName,
		CompanyID:       companyID,
		TransactionDate: now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Approve marks the transaction as approved
func (t *Transaction) Approve(approverID primitive.ObjectID, approverName string) {
	t.Status = TransactionStatusApproved
	t.ApprovedByID = approverID
	t.ApprovedByName = approverName
	t.ApprovedAt = time.Now()
	t.UpdatedAt = time.Now()
}

// Reject marks the transaction as rejected
func (t *Transaction) Reject(approverID primitive.ObjectID, approverName, reason string) {
	t.Status = TransactionStatusRejected
	t.ApprovedByID = approverID
	t.ApprovedByName = approverName
	t.RejectionReason = reason
	t.UpdatedAt = time.Now()
}

// IsPending checks if transaction is pending approval
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// IsApproved checks if transaction is approved
func (t *Transaction) IsApproved() bool {
	return t.Status == TransactionStatusApproved
}

// TransactionTypeDisplayName returns human-readable name
func TransactionTypeDisplayName(t TransactionType) string {
	switch t {
	case TransactionTypeIncome:
		return "Income"
	case TransactionTypeExpense:
		return "Expense"
	case TransactionTypeTransfer:
		return "Transfer"
	default:
		return string(t)
	}
}

// TransactionStatusDisplayName returns human-readable status
func TransactionStatusDisplayName(s TransactionStatus) string {
	switch s {
	case TransactionStatusPending:
		return "Pending"
	case TransactionStatusApproved:
		return "Approved"
	case TransactionStatusRejected:
		return "Rejected"
	default:
		return string(s)
	}
}

// IsValidTransactionType checks if the type is valid
func IsValidTransactionType(t string) bool {
	switch TransactionType(t) {
	case TransactionTypeIncome, TransactionTypeExpense, TransactionTypeTransfer:
		return true
	}
	return false
}

// IsValidTransactionStatus checks if the status is valid
func IsValidTransactionStatus(s string) bool {
	switch TransactionStatus(s) {
	case TransactionStatusPending, TransactionStatusApproved, TransactionStatusRejected:
		return true
	}
	return false
}

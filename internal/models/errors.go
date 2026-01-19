// Package models defines MongoDB models for the application
package models

import "errors"

// Common errors
var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidTransaction  = errors.New("invalid transaction")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrBudgetExceeded      = errors.New("budget limit exceeded")
)

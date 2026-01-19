// Package models defines MongoDB models for the application
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionCreate  AuditAction = "create"
	AuditActionUpdate  AuditAction = "update"
	AuditActionDelete  AuditAction = "delete"
	AuditActionApprove AuditAction = "approve"
	AuditActionReject  AuditAction = "reject"
	AuditActionLogin   AuditAction = "login"
	AuditActionLogout  AuditAction = "logout"
)

// AuditEntity represents the entity being audited
type AuditEntity string

const (
	AuditEntityUser        AuditEntity = "user"
	AuditEntityAccount     AuditEntity = "account"
	AuditEntityTransaction AuditEntity = "transaction"
	AuditEntityCategory    AuditEntity = "category"
	AuditEntityBudget      AuditEntity = "budget"
	AuditEntityCompany     AuditEntity = "company"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID        primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	Action    AuditAction            `json:"action" bson:"action"`
	Entity    AuditEntity            `json:"entity" bson:"entity"`
	EntityID  primitive.ObjectID     `json:"entity_id,omitempty" bson:"entity_id,omitempty"`
	UserID    primitive.ObjectID     `json:"user_id" bson:"user_id"`
	UserName  string                 `json:"user_name" bson:"user_name"`
	UserEmail string                 `json:"user_email" bson:"user_email"`
	CompanyID primitive.ObjectID     `json:"company_id" bson:"company_id"`
	Changes   map[string]interface{} `json:"changes,omitempty" bson:"changes,omitempty"`
	IPAddress string                 `json:"ip_address" bson:"ip_address"`
	UserAgent string                 `json:"user_agent" bson:"user_agent"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
}

// NewAuditLog creates a new audit log entry
func NewAuditLog(action AuditAction, entity AuditEntity, entityID, userID, companyID primitive.ObjectID, userName, userEmail string) *AuditLog {
	return &AuditLog{
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		UserID:    userID,
		UserName:  userName,
		UserEmail: userEmail,
		CompanyID: companyID,
		CreatedAt: time.Now(),
	}
}

// WithChanges sets the changes field
func (a *AuditLog) WithChanges(changes map[string]interface{}) *AuditLog {
	a.Changes = changes
	return a
}

// WithIPAddress sets the IP address
func (a *AuditLog) WithIPAddress(ip string) *AuditLog {
	a.IPAddress = ip
	return a
}

// WithUserAgent sets the user agent
func (a *AuditLog) WithUserAgent(ua string) *AuditLog {
	a.UserAgent = ua
	return a
}

// AuditActionDisplayName returns human-readable action name
func AuditActionDisplayName(a AuditAction) string {
	switch a {
	case AuditActionCreate:
		return "Created"
	case AuditActionUpdate:
		return "Updated"
	case AuditActionDelete:
		return "Deleted"
	case AuditActionApprove:
		return "Approved"
	case AuditActionReject:
		return "Rejected"
	case AuditActionLogin:
		return "Logged In"
	case AuditActionLogout:
		return "Logged Out"
	default:
		return string(a)
	}
}

// AuditEntityDisplayName returns human-readable entity name
func AuditEntityDisplayName(e AuditEntity) string {
	switch e {
	case AuditEntityUser:
		return "User"
	case AuditEntityAccount:
		return "Account"
	case AuditEntityTransaction:
		return "Transaction"
	case AuditEntityCategory:
		return "Category"
	case AuditEntityBudget:
		return "Budget"
	case AuditEntityCompany:
		return "Company"
	default:
		return string(e)
	}
}

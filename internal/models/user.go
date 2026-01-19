// Package models mongo user model db
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email           string             `json:"email" bson:"email"`
	PasswordHash    string             `json:"-" bson:"password_hash"`
	Name            string             `json:"name" bson:"name"`
	Role            string             `json:"role" bson:"role"`                         // RBAC role
	CompanyID       primitive.ObjectID `json:"company_id" bson:"company_id"`             // Company association
	AvatarFileID    primitive.ObjectID `json:"avatar_file_id,omitempty" bson:"avatar_file_id,omitempty"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
	LastLoginAt     time.Time          `json:"last_login_at" bson:"last_login_at"`
	// Email verification
	EmailVerified    bool      `json:"email_verified" bson:"email_verified"`
	VerifyToken      string    `json:"-" bson:"verify_token,omitempty"`
	VerifyExpiresAt  time.Time `json:"-" bson:"verify_expires_at,omitempty"`
	LastEmailSentAt  time.Time `json:"-" bson:"last_email_sent_at,omitempty"`
	// Password reset
	ResetToken      string             `json:"-" bson:"reset_token,omitempty"`
	ResetExpiresAt  time.Time          `json:"-" bson:"reset_expires_at,omitempty"`
}


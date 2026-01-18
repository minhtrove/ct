// Package auth handles session management
package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SessionManager handles user sessions
type SessionManager struct{}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

// SetSession sets the user session cookie
func (s *SessionManager) SetSession(c *fiber.Ctx, user *models.User) error {
	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    user.ID.Hex(),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https", // Only use HTTPS in production
		SameSite: "lax",
		MaxAge:   86400 * 7, // 7 days
	})
	return nil
}

// GetSession retrieves the current user from session
func (s *SessionManager) GetSession(c *fiber.Ctx) (primitive.ObjectID, error) {
	userID := c.Cookies("user_id")
	if userID == "" {
		return primitive.NilObjectID, fiber.NewError(fiber.StatusUnauthorized, "Not authenticated")
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return primitive.NilObjectID, fiber.NewError(fiber.StatusUnauthorized, "Invalid session")
	}

	return objectID, nil
}

// ClearSession removes the user session cookie
func (s *SessionManager) ClearSession(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "user_id",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HTTPOnly: true,
		SameSite: "lax",
	})
	return nil
}

// IsAuthenticated checks if user is authenticated
func (s *SessionManager) IsAuthenticated(c *fiber.Ctx) bool {
	_, err := s.GetSession(c)
	return err == nil
}

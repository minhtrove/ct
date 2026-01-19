// Package middleware provides HTTP middleware for authentication and authorization
package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/auth"
	"github.com/minhtranin/ct/internal/handler"
	"github.com/minhtranin/ct/internal/view/shared/toast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RequireAuth checks if the user is authenticated and loads user info into context
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log for debugging
		fmt.Println("[DEBUG] RequireAuth called for path:", c.Path())

		userID, err := handler.GetSession(c)
		if err != nil || userID == "" {
			// Return toast for HTMX requests, redirect for regular requests
			if c.Get("HX-Requested-With") == "true" || c.Get("HX-Request") == "true" {
				c.Set("Content-Type", "text/html")
				return toast.Toast(toast.Props{
					Title:         "Authentication required",
					Description:  "Please sign in to continue",
					Variant:       toast.VariantError,
					Position:      toast.PositionTopLeft,
					Duration:      5000,
					Dismissible:   true,
					ShowIndicator: true,
					Icon:          true,
				}).Render(c.Context(), c.Response().BodyWriter())
			}
			return c.Redirect("/signin")
		}

		// Fetch user from database and store in context
		db := handler.GetDB()
		usersCollection := db.Database("ct").Collection("users")
		var user struct {
			ID    string `bson:"_id"`
			Role  string `bson:"role"`
			Email string `bson:"email"`
		}

		// Convert hex string to ObjectID for MongoDB query
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return c.Redirect("/signin")
		}

		err = usersCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&user)
		if err != nil {
			// User not found in database - return toast
			if c.Get("HX-Requested-With") == "true" || c.Get("HX-Request") == "true" {
				c.Set("Content-Type", "text/html")
				return toast.Toast(toast.Props{
					Title:         "Session expired",
					Description:  "Please sign in again",
					Variant:       toast.VariantError,
					Position:      toast.PositionTopLeft,
					Duration:      5000,
					Dismissible:   true,
					ShowIndicator: true,
					Icon:          true,
				}).Render(c.Context(), c.Response().BodyWriter())
			}
			return c.Redirect("/signin")
		}

		// Store user info in context for middleware to use
		c.Locals("userID", userID)
		c.Locals("userRole", user.Role)
		c.Locals("userEmail", user.Email)

		return c.Next()
	}
}

// RequireRole checks if the user has the required role level
func RequireRole(requiredLevel int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if user info is in context (set by RequireAuth)
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Authentication required")
		}

		roleStr := userRole.(string)
		userLevel := auth.GetRoleLevel(roleStr)

		if userLevel < requiredLevel {
			// Return toast for HTMX requests
			if c.Get("HX-Requested-With") == "true" || c.Get("HX-Request") == "true" {
				c.Set("Content-Type", "text/html")
				return toast.Toast(toast.Props{
					Title:         "Access denied",
					Description:  "You don't have permission to access this page",
					Variant:       toast.VariantError,
					Position:      toast.PositionTopLeft,
					Duration:      5000,
					Dismissible:   true,
					ShowIndicator: true,
					Icon:          true,
				}).Render(c.Context(), c.Response().BodyWriter())
			}
			return c.Redirect("/access-denied")
		}

		return c.Next()
	}
}

// RequirePermission checks if the user has a specific permission
func RequirePermission(permissionCheck func(role string) bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Authentication required")
		}

		roleStr := userRole.(string)
		if !permissionCheck(roleStr) {
			// Return toast for HTMX requests
			if c.Get("HX-Requested-With") == "true" || c.Get("HX-Request") == "true" {
				c.Set("Content-Type", "text/html")
				return toast.Toast(toast.Props{
					Title:         "Access denied",
					Description:  "You don't have permission to access this page",
					Variant:       toast.VariantError,
					Position:      toast.PositionTopLeft,
					Duration:      5000,
					Dismissible:   true,
					ShowIndicator: true,
					Icon:          true,
				}).Render(c.Context(), c.Response().BodyWriter())
			}
			return c.Redirect("/access-denied")
		}

		return c.Next()
	}
}

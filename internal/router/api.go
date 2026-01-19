// Package router handles the routing for the API endpoints.
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/handler"
)

func registerAPI(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/health", handler.Health)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/signup", handler.SignUp)
	auth.Post("/signin", handler.SignIn)
	auth.Post("/verify-code", handler.VerifyCode)
	auth.Post("/resend-verification", handler.ResendVerification)
	auth.Post("/forgot-password", handler.ForgotPassword)
	auth.Post("/reset-password", handler.ResetPassword)
	auth.Get("/verify-email", handler.VerifyEmail)
	auth.Get("/logout", handler.Logout)

	// Account routes
	api.Post("/accounts", handler.CreateAccount)
	api.Post("/accounts/:id", handler.UpdateAccount) // POST for HTML forms
	api.Post("/accounts/:id/delete", handler.DeleteAccount)

	// Category routes
	api.Post("/categories", handler.CreateCategory)
	api.Post("/categories/:id", handler.UpdateCategory) // POST for HTML forms
	api.Post("/categories/:id/delete", handler.DeleteCategory)

	// Budget routes
	api.Post("/budgets", handler.CreateBudget)
	api.Post("/budgets/:id", handler.UpdateBudget) // POST for HTML forms
	api.Post("/budgets/:id/delete", handler.DeleteBudget)

	// Transaction routes
	api.Post("/transactions", handler.CreateTransaction)
	api.Post("/transactions/:id", handler.UpdateTransaction) // POST for HTML forms
	api.Post("/transactions/:id/delete", handler.DeleteTransaction)
	api.Post("/transactions/:id/approve", handler.ApproveTransaction)
	api.Post("/transactions/:id/reject", handler.RejectTransaction)

	// User management routes
	api.Post("/users/:id/role", handler.UpdateUserRole)
}

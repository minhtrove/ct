// Package router handles the routing for the API endpoints.
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/auth"
	"github.com/minhtranin/ct/internal/handler"
	"github.com/minhtranin/ct/internal/middleware"
)

func registerAPI(app *fiber.App) {
	// Public routes - NO auth middleware
	app.Post("/api/auth/signup", handler.SignUp)
	app.Post("/api/auth/signin", handler.SignIn)
	app.Post("/api/auth/verify-code", handler.VerifyCode)
	app.Post("/api/auth/resend-verification", handler.ResendVerification)
	app.Post("/api/auth/forgot-password", handler.ForgotPassword)
	app.Post("/api/auth/reset-password", handler.ResetPassword)
	app.Get("/api/auth/verify-email", handler.VerifyEmail)
	app.Get("/api/auth/logout", handler.Logout)
	app.Get("/api/health", handler.Health)

	// Protected routes - apply RequireAuth individually
	// Transaction routes - employee+
	app.Post("/api/transactions", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleEmployee]), handler.CreateTransaction)
	app.Post("/api/transactions/:id", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleEmployee]), handler.UpdateTransaction)
	app.Post("/api/transactions/:id/delete", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleEmployee]), handler.DeleteTransaction)

	// Approval routes - holder+
	app.Post("/api/transactions/:id/approve", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleHolder]), handler.ApproveTransaction)
	app.Post("/api/transactions/:id/reject", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleHolder]), handler.RejectTransaction)

	// Account routes - admin+
	app.Post("/api/accounts", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.CreateAccount)
	app.Post("/api/accounts/:id", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.UpdateAccount)
	app.Post("/api/accounts/:id/delete", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.DeleteAccount)

	// Category routes - admin+
	app.Post("/api/categories", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.CreateCategory)
	app.Post("/api/categories/:id", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.UpdateCategory)
	app.Post("/api/categories/:id/delete", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.DeleteCategory)

	// Budget routes - admin+
	app.Post("/api/budgets", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.CreateBudget)
	app.Post("/api/budgets/:id", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.UpdateBudget)
	app.Post("/api/budgets/:id/delete", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), handler.DeleteBudget)

	// User management routes - manager+
	app.Post("/api/users/:id/role", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleManager]), handler.UpdateUserRole)
}

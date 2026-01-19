package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/auth"
	"github.com/minhtranin/ct/internal/middleware"
	"github.com/minhtranin/ct/internal/page"
)

func registerPages(r fiber.Router) {
	// Public pages - NO auth middleware
	r.Get("/", page.Home)
	r.Get("/signin", page.SignIn)
	r.Get("/signup", page.SignUp)
	r.Get("/verify-email", page.VerifyEmailPage)
	r.Get("/forgot-password", page.ForgotPassword)
	r.Get("/reset-password", page.ResetPassword)
	r.Get("/logout", page.Logout)

	// Protected pages - apply RequireAuth individually to avoid catching API routes
	// Dashboard (all authenticated users)
	r.Get("/dashboard", middleware.RequireAuth(), page.Dashboard)

	// Access Denied page
	r.Get("/access-denied", middleware.RequireAuth(), page.AccessDenied)

	// Storybook (developer only)
	r.Get("/story", middleware.RequireAuth(), middleware.RequirePermission(auth.IsDeveloper), page.Story)

	// Income & Expense Management pages
	r.Get("/transactions", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleEmployee]), page.TransactionsPage)

	// Approval pages (holder+)
	r.Get("/approvals", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleHolder]), page.ApprovalsPage)

	// Reporting pages (accountant+)
	r.Get("/reports", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAccountant]), page.ReportsPage)
	r.Get("/audit", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAccountant]), page.AuditPage)

	// Management pages (admin+)
	r.Get("/accounts", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), page.AccountsPage)
	r.Get("/categories", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), page.CategoriesPage)
	r.Get("/budgets", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), page.BudgetsPage)
	r.Get("/settings", middleware.RequireAuth(), middleware.RequireRole(auth.RoleLevel[auth.RoleAdmin]), page.SettingsPage)

	// Team page - employee+
	r.Get("/team", middleware.RequireAuth(), page.TeamPage)
}

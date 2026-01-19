package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/page"
)

func registerPages(r fiber.Router) {
	// Public pages
	r.Get("/", page.Home)
	r.Get("/signin", page.SignIn)
	r.Get("/signup", page.SignUp)
	r.Get("/verify-email", page.VerifyEmailPage)
	r.Get("/forgot-password", page.ForgotPassword)
	r.Get("/reset-password", page.ResetPassword)
	r.Get("/logout", page.Logout)

	// Authenticated pages
	r.Get("/dashboard", page.Dashboard)
	r.Get("/story", page.Story)

	// Income & Expense Management pages
	r.Get("/transactions", page.TransactionsPage)
	r.Get("/accounts", page.AccountsPage)
	r.Get("/categories", page.CategoriesPage)
	r.Get("/budgets", page.BudgetsPage)
	r.Get("/reports", page.ReportsPage)
	r.Get("/audit", page.AuditPage)
	r.Get("/approvals", page.ApprovalsPage)
	r.Get("/settings", page.SettingsPage)
	r.Get("/team", page.TeamPage)
}

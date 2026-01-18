// Package router handles the routing for the API endpoints.
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/handler"
)

func registerAPI (app *fiber.App) {
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
}


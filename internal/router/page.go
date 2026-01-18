package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/page"
)

func registerPages(r fiber.Router) {
	r.Get("/", page.Home)
	r.Get("/signin", page.SignIn)
	r.Get("/signup", page.SignUp)
	r.Get("/dashboard", page.Dashboard)
	r.Get("/verify-email", page.VerifyEmailPage)
	r.Get("/forgot-password", page.ForgotPassword)
	r.Get("/reset-password", page.ResetPassword)
	r.Get("/logout", page.Logout)
	r.Get("/story", page.Story)
}

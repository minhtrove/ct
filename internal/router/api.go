// Package router handles the routing for the API endpoints.
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/handler"
)

func registerAPI (app *fiber.App) {
	api := app.Group("/api")
	api.Get("/health", handler.Health)
}


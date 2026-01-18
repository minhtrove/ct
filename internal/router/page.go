package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/page"
)

func registerPages(r fiber.Router) {
	r.Get("/", page.Home)
	r.Get("/story", page.Story)
}

package router

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App) {
	registerPages(app)
	registerAPI(app)
}

// Package handler defines HTTP handlers for the web application.
package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/render"
	"github.com/minhtranin/ct/internal/view/layouts"
	view "github.com/minhtranin/ct/internal/view/pages"
)

func Home(f *fiber.Ctx) error {
	return render.HTML(
		f,
		layouts.Base("home", view.HomePage("")),
	)
}

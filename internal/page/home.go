// Package page defines HTTP handlers for the web application.
package page

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minhtranin/ct/internal/render"
	view "github.com/minhtranin/ct/internal/view/components"
	"github.com/minhtranin/ct/internal/view/layouts"
)

func Home(f *fiber.Ctx) error {
	return render.HTML(
		f,
		layouts.Base("home", view.HomePage("")),
	)
}

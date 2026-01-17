// Package render provides helper functions to render templated components in Fiber handlers.
package render

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func HTML(c *fiber.Ctx, comp templ.Component) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return comp.Render(c.Context(), c.Response().BodyWriter())
}

package notfound

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
)

func New() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusNotFound).
			JSON(controllers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
}

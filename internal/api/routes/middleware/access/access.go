package access

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
)

func AdminOnly(ctx *fiber.Ctx) error {
	payload, ok := ctx.Locals("user").(*dto.TokenPayload)
	if !ok {
		return ctx.Status(http.StatusForbidden).
			JSON(handlers.HTTPError{Error: fiber.ErrForbidden.Error()})
	}
	for _, role := range payload.Roles {
		if role == "admin" {
			return ctx.Next()
		}
	}
	return ctx.Status(http.StatusForbidden).
		JSON(handlers.HTTPError{Error: fiber.ErrForbidden.Error()})
}

package access

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
)

func AdminOnly(ctx *fiber.Ctx) error {
	payload, ok := ctx.Locals("user").(*dto.TokenPayload)
	if !ok {
		return fiber.ErrForbidden
	}
	for _, role := range payload.Roles {
		if role == "admin" {
			return ctx.Next()
		}
	}
	return fiber.ErrForbidden
}

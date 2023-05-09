package notfound

import "github.com/gofiber/fiber/v2"

func New() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return fiber.ErrNotFound
	}
}

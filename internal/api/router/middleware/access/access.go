package access

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/pkg/errors"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
)

func AdminOnly(l logger.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		payload, ok := ctx.Locals("user").(*params.TokenPayload)
		if !ok {
			l.Error(ctx, logger.Client, errors.ErrFailedToGetTokenPayload)
			return ctx.Status(http.StatusForbidden).
				JSON(controllers.HTTPError{Error: fiber.ErrForbidden.Error()})
		}
		for _, role := range payload.Roles {
			if role == "admin" {
				return ctx.Next()
			}
		}
		l.Error(ctx, logger.Client, errors.ErrNotAdmin)
		return ctx.Status(http.StatusForbidden).
			JSON(controllers.HTTPError{Error: fiber.ErrForbidden.Error()})
	}
}

func OwnerOnly(l logger.Logger) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		id, ok := ctx.Locals(context.UserID).(int)
		if !ok {
			l.Error(ctx, logger.Client, errors.ErrFailedToGetUserId)
			return ctx.Status(http.StatusForbidden).
				JSON(controllers.HTTPError{Error: fiber.ErrForbidden.Error()})
		}
		payload, ok := ctx.Locals("user").(*params.TokenPayload)
		if !ok {
			l.Error(ctx, logger.Client, errors.ErrFailedToGetTokenPayload)
			return ctx.Status(http.StatusForbidden).
				JSON(controllers.HTTPError{Error: fiber.ErrForbidden.Error()})
		}
		if payload.UserID == id {
			return ctx.Next()
		}
		for _, role := range payload.Roles {
			if role == "admin" {
				return ctx.Next()
			}
		}
		l.Error(ctx, logger.Client, errors.ErrNotOwner)
		return ctx.Status(http.StatusForbidden).
			JSON(controllers.HTTPError{Error: fiber.ErrForbidden.Error()})
	}
}

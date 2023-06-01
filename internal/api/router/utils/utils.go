package utils

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/validator"
)

func ParseIDFromPath(ctx *fiber.Ctx, key string) (int, error) {
	id, err := strconv.Atoi(ctx.Params(key))
	if err != nil {
		return -1, err
	}
	return id, nil
}

func ParseFilterAndValidate(ctx *fiber.Ctx, filter interface{}) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if err := validator.Validate(filter); err != nil {
		return err
	}
	return nil
}

func ParseBody(ctx *fiber.Ctx, body interface{}) error {
	if err := ctx.BodyParser(body); err != nil {
		return fmt.Errorf("failed to parse body: %w", err)
	}
	return nil
}

func ParseBodyAndValidate(ctx *fiber.Ctx, body interface{}) error {
	if err := ParseBody(ctx, body); err != nil {
		return err
	}
	if err := validator.Validate(body); err != nil {
		return err
	}
	return nil
}

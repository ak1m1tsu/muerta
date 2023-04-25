package common

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetIdByFiberCtx(ctx *fiber.Ctx) (int, error) {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return -1, err
	}
	return id, nil
}

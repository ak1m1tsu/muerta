package common

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
)

func GetIdByFiberCtx(ctx *fiber.Ctx) (int, error) {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetNameByFiberCtx(ctx *fiber.Ctx) (string, error) {
	param := ctx.Params("name")
	if param == "" {
		return "", fmt.Errorf("name is required")
	}
	return param, nil
}

func GetRecipeFilterByFiberCtx(ctx *fiber.Ctx, filter *dto.RecipeFilterDTO) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.Paging == nil {
		filter.Paging = &dto.Paging{
			Limit:  10,
			Offset: 0,
		}
	}
	if filter.Paging.Limit <= 0 {
		filter.Paging.Limit = 10
	}
	if filter.Paging.Offset < 0 {
		filter.Paging.Offset = 0
	}
	return nil
}

func GetUserFilterByFiberCtx(ctx *fiber.Ctx, filter *dto.UserFilterDTO) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.Paging == nil {
		filter.Paging = &dto.Paging{
			Limit:  10,
			Offset: 0,
		}
	}
	if filter.Paging.Limit <= 0 {
		filter.Paging.Limit = 10
	}
	if filter.Paging.Offset < 0 {
		filter.Paging.Offset = 0
	}
	return nil
}

func GetSettingFilterByFiberCtx(ctx *fiber.Ctx, filter *dto.SettingFilterDTO) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.Paging == nil {
		filter.Paging = &dto.Paging{
			Limit:  10,
			Offset: 0,
		}
	}
	if filter.Paging.Limit <= 0 {
		filter.Paging.Limit = 10
	}
	if filter.Paging.Offset < 0 {
		filter.Paging.Offset = 0
	}
	return nil
}

func GetStorageFilterByFiberCtx(ctx *fiber.Ctx, filter *dto.StorageFilterDTO) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.Paging == nil {
		filter.Paging = &dto.Paging{
			Limit:  10,
			Offset: 0,
		}
	}
	if filter.Paging.Limit <= 0 {
		filter.Paging.Limit = 10
	}
	if filter.Paging.Offset < 0 {
		filter.Paging.Offset = 0
	}
	return nil
}

func GetProductFilterByFiberCtx(ctx *fiber.Ctx, filter *dto.ProductFilterDTO) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.Paging == nil {
		filter.Paging = &dto.Paging{
			Limit:  10,
			Offset: 0,
		}
	}
	if filter.Paging.Limit <= 0 {
		filter.Paging.Limit = 10
	}
	if filter.Paging.Offset < 0 {
		filter.Paging.Offset = 0
	}
	return nil
}

func GetRoleFilterByFiberCtx(ctx *fiber.Ctx, filter *dto.RoleFilterDTO) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.Paging == nil {
		filter.Paging = &dto.Paging{
			Limit:  10,
			Offset: 0,
		}
	}
	if filter.Paging.Limit <= 0 {
		filter.Paging.Limit = 10
	}
	if filter.Paging.Offset < 0 {
		filter.Paging.Offset = 0
	}
	return nil
}

func GetCategoryFilterByFiberCtx(ctx *fiber.Ctx, filter *dto.CategoryFilterDTO) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.Paging == nil {
		filter.Paging = &dto.Paging{
			Limit:  10,
			Offset: 0,
		}
	}
	if filter.Paging.Limit <= 0 {
		filter.Paging.Limit = 10
	}
	if filter.Paging.Offset < 0 {
		filter.Paging.Offset = 0
	}
	return nil
}

func GetTipFilterByFiberCtx(ctx *fiber.Ctx, filter *dto.TipFilterDTO) error {
	if err := ctx.QueryParser(filter); err != nil {
		return fmt.Errorf("failed to parse query: %w", err)
	}
	if filter.Paging == nil {
		filter.Paging = &dto.Paging{
			Limit:  10,
			Offset: 0,
		}
	}
	if filter.Paging.Limit <= 0 {
		filter.Paging.Limit = 10
	}
	if filter.Paging.Offset < 0 {
		filter.Paging.Offset = 0
	}
	return nil
}

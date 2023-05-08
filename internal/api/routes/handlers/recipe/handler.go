package recipe

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/services/recipe"
)

type RecipesHandler struct {
	svc recipe.RecipeServicer
	log *log.Logger
}

func New(svc recipe.RecipeServicer, log *log.Logger) *RecipesHandler {
	return &RecipesHandler{
		svc: svc,
		log: log,
	}
}

func (h *RecipesHandler) CreateRecipe(ctx *fiber.Ctx) error {
	var payload *dto.CreateRecipeDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateRecipe(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *RecipesHandler) FindRecipeByID(ctx *fiber.Ctx) error {
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	dto, err := h.svc.FindRecipeByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"recipe": dto},
	})
}

func (h *RecipesHandler) FindRecipes(ctx *fiber.Ctx) error {
	filter := new(dto.RecipeFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	dtos, err := h.svc.FindRecipes(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"recipes": dtos},
	})
}

func (h *RecipesHandler) UpdateRecipe(ctx *fiber.Ctx) error {
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	payload := new(dto.UpdateRecipeDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateRecipe(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *RecipesHandler) DeleteRecipe(ctx *fiber.Ctx) error {
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	if err := h.svc.DeleteRecipe(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *RecipesHandler) RestoreRecipe(ctx *fiber.Ctx) error {
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	if err := h.svc.RestoreRecipe(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

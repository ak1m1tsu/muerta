package recipe

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
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
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *RecipesHandler) FindRecipeByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"recipe": result},
	))
}

func (h *RecipesHandler) FindRecipes(ctx *fiber.Ctx) error {
	filter := new(dto.RecipeFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindRecipes(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"recipes": result, "count": count},
	))
}

func (h *RecipesHandler) UpdateRecipe(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
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
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *RecipesHandler) DeleteRecipe(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	if err := h.svc.DeleteRecipe(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *RecipesHandler) RestoreRecipe(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	if err := h.svc.RestoreRecipe(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *RecipesHandler) FindRecipeIngredients(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeIngredients(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"ingredients": result},
	))
}

func (h *RecipesHandler) CreateRecipeIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	var payload *dto.CreateRecipeIngredientDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.CreateRecipeIngredient(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"ingredient": result},
	))
}

func (h *RecipesHandler) UpdateRecipeIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	var payload *dto.UpdateRecipeIngredientDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.UpdateRecipeIngredient(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"ingredient": result},
	))
}

func (h *RecipesHandler) DeleteRecipeIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	var payload *dto.DeleteRecipeIngredientDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.DeleteRecipeIngredient(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *RecipesHandler) FindRecipeSteps(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeSteps(ctx.Context(), recipeId)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"steps": result},
	))
}

func (h *RecipesHandler) CreateRecipeStep(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	stepId := ctx.Locals(context.StepID).(int)
	var payload *dto.CreateRecipeStepDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.CreateRecipeStep(ctx.Context(), recipeId, stepId, payload.Place)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"step": result},
	))
}

func (h *RecipesHandler) DeleteRecipeStep(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	stepId := ctx.Locals(context.StepID).(int)
	var payload *dto.DeleteRecipeStepDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.DeleteRecipeStep(ctx.Context(), recipeId, stepId, payload.Place); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

package recipe

import (
	"net/http"

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

// CreateRecipe creates a new recipe with the provided data.
//
//	@Summary		Create a new recipe
//	@Description	Creates a new recipe with the provided data
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateRecipeDTO	true	"Recipe data"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/recipes [post]
func (h *RecipesHandler) CreateRecipe(ctx *fiber.Ctx) error {
	var payload *dto.CreateRecipeDTO
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateRecipe(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindRecipeByID finds a recipe by its ID.
//
//	@Summary		Find recipe by ID
//	@Description	Find a recipe by its ID
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int	true	"Recipe ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id} [get]
func (h *RecipesHandler) FindRecipeByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusNotFound).
			JSON(handlers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"recipe": result}})
}

// FindRecipes finds recipes based on the provided filter.
//
//	@Summary		Find recipes
//	@Description	Find recipes based on the provided filter
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.RecipeFilterDTO	false	"Filter recipes by name, category, or ingredients"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/recipes [get]
func (h *RecipesHandler) FindRecipes(ctx *fiber.Ctx) error {
	filter := new(dto.RecipeFilterDTO)
	if err := common.ParseFilterAndValidate(ctx, filter); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.FindRecipes(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"recipes": result, "count": count}})
}

// UpdateRecipe updates a recipe by its ID.
//
//	@Summary		Update a recipe
//	@Description	Update a recipe by its ID.
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int					true	"Recipe ID"
//	@Param			payload		body		dto.UpdateRecipeDTO	true	"Recipe data to update"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id} [put]
func (h *RecipesHandler) UpdateRecipe(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	payload := new(dto.UpdateRecipeDTO)
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateRecipe(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// DeleteRecipe deletes a recipe by ID.
//
//	@Summary		Delete a recipe
//	@Description	Delete a recipe by ID
//	@ID				delete-recipe-by-id
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int	true	"Recipe ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id} [delete]
func (h *RecipesHandler) DeleteRecipe(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	if err := h.svc.DeleteRecipe(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// RestoreRecipe restores a deleted recipe by ID.
//
//	@Summary		Restore a deleted recipe by ID
//	@Description	Restores a deleted recipe by ID
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int	true	"Recipe ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/restore [patch]
func (h *RecipesHandler) RestoreRecipe(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	if err := h.svc.RestoreRecipe(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *RecipesHandler) FindRecipeIngredients(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeIngredients(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"ingredients": result}})
}

// CreateIngredient creates an ingredient of a recipe.
//
//	@Summary		Find recipe ingredients by recipe ID
//	@Description	Find recipe ingredients by recipe ID
//	@Tags			Recipes
//	@Param			recipe_id	path	int	true	"Recipe ID"
//	@Produce		json
//	@Success		200	{object}	handlers.HTTPSuccess
//	@Failure		502	{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/ingredients [get]
func (h *RecipesHandler) CreateIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	var payload *dto.CreateIngredientDTO
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.CreateIngredient(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"ingredient": result}})
}

// UpdateIngredient updates an ingredient of a recipe.
//
//	@Summary		Update an ingredient of a recipe
//	@Description	Update an ingredient of a recipe
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int						true	"Recipe ID"
//	@Param			payload		body		dto.UpdateIngredientDTO	true	"Update ingredient payload"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/ingredients [put]
func (h *RecipesHandler) UpdateIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	var payload *dto.UpdateIngredientDTO
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.UpdateIngredient(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"ingredient": result}})
}

// DeleteIngredient deletes an ingredient from a recipe.
//
//	@Summary		Delete an ingredient from a recipe
//	@Description	Delete an ingredient from a recipe
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int						true	"Recipe ID"
//	@Param			payload		body		dto.DeleteIngredientDTO	true	"Delete ingredient payload"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/ingredients [delete]
func (h *RecipesHandler) DeleteIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	var payload *dto.DeleteIngredientDTO
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.DeleteIngredient(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindRecipeSteps finds all the steps for a recipe.
//
//	@Summary		Find all recipe steps
//	@Description	Find all steps for a recipe by recipe ID
//	@ID				find-recipe-steps
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int	true	"Recipe ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/steps [get]
func (h *RecipesHandler) FindRecipeSteps(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeSteps(ctx.Context(), recipeId)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"steps": result}})
}

// CreateRecipeStep creates a new recipe step for a given recipe and step ID.
//
//	@Summary	Create a recipe step
//	@Tags		Recipes
//	@Param		recipe_id	path		int						true	"Recipe ID"
//	@Param		step_id		path		int						true	"Step ID"
//	@Param		payload		body		dto.CreateRecipeStepDTO	true	"Request body"
//	@Success	200			{object}	handlers.HTTPSuccess
//	@Failure	400			{object}	handlers.HTTPError
//	@Failure	502			{object}	handlers.HTTPError
//	@Router		/recipes/{recipe_id}/steps/{step_id} [post]
func (h *RecipesHandler) CreateRecipeStep(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	stepId := ctx.Locals(context.StepID).(int)
	var payload *dto.CreateRecipeStepDTO
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.CreateRecipeStep(ctx.Context(), recipeId, stepId, payload.Place)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"step": result}})
}

// DeleteRecipeStep deletes a recipe step.
//
//	@Summary		Delete a recipe step.
//	@Description	Deletes the specified recipe step of a recipe.
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int	true	"Recipe ID."
//	@Param			step_id		path		int	true	"Step ID."
//	@Param			place		body		int	true	"Step place."	mininum(1)
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/steps/{step_id} [delete]
func (h *RecipesHandler) DeleteRecipeStep(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	stepId := ctx.Locals(context.StepID).(int)
	var payload *dto.DeleteRecipeStepDTO
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.DeleteRecipeStep(ctx.Context(), recipeId, stepId, payload.Place); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

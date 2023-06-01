package recipe

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	"github.com/romankravchuk/muerta/internal/services/recipe"
)

type RecipesController struct {
	svc recipe.RecipeServicer
	log logger.Logger
}

func New(svc recipe.RecipeServicer, log logger.Logger) *RecipesController {
	return &RecipesController{
		svc: svc,
		log: log,
	}
}

// Create creates a new recipe with the provided data.
//
//	@Summary		Create a new recipe
//	@Description	Creates a new recipe with the provided data
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateRecipe	true	"Recipe data"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/recipes [post]
//	@Security		Bearer
func (h *RecipesController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateRecipe)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateRecipe(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// FindOne finds a recipe by its ID.
//
//	@Summary		Find recipe by ID
//	@Description	Finds a recipe by its ID
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int	true	"Recipe ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id} [get]
func (h *RecipesController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusNotFound).
			JSON(controllers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{
		Success: true,
		Data:    controllers.Data{"recipe": result},
	})
}

// FindMany finds recipes based on the provided filter.
//
//	@Summary		Find recipes
//	@Description	Finds recipes based on the provided filter
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.RecipeFilter	false	"Filter recipes by name, category, or ingredients"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/recipes [get]
func (h *RecipesController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.RecipeFilter)
	if err := utils.ParseFilterAndValidate(ctx, filter); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.FindRecipes(ctx.Context(), filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(
		controllers.HTTPSuccess{
			Success: true,
			Data:    controllers.Data{"recipes": result, "count": count},
		},
	)
}

// Update updates a recipe by its ID.
//
//	@Summary		Update a recipe
//	@Description	Update a recipe by its ID.
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int					true	"Recipe ID"
//	@Param			payload		body		dto.UpdateRecipe	true	"Recipe data to update"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id} [put]
//	@Security		Bearer
func (h *RecipesController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	payload := new(params.UpdateRecipe)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateRecipe(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Delete deletes a recipe by ID.
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
//	@Security		Bearer
func (h *RecipesController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	if err := h.svc.DeleteRecipe(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Restore restores a deleted recipe by ID.
//
//	@Summary		Restore a deleted recipe by ID
//	@Description	Restores a deleted recipe by ID
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int	true	"Recipe ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id} [patch]
//	@Security		Bearer
func (h *RecipesController) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	if err := h.svc.RestoreRecipe(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

func (h *RecipesController) FindRecipeIngredients(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeIngredients(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{
		Success: true,
		Data:    controllers.Data{"ingredients": result},
	})
}

// AddIngredient adds an ingredient to a recipe.
//
//	@Summary		Add an ingredient to a recipe
//	@Description	Adds n ingredient to a recipe
//	@Tags			Recipes
//	@Produce		json
//	@Accept			json
//	@Param			recipe_id	path		int						true	"Recipe ID"
//	@Param			payload		body		dto.CreateIngredient	true	"Create ingredient payload"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/ingredients [get]
func (h *RecipesController) AddIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	payload := new(params.CreateIngredient)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.CreateIngredient(ctx.Context(), id, payload)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"ingredient": result}})
}

// UpdateIngredient updates an ingredient of a recipe.
//
//	@Summary		Update an ingredient of a recipe
//	@Description	Update an ingredient of a recipe
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int						true	"Recipe ID"
//	@Param			payload		body		dto.UpdateIngredient	true	"Update ingredient payload"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/ingredients [put]
//	@Security		Bearer
func (h *RecipesController) UpdateIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	payload := new(params.UpdateIngredient)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.UpdateIngredient(ctx.Context(), id, payload)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"ingredient": result}})
}

// RemoveIngredient removes an ingredient from a recipe.
//
//	@Summary		Remove an ingredient from a recipe
//	@Description	Removes an ingredient from a recipe
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int						true	"Recipe ID"
//	@Param			payload		body		dto.DeleteIngredient	true	"Delete ingredient payload"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/ingredients [delete]
//	@Security		Bearer
func (h *RecipesController) RemoveIngredient(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RecipeID).(int)
	payload := new(params.DeleteIngredient)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.DeleteIngredient(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// FindSteps finds all the steps for a recipe.
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
func (h *RecipesController) FindSteps(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	result, err := h.svc.FindRecipeSteps(ctx.Context(), recipeId)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"steps": result}})
}

// AddStep creates a new recipe step for a given recipe and step ID.
//
//	@Summary	Create a recipe step
//	@Tags		Recipes
//	@Param		recipe_id	path		int						true	"Recipe ID"
//	@Param		step_id		path		int						true	"Step ID"
//	@Param		payload		body		dto.CreateRecipeStep	true	"Request body"
//	@Success	200			{object}	handlers.HTTPSuccess
//	@Failure	400			{object}	handlers.HTTPError
//	@Failure	502			{object}	handlers.HTTPError
//	@Router		/recipes/{recipe_id}/steps/{step_id} [post]
//	@Security	Bearer
func (h *RecipesController) AddStep(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	stepId := ctx.Locals(context.StepID).(int)
	payload := new(params.CreateRecipeStep)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.CreateRecipeStep(ctx.Context(), recipeId, stepId, payload.Place)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"step": result}})
}

// RemoveStep removes a recipe step.
//
//	@Summary		Remove a recipe step.
//	@Description	Removes the specified recipe step of a recipe.
//	@Tags			Recipes
//	@Accept			json
//	@Produce		json
//	@Param			recipe_id	path		int						true	"Recipe ID."
//	@Param			step_id		path		int						true	"Step ID."
//	@Param			playload	body		dto.DeleteRecipeStep	true	"Request body"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/recipes/{recipe_id}/steps/{step_id} [delete]
//	@Security		Bearer
func (h *RecipesController) RemoveStep(ctx *fiber.Ctx) error {
	recipeId := ctx.Locals(context.RecipeID).(int)
	stepId := ctx.Locals(context.StepID).(int)
	payload := new(params.DeleteRecipeStep)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.DeleteRecipeStep(ctx.Context(), recipeId, stepId, payload.Place); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

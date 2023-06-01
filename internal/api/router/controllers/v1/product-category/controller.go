package productcategory

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/category"
)

type ProductCategoryController struct {
	svc service.CategoryServicer
	log logger.Logger
}

func New(svc service.CategoryServicer, log logger.Logger) *ProductCategoryController {
	return &ProductCategoryController{
		svc: svc,
		log: log,
	}
}

// Create creates a new product category
//
//	@Summary		Create product category
//	@Description	Creates a new product category
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateProductCategory	true	"Payload for creating product category"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/product-categories [post]
//	@Security		Bearer
func (h *ProductCategoryController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateProductCategory)
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
	if err := h.svc.CreateCategory(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// FindOne finds a product category by ID
//
//	@Summary		Find a product category by ID
//	@Description	Finds a product category by ID
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int	true	"Category ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/product-categories/{category_id} [get]
func (h *ProductCategoryController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	result, err := h.svc.FindCategoryByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusNotFound).
			JSON(controllers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{
		Success: true,
		Data:    controllers.Data{"category": result},
	})
}

// FindMany finds product categories with optional filters
//
//	@Summary		Find product categories
//	@Description	Finds product categories with optional filters
//	@ID				find-product-categories
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.ProductCategoryFilter	false	"Filter criteria for product categories"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/product-categories [get]
func (h *ProductCategoryController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.ProductCategoryFilter)
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
	result, err := h.svc.FindCategorys(ctx.Context(), filter)
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
			Data:    controllers.Data{"categories": result, "count": count},
		},
	)
}

// Update updates an existing product category by providing the ID and updated fields in the request body.
//
//	@Summary		Update an existing product category by ID
//	@Description	Updates an existing product category by providing the ID and updated fields in the request body
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int							true	"Product Category ID"
//	@Param			payload		body		dto.UpdateProductCategory	true	"Updated Product Category Fields"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/product-categories/{category_id} [put]
//	@Security		Bearer
func (h *ProductCategoryController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	payload := new(params.UpdateProductCategory)
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
	if err := h.svc.UpdateCategory(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Delete deletes a product category by ID
//
//	@Summary		Delete a product category
//	@Description	Deletes a product category by ID
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int	true	"Category ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/product-categories/{category_id} [delete]
//	@Security		Bearer
func (h *ProductCategoryController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.DeleteCategory(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Restore restores a previously deleted product category with the given ID
//
//	@Summary		Restore a deleted product category
//	@Description	Restores a previously deleted product category with the given ID
//	@Tags			Product Categories
//	@Param			category_id	path		int	true	"Product category ID to be restored"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/product-categories/{category_id} [patch]
//	@Security		Bearer
func (h *ProductCategoryController) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.RestoreCategory(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

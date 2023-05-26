package productcategory

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/category"
)

type CategoryHandler struct {
	svc service.CategoryServicer
	log logger.Logger
}

func New(svc service.CategoryServicer, log logger.Logger) *CategoryHandler {
	return &CategoryHandler{
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
func (h *CategoryHandler) Create(ctx *fiber.Ctx) error {
	payload := new(dto.CreateProductCategory)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateCategory(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *CategoryHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	result, err := h.svc.FindCategoryByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusNotFound).
			JSON(handlers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data:    handlers.Data{"category": result},
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
func (h *CategoryHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.ProductCategoryFilter)
	if err := common.ParseFilterAndValidate(ctx, filter); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.FindCategorys(ctx.Context(), filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(
		handlers.HTTPSuccess{
			Success: true,
			Data:    handlers.Data{"categories": result, "count": count},
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
func (h *CategoryHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	payload := new(dto.UpdateProductCategory)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateCategory(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *CategoryHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.DeleteCategory(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *CategoryHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.RestoreCategory(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

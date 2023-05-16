package productcategory

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/category"
)

type CategoryHandler struct {
	svc service.CategoryServicer
	log *log.Logger
}

func New(svc service.CategoryServicer, log *log.Logger) *CategoryHandler {
	return &CategoryHandler{
		svc: svc,
		log: log,
	}
}

// CreateProductCategory creates a new product category
//
//	@Summary		Create product category
//	@Description	Create a new product category
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateProductCategoryDTO	true	"Payload for creating product category"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/product-categories [post]
func (h *CategoryHandler) CreateProductCategory(ctx *fiber.Ctx) error {
	var payload *dto.CreateProductCategoryDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateCategory(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindProductCategoryByID finds a product category by ID
//
//	@Summary		Find a product category by ID
//	@Description	Get a product category by ID
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int	true	"Category ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/product-categories/{category_id} [get]
func (h *CategoryHandler) FindProductCategoryByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	result, err := h.svc.FindCategoryByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusNotFound).
			JSON(handlers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true,
		Data: handlers.Data{"category": result},
	})
}

// FindProductCategories finds product categories with optional filters
//
//	@Summary		Find product categories
//	@Description	Find product categories with optional filters
//	@ID				find-product-categories
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.ProductCategoryFilterDTO	false	"Filter criteria for product categories"
//	@Param			limit	query		int								false	"Limit the number of results returned"
//	@Param			offset	query		int								false	"Offset for pagination"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/product-categories [get]
func (h *CategoryHandler) FindProductCategories(ctx *fiber.Ctx) error {
	filter := new(dto.ProductCategoryFilterDTO)
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
	result, err := h.svc.FindCategorys(ctx.Context(), filter)
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
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"categories": result, "count": count}})
}

// UpdateProductCategory updates an existing product category by providing the ID and updated fields in the request body.
//
//	@Summary		Update an existing product category by ID
//	@Description	Update an existing product category by providing the ID and updated fields in the request body
//	@Tags			Product Categories
//	@Accept			json
//	@Produce		json
//	@Param			category_id	path		int								true	"Product Category ID"
//	@Param			payload		body		dto.UpdateProductCategoryDTO	true	"Updated Product Category Fields"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/product-categories/{category_id} [put]
func (h *CategoryHandler) UpdateProductCategory(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	payload := new(dto.UpdateProductCategoryDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateCategory(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// DeleteProductCategroy deletes a product category by ID
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
func (h *CategoryHandler) DeleteProductCategroy(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.DeleteCategory(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// RestoreProductCategroy restores a previously deleted product category with the given ID
//
//	@Summary		Restore a deleted product category
//	@Description	Restores a previously deleted product category with the given ID
//	@Tags			Product Categories
//	@Param			category_id	path		int	true	"Product category ID to be restored"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/product-categories/{category_id} [patch]
func (h *CategoryHandler) RestoreProductCategroy(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.RestoreCategory(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

package product

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/product"
)

type ProductHandler struct {
	svc service.ProductServicer
	log logger.Logger
}

func New(svc service.ProductServicer, log logger.Logger) *ProductHandler {
	return &ProductHandler{
		svc: svc,
		log: log,
	}
}

// Create creates a new product
//
//	@Summary		Create a new product
//	@Description	Create a new product with the given details
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateProduct	true	"Product details"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/products [post]
//	@Security		Bearer
func (h *ProductHandler) Create(ctx *fiber.Ctx) error {
	payload := new(dto.CreateProduct)
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
	if err := h.svc.CreateProduct(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindOne finds a product by id
//
//	@Summary		Get a product by ID
//	@Description	Retrieve the details of a product with the specified ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id} [get]
func (h *ProductHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	result, err := h.svc.FindProductByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data:    handlers.Data{"product": result},
	})
}

// FindMany finds products by filter
//
//	@Summary		Get a list of products
//	@Description	Retrieve a list of products with optional filters
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.ProductFilter	false	"Product filter parameters"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/products [get]
func (h *ProductHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.ProductFilter)
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
	result, err := h.svc.FindProducts(ctx.Context(), filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data:    handlers.Data{"products": result, "count": count},
	})
}

// Update updates a product
//
//	@Summary		Update a product
//	@Description	Update an existing product with new details
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int					true	"Product ID"
//	@Param			payload		body		dto.UpdateProduct	true	"New product details"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id} [put]
//	@Security		Bearer
func (h *ProductHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	payload := new(dto.UpdateProduct)
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
	if err := h.svc.UpdateProduct(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Delete deletes a product
//
//	@Summary		Delete a product
//	@Description	Delete an existing product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id} [delete]
//	@Security		Bearer
func (h *ProductHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	if err := h.svc.DeleteProduct(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Restore restores a product
//
//	@Summary		Restore a deleted product
//	@Description	Restore a deleted product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id}/ [patch]
//	@Security		Bearer
func (h *ProductHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	if err := h.svc.RestoreProduct(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindCategories finds product categories
//
//	@Summary		Get categories of a product
//	@Description	Get the categories of a product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id}/categories [get]
func (h *ProductHandler) FindCategories(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	categories, err := h.svc.FindProductCategories(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Data: handlers.Data{"categories": categories}})
}

// FindRecipes finds product recipes
//
//	@Summary		Get recipes of a product
//	@Description	Get the recipes of a product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id}/recipes [get]
func (h *ProductHandler) FindRecipes(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	recipes, err := h.svc.FindProductRecipes(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Data: handlers.Data{"recipes": recipes}})
}

// AddCategory adds category to product
//
//	@Summary		Add category to product
//	@Description	Adds category to product given the product ID and category ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Param			category_id	path		int	true	"Category ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id}/categories/{category_id} [post]
//	@Security		Bearer
func (h *ProductHandler) AddCategory(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	categoryID := ctx.Locals(context.CategoryID).(int)
	result, err := h.svc.CreateCategory(ctx.Context(), productID, categoryID)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})

	}
	return ctx.JSON(handlers.HTTPSuccess{
		Data: handlers.Data{"category": result},
	})
}

// RemoveCategory removes category from product
//
//	@Summary		Remove a category from a product
//	@Description	Removes a category from a product given the product ID and category ID
//	@Tags			Products
//	@Param			product_id	path		integer	true	"Product ID"
//	@Param			category_id	path		integer	true	"Category ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id}/categories/{category_id} [delete]
//	@Security		Bearer
func (h *ProductHandler) RemoveCategory(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	categoryID := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.DeleteCategory(ctx.Context(), productID, categoryID); err != nil {
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindTips finds product tips
//
//	@Summary		Find tips for a product
//	@Description	Finds tips for a product given the product ID
//	@Tags			Products
//	@Param			product_id	path		integer	true	"Product ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id}/tips [get]
func (h *ProductHandler) FindTips(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	result, err := h.svc.FindProductTips(ctx.Context(), productID)
	if err != nil {
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data:    handlers.Data{"tips": result},
	})
}

// AddTip adds tip to product
//
//	@Summary		Add a tip for a product
//	@Description	Adds a tip for a product given the product ID and tip ID
//	@Tags			Products
//	@Param			product_id	path		integer	true	"Product ID"
//	@Param			tip_id		path		integer	true	"Tip ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id}/tip/{tip_id} [post]
//	@Security		Bearer
func (h *ProductHandler) AddTip(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	result, err := h.svc.CreateProductTip(ctx.Context(), productID, tipID)
	if err != nil {
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data:    handlers.Data{"tip": result},
	})
}

// RemoveTip removes tip from product
//
//	@Summary		Remove a tip from a product
//	@Description	Removes a tip from a product given the product ID and tip ID
//	@Tags			Products
//	@Param			product_id	path		integer	true	"Product ID"
//	@Param			tip_id		path		integer	true	"Tip ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/products/{product_id}/tip/{tip_id} [delete]
//	@Security		Bearer
func (h *ProductHandler) RemoveTip(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	err := h.svc.DeleteProductTip(ctx.Context(), productID, tipID)
	if err != nil {
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

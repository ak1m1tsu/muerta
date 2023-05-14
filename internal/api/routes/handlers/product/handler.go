package product

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/product"
)

type ProductHandler struct {
	svc service.ProductServicer
	log *log.Logger
}

func New(svc service.ProductServicer, log *log.Logger) *ProductHandler {
	return &ProductHandler{
		svc: svc,
		log: log,
	}
}

// CreateProduct creates a new product
//
//	@Summary		Create a new product
//	@Description	Create a new product with the given details
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product	body		dto.CreateProductDTO	true	"Product details"
//	@Success		200		{object}	handlers.SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		502		{object}	ErrorResponse
//	@Router			/products [post]
func (h *ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	var payload *dto.CreateProductDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	if err := h.svc.CreateProduct(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse())
}

// FindProductByID finds a product by id
//
//	@Summary		Get a product by ID
//	@Description	Retrieve the details of a product with the specified ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		502			{object}	ErrorResponse
//	@Router			/products/{product_id} [get]
func (h *ProductHandler) FindProductByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	result, err := h.svc.FindProductByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"product": result},
	))
}

// FindProducts finds products by filter
//
//	@Summary		Get a list of products
//	@Description	Retrieve a list of products with optional filters
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.ProductFilterDTO	false	"Product filter parameters"
//	@Success		200		{object}	handlers.SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		502		{object}	ErrorResponse
//	@Router			/products [get]
func (h *ProductHandler) FindProducts(ctx *fiber.Ctx) error {
	filter := new(dto.ProductFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	result, err := h.svc.FindProducts(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"products": result, "count": count},
	))
}

// UpdateProduct updates a product
//
//	@Summary		Update a product
//	@Description	Update an existing product with new details
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int						true	"Product ID"
//	@Param			product		body		dto.UpdateProductDTO	true	"New product details"
//	@Success		200			{object}	fiber.Map
//	@Failure		400			{object}	ErrorResponse
//	@Failure		502			{object}	ErrorResponse
//	@Router			/products/{product_id} [put]
func (h *ProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	payload := new(dto.UpdateProductDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).JSON(handlers.ErrorResponse(fiber.ErrBadRequest))
	}
	if err := h.svc.UpdateProduct(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse())
}

// DeleteProduct deletes a product
//
//	@Summary		Delete a product
//	@Description	Delete an existing product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	ErrorResponse
//	@Router			/products/{product_id} [delete]
func (h *ProductHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	if err := h.svc.DeleteProduct(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse())
}

// RestoreProduct restores a product
//
//	@Summary		Restore a deleted product
//	@Description	Restore a deleted product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	ErrorResponse
//	@Router			/products/{product_id}/ [patch]
func (h *ProductHandler) RestoreProduct(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	if err := h.svc.RestoreProduct(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse())
}

// FindProductCategories finds product categories
//
//	@Summary		Get categories of a product
//	@Description	Get the categories of a product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	ErrorResponse
//	@Router			/products/{product_id}/categories [get]
func (h *ProductHandler) FindProductCategories(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	categories, err := h.svc.FindProductCategories(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"categories": categories},
	))
}

// FindProductRecipes finds product recipes
//
//	@Summary		Get recipes of a product
//	@Description	Get the recipes of a product by ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	ErrorResponse
//	@Router			/products/{product_id}/recipes [get]
func (h *ProductHandler) FindProductRecipes(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	recipes, err := h.svc.FindProductRecipes(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"recipes": recipes},
	))
}

// CreateCategory adds category to product
//
//	@Summary		Create a category for a product
//	@Description	Creates a new category for a product by product ID and category ID
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			product_id	path		int	true	"Product ID"
//	@Param			category_id	path		int	true	"Category ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	handlers.ErrorResponse
//	@Router			/api/v1/products/{product_id}/categories/{category_id} [post]
func (h *ProductHandler) CreateCategory(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	categoryID := ctx.Locals(context.CategoryID).(int)
	result, err := h.svc.CreateCategory(ctx.Context(), productID, categoryID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))

	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"category": result},
	))
}

// DeleteCategory removes category from product
//
//	@Summary		Delete a category from a product
//	@Description	Deletes a category from a product given the product ID and category ID
//	@Tags			Products
//	@Param			product_id	path		integer	true	"Product ID"
//	@Param			category_id	path		integer	true	"Category ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	handlers.ErrorResponse
//	@Router			/product/{product_id}/category/{category_id} [delete]
func (h *ProductHandler) DeleteCategory(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	categoryID := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.DeleteCategory(ctx.Context(), productID, categoryID); err != nil {
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse())
}

// FindProductTips finds product tips
//
//	@Summary		Find tips for a product
//	@Description	Finds tips for a product given the product ID
//	@Tags			Products
//	@Param			product_id	path		integer	true	"Product ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	handlers.ErrorResponse
//	@Router			/product/{product_id}/tips [get]
func (h *ProductHandler) FindProductTips(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	result, err := h.svc.FindProductTips(ctx.Context(), productID)
	if err != nil {
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"tips": result},
	))
}

// CreateProductTip adds tip to product
//
//	@Summary		Create a tip for a product
//	@Description	Creates a tip for a product given the product ID and tip ID
//	@Tags			Products
//	@Param			product_id	path		integer	true	"Product ID"
//	@Param			tip_id		path		integer	true	"Tip ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	handlers.ErrorResponse
//	@Router			/product/{product_id}/tip/{tip_id} [post]
func (h *ProductHandler) CreateProductTip(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	result, err := h.svc.CreateProductTip(ctx.Context(), productID, tipID)
	if err != nil {
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"tip": result},
	))
}

// DeleteProductTip removes tip from product
//
//	@Summary		Delete a tip from a product
//	@Description	Deletes a tip from a product given the product ID and tip ID
//	@Tags			Product
//	@Param			product_id	path		integer	true	"Product ID"
//	@Param			tip_id		path		integer	true	"Tip ID"
//	@Success		200			{object}	handlers.SuccessResponse
//	@Failure		502			{object}	handlers.ErrorResponse
//	@Router			/product/{product_id}/tip/{tip_id} [delete]
func (h *ProductHandler) DeleteProductTip(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	err := h.svc.DeleteProductTip(ctx.Context(), productID, tipID)
	if err != nil {
		return ctx.Status(http.StatusBadGateway).JSON(handlers.ErrorResponse(fiber.ErrBadGateway))
	}
	return ctx.JSON(handlers.SuccessResponse())
}

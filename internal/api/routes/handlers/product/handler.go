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

func (h *ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	var payload *dto.CreateProductDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateProduct(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *ProductHandler) FindProductByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	result, err := h.svc.FindProductByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"product": result},
	))
}

func (h *ProductHandler) FindProducts(ctx *fiber.Ctx) error {
	filter := new(dto.ProductFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindProducts(ctx.Context(), filter)
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
		handlers.Data{"products": result, "count": count},
	))
}

func (h *ProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	payload := new(dto.UpdateProductDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateProduct(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *ProductHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	if err := h.svc.DeleteProduct(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *ProductHandler) RestoreProduct(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	if err := h.svc.RestoreProduct(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *ProductHandler) FindProductCategories(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	categories, err := h.svc.FindProductCategories(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"categories": categories},
	})
}

func (h *ProductHandler) FindProductRecipes(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ProductID).(int)
	recipes, err := h.svc.FindProductRecipes(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"recipes": recipes},
	})
}

func (h *ProductHandler) AddProductCategory(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	categoryID := ctx.Locals(context.CategoryID).(int)
	result, err := h.svc.AddProductCategory(ctx.Context(), productID, categoryID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"category": result},
	))
}
func (h *ProductHandler) RemoveProductCategory(ctx *fiber.Ctx) error {
	productID := ctx.Locals(context.ProductID).(int)
	categoryID := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.RemoveProductCategory(ctx.Context(), productID, categoryID); err != nil {
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.SuccessResponse())
}

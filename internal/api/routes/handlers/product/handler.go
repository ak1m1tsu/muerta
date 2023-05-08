package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
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
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *ProductHandler) FindProductByID(ctx *fiber.Ctx) error {
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	dto, err := h.svc.FindProductByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"recipe": dto},
	})
}

func (h *ProductHandler) FindProducts(ctx *fiber.Ctx) error {
	filter := new(dto.ProductFilterDTO)
	if err := common.GetProductFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	dtos, err := h.svc.FindProducts(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"recipes": dtos},
	})
}

func (h *ProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
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
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	if err := h.svc.DeleteProduct(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *ProductHandler) RestoreProduct(ctx *fiber.Ctx) error {
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	if err := h.svc.RestoreProduct(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

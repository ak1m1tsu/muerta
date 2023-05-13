package productcategory

import (
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

func (h *CategoryHandler) Create(ctx *fiber.Ctx) error {
	var payload *dto.CreateProductCategoryDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateCategory(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *CategoryHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	result, err := h.svc.FindCategoryByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"category": result},
	))
}

func (h *CategoryHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.ProductCategoryFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindCategorys(ctx.Context(), filter)
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
		handlers.Data{"categories": result, "count": count},
	))
}

func (h *CategoryHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	payload := new(dto.UpdateProductCategoryDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateCategory(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *CategoryHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.DeleteCategory(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *CategoryHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.CategoryID).(int)
	if err := h.svc.RestoreCategory(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

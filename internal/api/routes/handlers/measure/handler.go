package measure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/measure"
)

type MeasureHandler struct {
	svc service.MeasureServicer
	log *log.Logger
}

func New(svc service.MeasureServicer, log *log.Logger) MeasureHandler {
	return MeasureHandler{
		svc: svc,
		log: log,
	}
}

func (h *MeasureHandler) CreateMeasure(ctx *fiber.Ctx) error {
	var payload *dto.CreateMeasureDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateMeasure(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *MeasureHandler) FindMeasureByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	dto, err := h.svc.FindMeasureByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"measure": dto},
	))
}

func (h *MeasureHandler) FindMeasures(ctx *fiber.Ctx) error {
	filter := new(dto.MeasureFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindMeasures(ctx.Context(), filter)
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
		handlers.Data{"measures": result, "count": count},
	))
}

func (h *MeasureHandler) UpdateMeasure(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	payload := new(dto.UpdateMeasureDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateMeasure(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *MeasureHandler) DeleteMeasure(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	if err := h.svc.DeleteMeasure(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

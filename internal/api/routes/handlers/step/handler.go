package step

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/step"
)

type StepHandler struct {
	svc service.StepServicer
	log *log.Logger
}

func New(svc service.StepServicer, log *log.Logger) *StepHandler {
	return &StepHandler{
		svc: svc,
		log: log,
	}
}

func (h *StepHandler) FindSteps(ctx *fiber.Ctx) error {
	filter := new(dto.StepFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindSteps(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"steps": result, "count": count},
	))
}

func (h *StepHandler) CreateStep(ctx *fiber.Ctx) error {
	var paylaod *dto.CreateStepDTO
	if err := ctx.BodyParser(paylaod); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(paylaod); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.CreateStep(ctx.Context(), paylaod)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"step": result},
	))
}

func (h *StepHandler) FindStep(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	result, err := h.svc.FindStep(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"step": result},
	))
}

func (h *StepHandler) UpdateStep(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	var payload *dto.UpdateStepDTO
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.UpdateStep(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"step": result},
	))
}

func (h *StepHandler) DeleteStep(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	if err := h.svc.DeleteStep(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *StepHandler) RestoreStep(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	result, err := h.svc.RestoreStep(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"step": result},
	))
}

package tip

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/tip"
)

type TipHandler struct {
	svc service.TipServicer
	log *log.Logger
}

func New(svc service.TipServicer, log *log.Logger) *TipHandler {
	return &TipHandler{
		svc: svc,
		log: log,
	}
}

func (h *TipHandler) CreateTip(ctx *fiber.Ctx) error {
	var payload *dto.CreateTipDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateTip(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *TipHandler) FindTipByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"tip": result},
	))
}

func (h *TipHandler) FindTips(ctx *fiber.Ctx) error {
	filter := new(dto.TipFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindTips(ctx.Context(), filter)
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
		handlers.Data{"tips": result, "count": count},
	))
}

func (h *TipHandler) UpdateTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	payload := new(dto.UpdateTipDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateTip(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *TipHandler) DeleteTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *TipHandler) RestoreTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	if err := h.svc.RestoreTip(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *TipHandler) FindTipStorages(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipStorages(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"storages": result},
	))
}

func (h *TipHandler) FindTipProducts(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipProducts(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"products": result},
	))
}

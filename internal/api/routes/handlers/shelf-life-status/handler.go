package shelflifestatus

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life-status"
)

type ShelfLifeStatusHandler struct {
	svc service.ShelfLifeStatusServicer
	log *log.Logger
}

func New(svc service.ShelfLifeStatusServicer, log *log.Logger) ShelfLifeStatusHandler {
	return ShelfLifeStatusHandler{
		svc: svc,
		log: log,
	}
}

func (h *ShelfLifeStatusHandler) CreateShelfLifeStatus(ctx *fiber.Ctx) error {
	var payload *dto.CreateShelfLifeStatusDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateShelfLifeStatus(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *ShelfLifeStatusHandler) FindShelfLifeStatusByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	result, err := h.svc.FindShelfLifeStatusByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"status": result},
	))
}

func (h *ShelfLifeStatusHandler) FindShelfLifeStatuses(ctx *fiber.Ctx) error {
	filter := new(dto.ShelfLifeStatusFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindShelfLifeStatuss(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"statues": result, "count": count},
	))
}

func (h *ShelfLifeStatusHandler) UpdateShelfLifeStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	payload := new(dto.UpdateShelfLifeStatusDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateShelfLifeStatus(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *ShelfLifeStatusHandler) DeleteShelfLifeStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	if err := h.svc.DeleteShelfLifeStatus(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

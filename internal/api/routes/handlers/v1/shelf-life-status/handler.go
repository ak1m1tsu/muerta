package shelflifestatus

import (
	"net/http"

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
	var payload *dto.CreateShelfLifeStatus
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateShelfLifeStatus(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *ShelfLifeStatusHandler) FindShelfLifeStatusByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	result, err := h.svc.FindShelfLifeStatusByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"status": result}})
}

func (h *ShelfLifeStatusHandler) FindShelfLifeStatuses(ctx *fiber.Ctx) error {
	filter := new(dto.ShelfLifeStatusFilter)
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
	result, err := h.svc.FindShelfLifeStatuss(ctx.Context(), filter)
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
	return ctx.JSON(
		handlers.HTTPSuccess{Success: true, Data: handlers.Data{"statues": result, "count": count}},
	)
}

func (h *ShelfLifeStatusHandler) UpdateShelfLifeStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	payload := new(dto.UpdateShelfLifeStatus)
	if err := common.ParseBodyAndValidate(ctx, &payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateShelfLifeStatus(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *ShelfLifeStatusHandler) DeleteShelfLifeStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	if err := h.svc.DeleteShelfLifeStatus(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

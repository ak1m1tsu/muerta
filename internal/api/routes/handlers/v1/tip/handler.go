package tip

import (
	"net/http"

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
	var payload *dto.CreateTip
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
	result, err := h.svc.CreateTip(ctx.Context(), payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tip": result}})
}

func (h *TipHandler) FindTipByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tip": result}})
}

func (h *TipHandler) FindTips(ctx *fiber.Ctx) error {
	filter := new(dto.TipFilter)
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
	result, err := h.svc.FindTips(ctx.Context(), filter)
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
		handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tips": result, "count": count}},
	)
}

func (h *TipHandler) UpdateTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	payload := new(dto.UpdateTip)
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
	if err := h.svc.UpdateTip(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *TipHandler) DeleteTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *TipHandler) RestoreTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	if err := h.svc.RestoreTip(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *TipHandler) FindTipStorages(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipStorages(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storages": result}})
}

func (h *TipHandler) FindTipProducts(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TipID).(int)
	result, err := h.svc.FindTipProducts(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"products": result}})
}

func (h *TipHandler) AddProductToTip(ctx *fiber.Ctx) error {
	tipID := ctx.Locals(context.TipID).(int)
	productID := ctx.Locals(context.ProductID).(int)
	result, err := h.svc.AddProductToTip(ctx.Context(), tipID, productID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"product": result}})
}

func (h *TipHandler) RemoveProductFromTip(ctx *fiber.Ctx) error {
	tipID := ctx.Locals(context.TipID).(int)
	productID := ctx.Locals(context.ProductID).(int)
	if err := h.svc.RemoveProductFromTip(ctx.Context(), tipID, productID); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *TipHandler) AddStorageToTip(ctx *fiber.Ctx) error {
	tipID := ctx.Locals(context.TipID).(int)
	storateID := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.AddStorageToTip(ctx.Context(), tipID, storateID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storage": result}})
}

func (h *TipHandler) RemoveStorageFromTip(ctx *fiber.Ctx) error {
	tipID := ctx.Locals(context.TipID).(int)
	storateID := ctx.Locals(context.StorageID).(int)
	if err := h.svc.RemoveStorageFromTip(ctx.Context(), tipID, storateID); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

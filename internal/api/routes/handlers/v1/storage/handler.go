package storage

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/storage"
)

type StorageHandler struct {
	svc service.StorageServicer
	log *log.Logger
}

func New(svc service.StorageServicer, log *log.Logger) *StorageHandler {
	return &StorageHandler{
		svc: svc,
		log: log,
	}
}

func (h *StorageHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.StorageFilterDTO)
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
	result, err := h.svc.FindStorages(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storages": result, "count": count}})
}

func (h *StorageHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.FindStorageByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storage": result}})
}

func (h *StorageHandler) Create(ctx *fiber.Ctx) error {
	var payload *dto.CreateStorageDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateStorage(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	var payload *dto.UpdateStorageDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateStorage(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	if err := h.svc.DeleteStorage(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	if err := h.svc.RestoreStorage(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageHandler) FindTips(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.FindTips(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tips": result}})
}

func (h *StorageHandler) CreateTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	result, err := h.svc.CreateTip(ctx.Context(), id, tipID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tip": result}})
}

func (h *StorageHandler) DeleteTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id, tipID); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageHandler) FindShelfLives(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.FindShelfLives(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf-lives": result}})
}

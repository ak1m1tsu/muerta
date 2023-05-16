package storagetype

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/storage-type"
)

type StorageTypeHandler struct {
	svc service.StorageTypeServicer
	log *log.Logger
}

func New(svc service.StorageTypeServicer, log *log.Logger) *StorageTypeHandler {
	return &StorageTypeHandler{
		svc: svc,
		log: log,
	}
}

func (h *StorageTypeHandler) CreateStorageType(ctx *fiber.Ctx) error {
	var payload *dto.CreateStorageTypeDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateStorageType(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageTypeHandler) FindStorageTypeByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindStorageTypeByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"type": result}})
}

func (h *StorageTypeHandler) FindStorageTypes(ctx *fiber.Ctx) error {
	filter := new(dto.StorageTypeFilterDTO)
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
	result, err := h.svc.FindStorageTypes(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"types": result, "count": count}})
}

func (h *StorageTypeHandler) UpdateStorageType(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	payload := new(dto.UpdateStorageTypeDTO)
	if err := ctx.BodyParser(payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateStorageType(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageTypeHandler) DeleteStorageType(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	if err := h.svc.DeleteStorageType(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageTypeHandler) FindStorages(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindStorages(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storages": result}})
}

func (h *StorageTypeHandler) FindTips(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindTips(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tips": result}})
}

func (h *StorageTypeHandler) CreateTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	result, err := h.svc.CreateTip(ctx.Context(), id, tipID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tip": result}})
}

func (h *StorageTypeHandler) DeleteTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id, tipID); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

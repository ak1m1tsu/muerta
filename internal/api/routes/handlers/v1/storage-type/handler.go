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

// Create godoc
//
//	@Summary		Create storage type
//	@Description	Create storage type
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			storage_type	body		dto.CreateStorageType	true	"Storage type"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/storage-types [post]
//	@Security		Bearer
func (h *StorageTypeHandler) Create(ctx *fiber.Ctx) error {
	var payload *dto.CreateStorageType
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
	if err := h.svc.CreateStorageType(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindOne godoc
//
//	@Summary		Find storage type by id
//	@Description	Find storage type by id
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			id_type	path		int	true	"Storage type id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storage-types/{id_type} [get]
func (h *StorageTypeHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindStorageTypeByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"type": result}})
}

// FindMany godoc
//
//	@Summary		Find storage types by filter
//	@Description	Find storage types by filter
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.StorageTypeFilter	true	"Filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storage-types [get]
func (h *StorageTypeHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.StorageTypeFilter)
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
		handlers.HTTPSuccess{Success: true, Data: handlers.Data{"types": result, "count": count}},
	)
}

// Update godoc
//
//	@Summary		Update storage type
//	@Description	Update storage type
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			id_type			path		int						true	"Storage type id"
//	@Param			storage_type	body		dto.UpdateStorageType	true	"Storage type"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		404				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/storage-types/{id_type} [put]
//	@Security		Bearer
func (h *StorageTypeHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	payload := new(dto.UpdateStorageType)
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
	if err := h.svc.UpdateStorageType(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Delete godoc
//
//	@Summary		Delete storage type
//	@Description	Delete storage type
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			id_type	path		int	true	"Storage type id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storage-types/{id_type} [delete]
//	@Security		Bearer
func (h *StorageTypeHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	if err := h.svc.DeleteStorageType(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindStorages godoc
//
//	@Summary		Find storages by storage type id
//	@Description	Find storages by storage type id
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			id_type	path		int	true	"Storage type id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storage-types/{id_type}/storages [get]
func (h *StorageTypeHandler) FindStorages(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindStorages(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storages": result}})
}

func (h *StorageTypeHandler) FindTips(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindTips(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tips": result}})
}

// AddTip godoc
//
//	@Summary		Add tip to storage type
//	@Description	Add tip to storage type
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			id_type	path		int	true	"Storage type id"
//	@Param			id_tip	path		int	true	"Tip id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storage-types/{id_type}/tips/{id_tip} [post]
//	@Security		Bearer
func (h *StorageTypeHandler) AddTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	result, err := h.svc.CreateTip(ctx.Context(), id, tipID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"tip": result}})
}

// RemoveTip godoc
//
//	@Summary		Remove tip from storage type
//	@Description	Remove tip from storage type
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			id_type	path		int	true	"Storage type id"
//	@Param			id_tip	path		int	true	"Tip id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storage-types/{id_type}/tips/{id_tip} [delete]
//	@Security		Bearer
func (h *StorageTypeHandler) RemoveTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id, tipID); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

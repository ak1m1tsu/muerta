package vaulttype

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/storage-type"
)

type VaultTypeController struct {
	svc service.StorageTypeServicer
	log logger.Logger
}

func New(svc service.StorageTypeServicer, log logger.Logger) *VaultTypeController {
	return &VaultTypeController{
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
func (h *VaultTypeController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateStorageType)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateStorageType(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *VaultTypeController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindStorageTypeByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"type": result}})
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
func (h *VaultTypeController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.StorageTypeFilter)
	if err := utils.ParseFilterAndValidate(ctx, filter); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.FindStorageTypes(ctx.Context(), filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(
		controllers.HTTPSuccess{Success: true, Data: controllers.Data{"types": result, "count": count}},
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
func (h *VaultTypeController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	payload := new(params.UpdateStorageType)
	if err := utils.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateStorageType(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *VaultTypeController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	if err := h.svc.DeleteStorageType(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *VaultTypeController) FindStorages(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindStorages(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"storages": result}})
}

// FindTips godoc
//
//	@Summary		Find tips by storage type id
//	@Description	Find tips by storage type id
//	@Tags			Storage Types
//	@Accept			json
//	@Produce		json
//	@Param			id_type	path		int	true	"Storage type id"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storage-types/{id_type}/tips [get]
func (h *VaultTypeController) FindTips(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	result, err := h.svc.FindTips(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"tips": result}})
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
func (h *VaultTypeController) AddTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	result, err := h.svc.CreateTip(ctx.Context(), id, tipID)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"tip": result}})
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
func (h *VaultTypeController) RemoveTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.TypeID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id, tipID); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

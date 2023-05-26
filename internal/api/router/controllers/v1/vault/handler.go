package vault

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/storage"
)

type VaultController struct {
	svc service.StorageServicer
	log logger.Logger
}

func New(svc service.StorageServicer, log logger.Logger) *VaultController {
	return &VaultController{
		svc: svc,
		log: log,
	}
}

// FindMany godoc
//
//	@Summary		Find storages
//	@Description	Find storages
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.StorageFilter	true	"Filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storages [get]
func (h *VaultController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.StorageFilter)
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
	result, err := h.svc.FindStorages(ctx.Context(), filter)
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
		controllers.HTTPSuccess{
			Success: true,
			Data:    controllers.Data{"storages": result, "count": count},
		},
	)
}

// FindOne godoc
//
//	@Summary		Find storage
//	@Description	Find storage
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			id_storage	path		int	true	"Storage ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/storages/{id_storage} [get]
func (h *VaultController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.FindStorageByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"storage": result}})
}

// Create godoc
//
//	@Summary		Create storage
//	@Description	Create storage
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateStorage	true	"Storage"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/storages [post]
//	@Security		Bearer
func (h *VaultController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateStorage)
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
	if err := h.svc.CreateStorage(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Update godoc
//
//	@Summary		Update storage
//	@Description	Update storage
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			id_storage	path		int					true	"Storage ID"
//	@Param			payload		body		dto.UpdateStorage	true	"Storage"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/storages/{id_storage} [put]
//	@Security		Bearer
func (h *VaultController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	payload := new(params.UpdateStorage)
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
	if err := h.svc.UpdateStorage(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Delete godoc
//
//	@Summary		Delete storage
//	@Description	Delete storage
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			id_storage	path		int	true	"Storage ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/storages/{id_storage} [delete]
//	@Security		Bearer
func (h *VaultController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	if err := h.svc.DeleteStorage(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Restore godoc
//
//	@Summary		Restore storage
//	@Description	Restore storage
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			id_storage	path		int	true	"Storage ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/storages/{id_storage} [patch]
//	@Security		Bearer
func (h *VaultController) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	if err := h.svc.RestoreStorage(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// FindTips godoc
//
//	@Summary		Find tips of storage
//	@Description	Find tips of storage
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			id_storage	path		int	true	"Storage ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/storages/{id_storage}/tips [get]
func (h *VaultController) FindTips(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
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
//	@Summary		Add tip to storage
//	@Description	Add tip
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			id_storage	path		int	true	"Storage ID"
//	@Param			id_tip		path		int	true	"Tip ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/storages/{id_storage}/tips/{id_tip} [post]
//	@Security		Bearer
func (h *VaultController) AddTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
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
//	@Summary		Remove tip from storage
//	@Description	Remove tip
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			id_storage	path		int	true	"Storage ID"
//	@Param			id_tip		path		int	true	"Tip ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/storages/{id_storage}/tips/{id_tip} [delete]
//	@Security		Bearer
func (h *VaultController) RemoveTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id, tipID); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// FindShelfLives godoc
//
//	@Summary		Find shelf lives
//	@Description	Find shelf lives
//	@Tags			Storages
//	@Accept			json
//	@Produce		json
//	@Param			id_storage	path		int	true	"Storage ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/storages/{id_storage}/shelf-lives [get]
func (h *VaultController) FindShelfLives(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.FindShelfLives(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"shelf-lives": result}})
}

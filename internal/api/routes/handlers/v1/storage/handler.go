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

// FindMany godoc
//
// @Summary      Find storages
// @Description  Find storages
// @Tags         Storages
// @Accept       json
// @Produce      json
// @Param        filter query dto.StorageFilter true "Filter"
// @Success      200  {object}  handlers.HTTPSuccess
// @Failure      400  {object}  handlers.HTTPError
// @Failure      500  {object}  handlers.HTTPError
// @Router       /storages [get]
func (h *StorageHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.StorageFilter)
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
		handlers.HTTPSuccess{
			Success: true,
			Data:    handlers.Data{"storages": result, "count": count},
		},
	)
}

// FindOne godoc
//
// @Summary      Find storage
// @Description  Find storage
// @Tags         Storages
// @Accept       json
// @Produce      json
// @Param        id_storage path int true "Storage ID"
// @Success      200  {object}  handlers.HTTPSuccess
// @Failure      400  {object}  handlers.HTTPError
// @Failure      500  {object}  handlers.HTTPError
// @Router       /storages/{id_storage} [get]
func (h *StorageHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.FindStorageByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storage": result}})
}

func (h *StorageHandler) Create(ctx *fiber.Ctx) error {
	var payload *dto.CreateStorage
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
	if err := h.svc.CreateStorage(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Update godoc
//
// @Summary      Update storage
// @Description  Update storage
// @Tags         Storages
// @Accept       json
// @Produce      json
// @Param        id_storage path int true "Storage ID"
// @Param        payload body dto.UpdateStorage true "Storage"
// @Success      200  {object}  handlers.HTTPSuccess
// @Failure      400  {object}  handlers.HTTPError
// @Failure      500  {object}  handlers.HTTPError
// @Router       /storages/{id_storage} [put]
// @Security     Bearer
func (h *StorageHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	var payload *dto.UpdateStorage
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
	if err := h.svc.UpdateStorage(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Delete godoc
//
// @Summary      Delete storage
// @Description  Delete storage
// @Tags         Storages
// @Accept       json
// @Produce      json
// @Param        id_storage path int true "Storage ID"
// @Success      200  {object}  handlers.HTTPSuccess
// @Failure      400  {object}  handlers.HTTPError
// @Failure      500  {object}  handlers.HTTPError
// @Router       /storages/{id_storage} [delete]
// @Security     Bearer
func (h *StorageHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	if err := h.svc.DeleteStorage(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}
// Restore godoc
//
// @Summary      Restore storage
// @Description  Restore storage
// @Tags         Storages
// @Accept       json
// @Produce      json
// @Param        id_storage path int true "Storage ID"
// @Success      200  {object}  handlers.HTTPSuccess
// @Failure      400  {object}  handlers.HTTPError
// @Failure      500  {object}  handlers.HTTPError
// @Router       /storages/{id_storage} [patch]
// @Security     Bearer
func (h *StorageHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	if err := h.svc.RestoreStorage(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *StorageHandler) FindTips(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
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
// @Summary      Add tip to storage
// @Description  Add tip
// @Tags         Storages
// @Accept       json
// @Produce      json
// @Param        id_storage path int true "Storage ID"
// @Param        id_tip path int true "Tip ID"
// @Success      200  {object}  handlers.HTTPSuccess
// @Failure      400  {object}  handlers.HTTPError
// @Failure      500  {object}  handlers.HTTPError
// @Router       /storages/{id_storage}/tips/{id_tip} [post]
// @Security     Bearer
func (h *StorageHandler) AddTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
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
// @Summary      Remove tip from storage
// @Description  Remove tip
// @Tags         Storages
// @Accept       json
// @Produce      json
// @Param        id_storage path int true "Storage ID"
// @Param        id_tip path int true "Tip ID"
// @Success      200  {object}  handlers.HTTPSuccess
// @Failure      400  {object}  handlers.HTTPError
// @Failure      500  {object}  handlers.HTTPError
// @Router       /storages/{id_storage}/tips/{id_tip} [delete]
// @Security     Bearer
func (h *StorageHandler) RemoveTip(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StorageID).(int)
	tipID := ctx.Locals(context.TipID).(int)
	if err := h.svc.DeleteTip(ctx.Context(), id, tipID); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindShelfLives godoc
//
// @Summary      Find shelf lives
// @Description  Find shelf lives
// @Tags         Storages
// @Accept       json
// @Produce      json
// @Param        id_storage path int true "Storage ID"
// @Success      200  {object}  handlers.HTTPSuccess
// @Failure      400  {object}  handlers.HTTPError
// @Failure      500  {object}  handlers.HTTPError
// @Router       /storages/{id_storage}/shelf-lives [get]
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

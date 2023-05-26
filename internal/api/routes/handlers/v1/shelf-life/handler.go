package shelflife

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life"
)

type ShelfLifeHandler struct {
	svc service.ShelfLifeServicer
	log logger.Logger
}

func New(svc service.ShelfLifeServicer, log logger.Logger) ShelfLifeHandler {
	return ShelfLifeHandler{
		svc: svc,
		log: log,
	}
}

// Create godoc
//
//	@Summary		Create shelf life
//	@Description	Create shelf life
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateShelfLife	true	"Shelf Life"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/shelf-lives [post]
//	@Security		Bearer
func (h *ShelfLifeHandler) Create(ctx *fiber.Ctx) error {
	payload := new(dto.CreateShelfLife)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateShelfLife(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindOne godoc
//
//	@Summary		Find shelf life by id
//	@Description	Find shelf life by id
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			shelf_life_id	path		int	true	"Shelf Life ID"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		404				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/shelf-lives/{shelf_life_id} [get]
func (h *ShelfLifeHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	result, err := h.svc.FindShelfLifeByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf_life": result}})
}

// FindMany godoc
//
//	@Summary		Find many shelf lifes
//	@Description	Find many shelf lifes
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.ShelfLifeFilter	true	"Shelf Life Filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/shelf-lives [get]
func (h *ShelfLifeHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.ShelfLifeFilter)
	if err := common.ParseFilterAndValidate(ctx, filter); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.FindShelfLifes(ctx.Context(), filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(
		handlers.HTTPSuccess{
			Success: true,
			Data:    handlers.Data{"shelf_lives": result, "count": count},
		},
	)
}

// Update godoc
//
//	@Summary		Update shelf life
//	@Description	Update shelf life
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			shelf_life_id	path		int					true	"Shelf Life ID"
//	@Param			payload			body		dto.UpdateShelfLife	true	"Shelf Life"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/shelf-lives/{shelf_life_id} [put]
//	@Security		Bearer
func (h *ShelfLifeHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	payload := new(dto.UpdateShelfLife)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateShelfLife(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Delete godoc
//
//	@Summary		Delete shelf life
//	@Description	Delete shelf life
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			shelf_life_id	path		int	true	"Shelf Life ID"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/shelf-lives/{shelf_life_id} [delete]
//	@Security		Bearer
func (h *ShelfLifeHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	if err := h.svc.DeleteShelfLife(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Restore godoc
//
//	@Summary		Restore shelf life
//	@Description	Restore shelf life
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			shelf_life_id	path		int	true	"Shelf Life ID"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/shelf-lives/{shelf_life_id} [patch]
//	@Security		Bearer
func (h *ShelfLifeHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	if err := h.svc.RestoreShelfLife(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindStatuses godoc
//
//	@Summary		Find shelf life statuses
//	@Description	Find shelf life statuses
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			shelf_life_id	path		int	true	"Shelf Life ID"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/shelf-lives/{shelf_life_id}/statuses [get]
func (h *ShelfLifeHandler) FindStatuses(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	result, err := h.svc.FindShelfLifeStatuses(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"statuses": result}})
}

// AddStatus godoc
//
//	@Summary		Add shelf life status
//	@Description	Add shelf life status
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			shelf_life_id	path		int	true	"Shelf Life ID"
//	@Param			status_id		path		int	true	"Status ID"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/shelf-lives/{shelf_life_id}/statuses/{status_id} [post]
//	@Security		Bearer
func (h *ShelfLifeHandler) AddStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	statusID := ctx.Locals(context.StatusID).(int)
	result, err := h.svc.CreateShelfLifeStatus(ctx.Context(), id, statusID)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"status": result}})
}

// RemoveStatus godoc
//
//	@Summary		Remove shelf life status
//	@Description	Remove shelf life status
//	@Tags			Shelf Lives
//	@Accept			json
//	@Produce		json
//	@Param			shelf_life_id	path		int	true	"Shelf Life ID"
//	@Param			status_id		path		int	true	"Status ID"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/shelf-lives/{shelf_life_id}/statuses/{status_id} [delete]
//	@Security		Bearer
func (h *ShelfLifeHandler) RemoveStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	statusID := ctx.Locals(context.StatusID).(int)
	if err := h.svc.DeleteShelfLifeStatus(ctx.Context(), id, statusID); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

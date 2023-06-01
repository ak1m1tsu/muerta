package shelflife

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life"
)

type ShelfLifeController struct {
	svc service.ShelfLifeServicer
	log logger.Logger
}

func New(svc service.ShelfLifeServicer, log logger.Logger) ShelfLifeController {
	return ShelfLifeController{
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
func (h *ShelfLifeController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateShelfLife)
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
	if err := h.svc.CreateShelfLife(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *ShelfLifeController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	result, err := h.svc.FindShelfLifeByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"shelf_life": result}})
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
func (h *ShelfLifeController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.ShelfLifeFilter)
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
	result, err := h.svc.FindShelfLifes(ctx.Context(), filter)
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
			Data:    controllers.Data{"shelf_lives": result, "count": count},
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
func (h *ShelfLifeController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	payload := new(params.UpdateShelfLife)
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
	if err := h.svc.UpdateShelfLife(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *ShelfLifeController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	if err := h.svc.DeleteShelfLife(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *ShelfLifeController) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	if err := h.svc.RestoreShelfLife(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *ShelfLifeController) FindStatuses(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	result, err := h.svc.FindShelfLifeStatuses(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"statuses": result}})
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
func (h *ShelfLifeController) AddStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	statusID := ctx.Locals(context.StatusID).(int)
	result, err := h.svc.CreateShelfLifeStatus(ctx.Context(), id, statusID)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"status": result}})
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
func (h *ShelfLifeController) RemoveStatus(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.ShelfLifeID).(int)
	statusID := ctx.Locals(context.StatusID).(int)
	if err := h.svc.DeleteShelfLifeStatus(ctx.Context(), id, statusID); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

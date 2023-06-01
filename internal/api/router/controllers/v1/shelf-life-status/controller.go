package shelflifestatus

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life-status"
)

type ShelfLifeStatusController struct {
	svc service.ShelfLifeStatusServicer
	log logger.Logger
}

func New(svc service.ShelfLifeStatusServicer, log logger.Logger) ShelfLifeStatusController {
	return ShelfLifeStatusController{
		svc: svc,
		log: log,
	}
}

// Create godoc
//
//	@Summary		Create shelf life status
//	@Description	Create shelf life status
//	@Tags			Shelf Life Status
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateShelfLifeStatus	true	"Shelf Life Status"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/shelf-life-statuses [post]
//	@Security		Bearer
func (h *ShelfLifeStatusController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateShelfLifeStatus)
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
	if err := h.svc.CreateShelfLifeStatus(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// FindOne godoc
//
//	@Summary		Find one shelf life status
//	@Description	Find one shelf life status
//	@Tags			Shelf Life Status
//	@Accept			json
//	@Produce		json
//
// /	@Param			id_status path int	true	"Shelf Life Status ID"
//
//	@Success		200	{object}	handlers.HTTPSuccess
//	@Failure		400	{object}	handlers.HTTPError
//	@Failure		500	{object}	handlers.HTTPError
//	@Router			/shelf-life-statuses/{id_status} [get]
func (h *ShelfLifeStatusController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	result, err := h.svc.FindShelfLifeStatusByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"status": result}})
}

// FindMnay godoc
//
//	@Summary		Find many shelf life status
//	@Description	Find many shelf life status
//	@Tags			Shelf Life Status
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.ShelfLifeStatusFilter	false	"Filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/shelf-life-statuses [get]
func (h *ShelfLifeStatusController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.ShelfLifeStatusFilter)
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
	result, err := h.svc.FindShelfLifeStatuss(ctx.Context(), filter)
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
		controllers.HTTPSuccess{Success: true, Data: controllers.Data{"statues": result, "count": count}},
	)
}

// Update godoc
//
//	@Summary		Update shelf life status
//	@Description	Update shelf life status
//	@Tags			Shelf Life Status
//	@Accept			json
//	@Produce		json
//	@Param			id_status	path		int							true	"Shelf life status ID"
//	@Param			body		body		dto.UpdateShelfLifeStatus	true	"Shelf life status"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/shelf-life-statuses/{id_status} [put]
//	@Security		Bearer
func (h *ShelfLifeStatusController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	payload := new(params.UpdateShelfLifeStatus)
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
	if err := h.svc.UpdateShelfLifeStatus(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Delete godoc
//
//	@Summary		Delete shelf life status
//	@Description	Delete shelf life status
//	@Tags			Shelf Life Status
//	@Accept			json
//	@Produce		json
//	@Param			id_status	path		int	true	"Shelf life status ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/shelf-life-statuses/{id_status} [delete]
//	@Security		Bearer
func (h *ShelfLifeStatusController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	if err := h.svc.DeleteShelfLifeStatus(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

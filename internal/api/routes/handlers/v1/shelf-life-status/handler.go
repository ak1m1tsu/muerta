package shelflifestatus

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life-status"
)

type ShelfLifeStatusHandler struct {
	svc service.ShelfLifeStatusServicer
	log *log.Logger
}

func New(svc service.ShelfLifeStatusServicer, log *log.Logger) ShelfLifeStatusHandler {
	return ShelfLifeStatusHandler{
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
func (h *ShelfLifeStatusHandler) Create(ctx *fiber.Ctx) error {
	payload := new(dto.CreateShelfLifeStatus)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateShelfLifeStatus(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *ShelfLifeStatusHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	result, err := h.svc.FindShelfLifeStatusByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"status": result}})
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
func (h *ShelfLifeStatusHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.ShelfLifeStatusFilter)
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
	result, err := h.svc.FindShelfLifeStatuss(ctx.Context(), filter)
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
		handlers.HTTPSuccess{Success: true, Data: handlers.Data{"statues": result, "count": count}},
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
func (h *ShelfLifeStatusHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	payload := new(dto.UpdateShelfLifeStatus)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.ValidationError(ctx, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateShelfLifeStatus(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *ShelfLifeStatusHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StatusID).(int)
	if err := h.svc.DeleteShelfLifeStatus(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

package measure

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/measure"
)

type MeasureHandler struct {
	svc service.MeasureServicer
	log *log.Logger
}

func New(svc service.MeasureServicer, log *log.Logger) MeasureHandler {
	return MeasureHandler{
		svc: svc,
		log: log,
	}
}

// CreateMeasure creates a new measure record
//
//	@Summary		Create a new measure record
//	@Description	Create a new measure record based on the given payload.
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateMeasure	true	"Payload of the measure record to create"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/measures [post]
func (h *MeasureHandler) CreateMeasure(ctx *fiber.Ctx) error {
	var payload *dto.CreateMeasure
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
	if err := h.svc.CreateMeasure(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// FindMeasureByID finds a measure by ID
//
//	@Summary		Find a measure by ID
//	@Description	Find a measure by ID
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Measure ID"
//	@Success		200	{object}	handlers.HTTPSuccess
//	@Failure		404	{object}	handlers.HTTPError
//	@Failure		502	{object}	handlers.HTTPError
//	@Router			/measures/{id} [get]
func (h *MeasureHandler) FindMeasureByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	dto, err := h.svc.FindMeasureByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusNotFound).
			JSON(handlers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data:    handlers.Data{"measure": dto},
	})
}

// FindMeasures returns a list of measures and their count based on the given filter.
//
//	@Summary		Find measures
//	@Description	Returns a list of measures and their count based on the given filter.
//	@ID				find-measures
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.MeasureFilter	false	"Filter for measures"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/measures [get]
func (h *MeasureHandler) FindMeasures(ctx *fiber.Ctx) error {
	filter := new(dto.MeasureFilter)
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
	result, err := h.svc.FindMeasures(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{
		Success: true,
		Data:    handlers.Data{"measures": result, "count": count},
	})
}

// UpdateMeasure updates a measure with the given ID using the provided payload
//
//	@Summary		Update a measure
//	@Description	Update a measure with the given ID
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Measure ID"
//	@Param			payload	body		dto.UpdateMeasure	true	"Measure payload"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/measures/{id} [put]
func (h *MeasureHandler) UpdateMeasure(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	payload := new(dto.UpdateMeasure)
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
	if err := h.svc.UpdateMeasure(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// DeleteMeasure deletes a measure by ID.
//
//	@Summary		Delete a measure
//	@Description	Deletes a measure by ID.
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Measure ID"
//	@Success		200	{object}	handlers.HTTPSuccess
//	@Failure		400	{object}	handlers.HTTPError
//	@Failure		502	{object}	handlers.HTTPError
//	@Router			/measures/{id} [delete]
func (h *MeasureHandler) DeleteMeasure(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	if err := h.svc.DeleteMeasure(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

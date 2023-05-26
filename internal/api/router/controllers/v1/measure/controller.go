package measure

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/measure"
)

type MeasureController struct {
	svc service.MeasureServicer
	log logger.Logger
}

func New(svc service.MeasureServicer, log logger.Logger) MeasureController {
	return MeasureController{
		svc: svc,
		log: log,
	}
}

// Create creates a new measure record
//
//	@Summary		Create a new measure record
//	@Description	Creates a new measure record based on the given payload.
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateMeasure	true	"Payload of the measure record to create"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		502		{object}	handlers.HTTPError
//	@Router			/measures [post]
//	@Security		Bearer
func (h *MeasureController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateMeasure)
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
	if err := h.svc.CreateMeasure(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{
		Success: true,
	})
}

// FindOne finds a measure by ID
//
//	@Summary		Find a measure by ID
//	@Description	Find a measure by ID
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			measure_id	path		int	true	"Measure ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/measures/{measure_id} [get]
func (h *MeasureController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	dto, err := h.svc.FindMeasureByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusNotFound).
			JSON(controllers.HTTPError{Error: fiber.ErrNotFound.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{
		Success: true,
		Data:    controllers.Data{"measure": dto},
	})
}

// FindMany returns a list of measures and their count based on the given filter.
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
func (h *MeasureController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.MeasureFilter)
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
	result, err := h.svc.FindMeasures(ctx.Context(), filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(controllers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{
		Success: true,
		Data:    controllers.Data{"measures": result, "count": count},
	})
}

// Update updates a measure with the given ID using the provided payload
//
//	@Summary		Update a measure
//	@Description	Updates a measure with the given ID
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			measure_id	path		int					true	"Measure ID"
//	@Param			payload		body		dto.UpdateMeasure	true	"Measure payload"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/measures/{measure_id} [put]
//	@Security		Bearer
func (h *MeasureController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	payload := new(params.UpdateMeasure)
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
	if err := h.svc.UpdateMeasure(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Delete deletes a measure by ID.
//
//	@Summary		Delete a measure
//	@Description	Deletes a measure by ID.
//	@Tags			Measures
//	@Accept			json
//	@Produce		json
//	@Param			measure_id	path		int	true	"Measure ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		502			{object}	handlers.HTTPError
//	@Router			/measures/{measure_id} [delete]
//	@Security		Bearer
func (h *MeasureController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.MeasureID).(int)
	if err := h.svc.DeleteMeasure(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

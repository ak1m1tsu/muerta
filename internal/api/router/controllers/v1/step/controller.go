package step

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/step"
)

type StepController struct {
	svc service.StepServicer
	log logger.Logger
}

func New(svc service.StepServicer, log logger.Logger) *StepController {
	return &StepController{
		svc: svc,
		log: log,
	}
}

// FindMany godoc
//
//	@Summary		Find many steps
//	@Description	Find many steps
//	@Tags			Steps
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.StepFilter	true	"Filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/steps [get]
func (h *StepController) FinaMany(ctx *fiber.Ctx) error {
	filter := new(params.StepFilter)
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
	result, err := h.svc.FindSteps(ctx.Context(), filter)
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
		controllers.HTTPSuccess{Success: true, Data: controllers.Data{"steps": result, "count": count}},
	)
}

// Create godoc
//
//	@Summary		Create a step
//	@Description	Create a step
//	@Tags			Steps
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateStep	true	"CreateStep"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/steps [post]
//	@Security		Bearer
func (h *StepController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateStep)
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
	result, err := h.svc.CreateStep(ctx.Context(), payload)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"step": result}})
}

// FindOne godoc
//
//	@Summary		Find one step
//	@Description	Find one step
//	@Tags			Steps
//	@Accept			json
//	@Produce		json
//	@Param			id_step	path		int	true	"Step ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/steps/{id_step} [get]
func (h *StepController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	result, err := h.svc.FindStep(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"step": result}})
}

// Update godoc
//
//	@Summary		Update a step
//	@Description	Update a step
//	@Tags			Steps
//	@Accept			json
//	@Produce		json
//	@Param			id_step	path		int				true	"Step ID"
//	@Param			payload	body		dto.UpdateStep	true	"UpdateStep"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/steps/{id_step} [put]
//	@Security		Bearer
func (h *StepController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	payload := new(params.UpdateStep)
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
	result, err := h.svc.UpdateStep(ctx.Context(), id, payload)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"step": result}})
}

// Delete godoc
//
//	@Summary		Delete a step
//	@Description	Delete a step
//	@Tags			Steps
//	@Accept			json
//	@Produce		json
//	@Param			id_step	path		int	true	"Step ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/steps/{id_step} [delete]
//	@Security		Bearer
func (h *StepController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	if err := h.svc.DeleteStep(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Restore godoc
//
//	@Summary		Restore a step
//	@Description	Restore a step
//	@Tags			Steps
//	@Accept			json
//	@Produce		json
//	@Param			id_step	path		int	true	"Step ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/steps/{id_step} [patch]
//	@Security		Bearer
func (h *StepController) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	result, err := h.svc.RestoreStep(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"step": result}})
}

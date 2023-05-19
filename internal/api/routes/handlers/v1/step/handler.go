package step

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/step"
)

type StepHandler struct {
	svc service.StepServicer
	log *log.Logger
}

func New(svc service.StepServicer, log *log.Logger) *StepHandler {
	return &StepHandler{
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
func (h *StepHandler) FinaMany(ctx *fiber.Ctx) error {
	filter := new(dto.StepFilter)
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
	result, err := h.svc.FindSteps(ctx.Context(), filter)
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
		handlers.HTTPSuccess{Success: true, Data: handlers.Data{"steps": result, "count": count}},
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
func (h *StepHandler) Create(ctx *fiber.Ctx) error {
	payload := new(dto.CreateStep)
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
	result, err := h.svc.CreateStep(ctx.Context(), payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"step": result}})
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
func (h *StepHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	result, err := h.svc.FindStep(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"step": result}})
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
func (h *StepHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	payload := new(dto.UpdateStep)
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
	result, err := h.svc.UpdateStep(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"step": result}})
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
func (h *StepHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	if err := h.svc.DeleteStep(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *StepHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.StepID).(int)
	result, err := h.svc.RestoreStep(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"step": result}})
}

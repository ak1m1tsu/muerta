package role

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/role"
)

type RoleController struct {
	svc service.RoleServicer
	log logger.Logger
}

func New(svc service.RoleServicer, log logger.Logger) *RoleController {
	return &RoleController{
		svc: svc,
		log: log,
	}
}

// FindMany godoc
//
//	@Summary		Find many roles
//	@Description	Find many roles
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.RoleFilter	true	"Filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/roles [get]
func (h *RoleController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.RoleFilter)
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
	result, err := h.svc.FindRoles(ctx.Context(), filter)
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
		controllers.HTTPSuccess{Success: true, Data: controllers.Data{"roles": result, "count": count}},
	)
}

// FindOne godoc
//
//	@Summary		Find one role
//	@Description	Find one role
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			role_id	path		int	true	"Role ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/roles/{role_id} [get]
func (h *RoleController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RoleID).(int)
	result, err := h.svc.FindRoleByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"roles": result}})
}

// FindOne godoc
//
//	@Summary		Create role
//	@Description	Create role
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			role	body		dto.CreateRole	true	"Role"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/roles [post]
//	@Security		Bearer
func (h *RoleController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateRole)
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
	if err := h.svc.CreateRole(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Update godoc
//
//	@Summary		Update role
//	@Description	Update role
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			role_id	path		int				true	"Role ID"
//	@Param			role	body		dto.UpdateRole	true	"Role"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/roles/{role_id} [put]
//	@Security		Bearer
func (h *RoleController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RoleID).(int)
	payload := new(params.UpdateRole)
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
	if err := h.svc.UpdateRole(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Delete godoc
//
//	@Summary		Delete role
//	@Description	Delete role
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			role_id	path		int	true	"Role ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/roles/{role_id} [delete]
//	@Security		Bearer
func (h *RoleController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RoleID).(int)
	if err := h.svc.DeleteRole(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Restore godoc
//
//	@Summary		Restore role
//	@Description	Restore role
//	@Tags			Roles
//	@Accept			json
//	@Produce		json
//	@Param			role_id	path		int	true	"Role ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/roles/{role_id} [patch]
//	@Security		Bearer
func (h *RoleController) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RoleID).(int)
	if err := h.svc.RestoreRole(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

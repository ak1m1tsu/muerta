package role

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/role"
)

type RoleHandler struct {
	svc service.RoleServicer
	log *log.Logger
}

func New(svc service.RoleServicer, log *log.Logger) *RoleHandler {
	return &RoleHandler{
		svc: svc,
		log: log,
	}
}

func (h *RoleHandler) FindRoles(ctx *fiber.Ctx) error {
	filter := new(dto.RoleFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindRoles(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"roles": result, "count": count},
	))
}

func (h *RoleHandler) FindRole(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RoleID).(int)
	result, err := h.svc.FindRoleByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"roles": result},
	))
}

func (h *RoleHandler) CreateRole(ctx *fiber.Ctx) error {
	var payload *dto.CreateRoleDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateRole(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *RoleHandler) UpdateRole(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RoleID).(int)
	var payload *dto.UpdateRoleDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateRole(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *RoleHandler) DeleteRole(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RoleID).(int)
	if err := h.svc.DeleteRole(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *RoleHandler) RestoreRole(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.RoleID).(int)
	if err := h.svc.RestoreRole(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.JSON(handlers.SuccessResponse())
}

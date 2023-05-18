package usersetting

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	usersetting "github.com/romankravchuk/muerta/internal/services/user-setting"
)

type UserSettingHandler struct {
	svc usersetting.UserSettingsServicer
	log *log.Logger
}

func New(svc usersetting.UserSettingsServicer, log *log.Logger) *UserSettingHandler {
	return &UserSettingHandler{
		svc: svc,
		log: log,
	}
}

func (h *UserSettingHandler) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	result, err := h.svc.FindSettingByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"setting": result}})
}

func (h *UserSettingHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.SettingFilter)
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
	result, err := h.svc.FindSettings(ctx.Context(), filter)
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
		handlers.HTTPSuccess{
			Success: true,
			Data:    handlers.Data{"settings": result, "count": count},
		},
	)
}

func (h *UserSettingHandler) Create(ctx *fiber.Ctx) error {
	var payload *dto.CreateSetting
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
	if err := h.svc.CreateSetting(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *UserSettingHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	payload := new(dto.UpdateSetting)
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
	if err := h.svc.UpdateSetting(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *UserSettingHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	if err := h.svc.DeleteSetting(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *UserSettingHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	result, err := h.svc.RestoreSetting(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"setting": result}})
}

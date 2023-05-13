package usersetting

import (
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
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"setting": result},
	))
}

func (h *UserSettingHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.SettingFilterDTO)
	if err := common.GetFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	result, err := h.svc.FindSettings(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	count, err := h.svc.Count(ctx.Context())
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse().WithData(
		handlers.Data{"settings": result, "count": count},
	))
}

func (h *UserSettingHandler) Create(ctx *fiber.Ctx) error {
	var payload *dto.CreateSettingDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.CreateSetting(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *UserSettingHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	payload := new(dto.UpdateSettingDTO)
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if err := h.svc.UpdateSetting(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *UserSettingHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	if err := h.svc.DeleteSetting(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

func (h *UserSettingHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	if err := h.svc.RestoreSetting(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(handlers.SuccessResponse())
}

package usersetting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	usersettings "github.com/romankravchuk/muerta/internal/services/user-settings"
)

type UserSettingHandler struct {
	svc usersettings.UserSettingsServicer
	log *log.Logger
}

func New(svc usersettings.UserSettingsServicer, log *log.Logger) *UserSettingHandler {
	return &UserSettingHandler{
		svc: svc,
		log: log,
	}
}

func (h *UserSettingHandler) FindByID(ctx *fiber.Ctx) error {
	id, err := common.GetIdByFiberCtx(ctx)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	setting, err := h.svc.FindSettingByID(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"setting": setting},
	})
}

func (h *UserSettingHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.SettingFilterDTO)
	if err := common.GetSettingFilterByFiberCtx(ctx, filter); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	settings, err := h.svc.FindSettings(ctx.Context(), filter)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"settings": settings},
	})
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
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

// func (h *UserSettingHandler) Update(ctx *fiber.Ctx) error {}
// func (h *UserSettingHandler) Delete(ctx *fiber.Ctx) error {}

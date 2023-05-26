package usersetting

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	usersetting "github.com/romankravchuk/muerta/internal/services/user-setting"
)

type UserSettingHandler struct {
	svc usersetting.UserSettingsServicer
	log logger.Logger
}

func New(svc usersetting.UserSettingsServicer, log logger.Logger) *UserSettingHandler {
	return &UserSettingHandler{
		svc: svc,
		log: log,
	}
}

// FindOne godoc
//
//	@Summary		Find setting
//	@Description	Find setting
//	@Tags			Settings
//	@Accept			json
//	@Produce		json
//	@Param			id_setting	path		integer	true	"Setting ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/settings/{id_setting} [get]
func (h *UserSettingHandler) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	result, err := h.svc.FindSettingByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"setting": result}})
}

// FindMany godoc
//
//	@Summary		Find settings
//	@Description	Find settings
//	@Tags			Settings
//	@Accept			json
//	@Produce		json
//	@Param			filter	query		dto.SettingFilter	true	"Filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/settings [get]
func (h *UserSettingHandler) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.SettingFilter)
	if err := common.ParseFilterAndValidate(ctx, filter); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	result, err := h.svc.FindSettings(ctx.Context(), filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	count, err := h.svc.Count(ctx.Context(), *filter)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
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

// Create godoc
//
//	@Summary		Create setting
//	@Description	Create setting
//	@Tags			Settings
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateSetting	true	"Setting"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/settings [post]
//	@Security		Bearer
func (h *UserSettingHandler) Create(ctx *fiber.Ctx) error {
	payload := new(dto.CreateSetting)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.CreateSetting(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Update godoc
//
//	@Summary		Update setting
//	@Description	Update setting
//	@Tags			Settings
//	@Accept			json
//	@Produce		json
//	@Param			id_setting	path		int					true	"Setting ID"
//	@Param			payload		body		dto.UpdateSetting	true	"Setting"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/settings/{id_setting} [put]
//	@Security		Bearer
func (h *UserSettingHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	payload := new(dto.UpdateSetting)
	if err := common.ParseBodyAndValidate(ctx, payload); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			h.log.Error(ctx, logger.Validation, err)
			return ctx.Status(http.StatusBadRequest).
				JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
		}
		h.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if err := h.svc.UpdateSetting(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *UserSettingHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	if err := h.svc.DeleteSetting(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

// Restore godoc
//
//	@Summary		Restore setting
//	@Description	Restore setting
//	@Tags			Settings
//	@Accept			json
//	@Produce		json
//	@Param			id_setting	path		int	true	"Setting ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/settings/{id_setting} [patch]
//	@Security		Bearer
func (h *UserSettingHandler) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	result, err := h.svc.RestoreSetting(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"setting": result}})
}

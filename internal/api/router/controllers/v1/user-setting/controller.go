package usersetting

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	usersetting "github.com/romankravchuk/muerta/internal/services/user-setting"
)

type UserSettingController struct {
	svc usersetting.UserSettingsServicer
	log logger.Logger
}

func New(svc usersetting.UserSettingsServicer, log logger.Logger) *UserSettingController {
	return &UserSettingController{
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
func (h *UserSettingController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	result, err := h.svc.FindSettingByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"setting": result}})
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
func (h *UserSettingController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.SettingFilter)
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
	result, err := h.svc.FindSettings(ctx.Context(), filter)
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
		controllers.HTTPSuccess{
			Success: true,
			Data:    controllers.Data{"settings": result, "count": count},
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
func (h *UserSettingController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateSetting)
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
	if err := h.svc.CreateSetting(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *UserSettingController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	payload := new(params.UpdateSetting)
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
	if err := h.svc.UpdateSetting(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

func (h *UserSettingController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	if err := h.svc.DeleteSetting(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
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
func (h *UserSettingController) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.SettingID).(int)
	result, err := h.svc.RestoreSetting(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"setting": result}})
}

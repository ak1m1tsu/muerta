package user

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/controllers"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/api/router/utils"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/user"
)

type UserController struct {
	svc service.UserServicer
	log logger.Logger
}

func New(svc service.UserServicer, log logger.Logger) *UserController {
	return &UserController{svc: svc, log: log}
}

// FindOne godoc
//
//	@Summary		Find user by id
//	@Description	Find user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user	path		int	true	"User ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users/{id_user} [get]
func (h *UserController) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	user, err := h.svc.FindUserByID(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Client, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"user": user}})
}

// FindMany godoc
//
//	@Summary		Find users
//	@Description	Find users
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.UserFilter	true	"Filter"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users [get]
func (h *UserController) FindMany(ctx *fiber.Ctx) error {
	filter := new(params.UserFilter)
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
	result, err := h.svc.FindUsers(ctx.Context(), filter)
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
		controllers.HTTPSuccess{Success: true, Data: controllers.Data{"users": result, "count": count}},
	)
}

// Create godoc
//
//	@Summary		Create user
//	@Description	Create user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateUser	true	"User"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users [post]
//	@Security		Bearer
func (h *UserController) Create(ctx *fiber.Ctx) error {
	payload := new(params.CreateUser)
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
	if err := h.svc.CreateUser(ctx.Context(), payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Update godoc
//
//	@Summary		Update user
//	@Description	Update user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.UpdateUser	true	"User"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users [put]
//	@Security		Bearer
func (h *UserController) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	payload := new(params.UpdateUser)
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
	if err := h.svc.UpdateUser(ctx.Context(), id, payload); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// Delete godoc
//
//	@Summary		Delete user
//	@Description	Delete user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handlers.HTTPSuccess
//	@Failure		400	{object}	handlers.HTTPError
//	@Failure		404	{object}	handlers.HTTPError
//	@Failure		500	{object}	handlers.HTTPError
//	@Router			/users [delete]
//	@Security		Bearer
func (h *UserController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	if err := h.svc.DeleteUser(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

func (h *UserController) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	if err := h.svc.RestoreUser(ctx.Context(), id); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(controllers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// FindSettings godoc
//
//	@Summary		Find user settings
//	@Description	Find user settings
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user	query		int	true	"User ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users/{id_user}/settings [get]
//	@Security		Bearer
func (h *UserController) FindSettings(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	result, err := h.svc.FindSettings(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"settings": result}})
}

// UpdateSettings godoc
//
//	@Summary		Update user settings
//	@Description	Update user settings
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user		query		int						true	"User ID"
//	@Param			id_setting	query		int						true	"User ID"
//	@Param			payload		body		dto.UpdateUserSetting	true	"User"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/users/{id_user}/settings/{id_setting} [put]
//	@Security		Bearer
func (h *UserController) UpdateSetting(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	payload := new(params.UpdateUserSetting)
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
	result, err := h.svc.UpdateSetting(ctx.Context(), id, payload)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"settings": result}})
}

// FindRoles godoc
//
//	@Summary		Find user roles
//	@Description	Find user roles
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user	query		int	true	"User ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users/{id_user}/roles [get]
func (h *UserController) FindRoles(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	result, err := h.svc.FindRoles(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"roles": result}})
}

// FindStorages godoc
//
//	@Summary		Find user storages
//	@Description	Find user storages
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user	query		int	true	"User ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users/{id_user}/storages [get]
func (h *UserController) FindStorages(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	result, err := h.svc.FindStorages(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"storages": result}})
}

// AddStorage godoc
//
//	@Summary		Add user storage
//	@Description	Add user storage
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user		query		int	true	"User ID"
//	@Param			id_storage	query		int	true	"Storage ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/users/{id_user}/storages/{id_storage} [post]
//	@Security		Bearer
func (h *UserController) AddStorage(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	storageID := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.AddStorage(ctx.Context(), id, storageID)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"storages": result}})
}

// RemoveStorage godoc
//
//	@Summary		Remove user storage
//	@Description	Remove user storage
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user		query		int	true	"User ID"
//	@Param			id_storage	query		int	true	"Storage ID"
//	@Success		200			{object}	handlers.HTTPSuccess
//	@Failure		400			{object}	handlers.HTTPError
//	@Failure		404			{object}	handlers.HTTPError
//	@Failure		500			{object}	handlers.HTTPError
//	@Router			/users/{id_user}/storages/{id_storage} [delete]
//	@Security		Bearer
func (h *UserController) RemoveStorage(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	storageID := ctx.Locals(context.StorageID).(int)
	err := h.svc.RemoveStorage(ctx.Context(), id, storageID)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

// FindShelfLives godoc
//
//	@Summary		Find user shelf lives
//	@Description	Find user shelf lives
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user	query		int	true	"User ID"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users/{id_user}/shelf-lives [get]
func (h *UserController) FindShelfLives(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	result, err := h.svc.FindShelfLives(ctx.Context(), id)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"shelf-lives": result}})
}

// CreateShelfLife godoc
//
//	@Summary		Create user shelf life
//	@Description	Create user shelf life
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user	query		int					true	"User ID"
//	@Param			payload	body		dto.CreateShelfLife	true	"User"
//	@Success		200		{object}	handlers.HTTPSuccess
//	@Failure		400		{object}	handlers.HTTPError
//	@Failure		404		{object}	handlers.HTTPError
//	@Failure		500		{object}	handlers.HTTPError
//	@Router			/users/{id_user}/shelf-lives [post]
//	@Security		Bearer
func (h *UserController) CreateShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	payload := new(params.CreateShelfLife)
	payload.UserID = id
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
	result, err := h.svc.CreateShelfLife(ctx.Context(), id, payload)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"shelf-life": result}})
}

// UpdateShelfLife godoc
//
//	@Summary		Update user shelf life
//	@Description	Update user shelf life
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user			query		int					true	"User ID"
//	@Param			id_shelf_life	query		int					true	"User ID"
//	@Param			payload			body		dto.UserShelfLife	true	"User"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		404				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/users/{id_user}/shelf-lives/{id_shelf_life} [put]
//	@Security		Bearer
func (h *UserController) UpdateShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	shelfLifeID := ctx.Locals(context.ShelfLifeID).(int)
	payload := new(params.UserShelfLife)
	payload.ShelfLifeID = shelfLifeID
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
	result, err := h.svc.UpdateShelfLife(ctx.Context(), id, payload)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"shelf-life": result}})
}

// RestoreShelfLife godoc
//
//	@Summary		Restore user shelf life
//	@Description	Restore user shelf life
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user			query		int	true	"User ID"
//	@Param			id_shelf_life	query		int	true	"Shelf Life ID"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		404				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/users/{id_user}/shelf-lives/{id_shelf_life} [patch]
//	@Security		Bearer
func (h *UserController) RestoreShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	shelfLifeID := ctx.Locals(context.ShelfLifeID).(int)
	result, err := h.svc.RestoreShelfLife(ctx.Context(), id, shelfLifeID)
	if err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true, Data: controllers.Data{"shelf-life": result}})
}

// DeleteShelfLife godoc
//
//	@Summary		Delete user shelf life
//	@Description	Delete user shelf life
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id_user			query		int	true	"User ID"
//	@Param			id_shelf_life	query		int	true	"Shelf Life ID"
//	@Success		200				{object}	handlers.HTTPSuccess
//	@Failure		400				{object}	handlers.HTTPError
//	@Failure		404				{object}	handlers.HTTPError
//	@Failure		500				{object}	handlers.HTTPError
//	@Router			/users/{id_user}/shelf-lives/{id_shelf_life} [delete]
//	@Security		Bearer
func (h *UserController) DeleteShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	shelfLifeID := ctx.Locals(context.ShelfLifeID).(int)
	if err := h.svc.DeleteShelfLife(ctx.Context(), id, shelfLifeID); err != nil {
		h.log.Error(ctx, logger.Server, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(controllers.HTTPSuccess{Success: true})
}

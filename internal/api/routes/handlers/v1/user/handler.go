package user

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/common"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/user"
)

type UserHanlder struct {
	svc service.UserServicer
	log *log.Logger
}

func New(svc service.UserServicer, log *log.Logger) *UserHanlder {
	return &UserHanlder{svc: svc, log: log}
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
func (h *UserHanlder) FindOne(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	user, err := h.svc.FindUserByID(ctx.Context(), id)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"user": user}})
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
func (h *UserHanlder) FindMany(ctx *fiber.Ctx) error {
	filter := new(dto.UserFilter)
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
	fmt.Printf("%+v\n", filter)
	result, err := h.svc.FindUsers(ctx.Context(), filter)
	fmt.Println(result)
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
		handlers.HTTPSuccess{Success: true, Data: handlers.Data{"users": result, "count": count}},
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
func (h *UserHanlder) Create(ctx *fiber.Ctx) error {
	payload := new(dto.CreateUser)
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
	if err := h.svc.CreateUser(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *UserHanlder) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	payload := new(dto.UpdateUser)
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
	if err := h.svc.UpdateUser(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *UserHanlder) Delete(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	if err := h.svc.DeleteUser(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *UserHanlder) Restore(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	if err := h.svc.RestoreUser(ctx.Context(), id); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *UserHanlder) FindSettings(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	result, err := h.svc.FindSettings(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"settings": result}})
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
func (h *UserHanlder) UpdateSetting(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	payload := new(dto.UpdateUserSetting)
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
	result, err := h.svc.UpdateSetting(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"settings": result}})
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
func (h *UserHanlder) FindRoles(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	result, err := h.svc.FindRoles(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"roles": result}})
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
func (h *UserHanlder) FindStorages(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	result, err := h.svc.FindStorages(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storages": result}})
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
func (h *UserHanlder) AddStorage(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	storageID := ctx.Locals(context.StorageID).(int)
	result, err := h.svc.AddStorage(ctx.Context(), id, storageID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storages": result}})
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
func (h *UserHanlder) RemoveStorage(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	storageID := ctx.Locals(context.StorageID).(int)
	err := h.svc.RemoveStorage(ctx.Context(), id, storageID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
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
func (h *UserHanlder) FindShelfLives(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	result, err := h.svc.FindShelfLives(ctx.Context(), id)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf-lives": result}})
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
func (h *UserHanlder) CreateShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	payload := new(dto.CreateShelfLife)
	payload.UserID = id
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
	result, err := h.svc.CreateShelfLife(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf-life": result}})
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
func (h *UserHanlder) UpdateShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	shelfLifeID := ctx.Locals(context.ShelfLifeID).(int)
	payload := new(dto.UserShelfLife)
	payload.ShelfLifeID = shelfLifeID
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
	result, err := h.svc.UpdateShelfLife(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf-life": result}})
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
func (h *UserHanlder) RestoreShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	shelfLifeID := ctx.Locals(context.ShelfLifeID).(int)
	result, err := h.svc.RestoreShelfLife(ctx.Context(), id, shelfLifeID)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf-life": result}})
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
func (h *UserHanlder) DeleteShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	shelfLifeID := ctx.Locals(context.ShelfLifeID).(int)
	if err := h.svc.DeleteShelfLife(ctx.Context(), id, shelfLifeID); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

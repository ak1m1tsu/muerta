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

func (h *UserHanlder) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	user, err := h.svc.FindUserByID(ctx.Context(), id)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrNotFound
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"user": user}})
}

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

func (h *UserHanlder) Create(ctx *fiber.Ctx) error {
	var payload *dto.CreateUser
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
	if err := h.svc.CreateUser(ctx.Context(), payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			JSON(handlers.HTTPError{Error: fiber.ErrBadGateway.Error()})
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

func (h *UserHanlder) Update(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	var payload *dto.UpdateUser
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(http.StatusBadRequest).
			JSON(handlers.HTTPError{Error: fiber.ErrBadRequest.Error()})
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
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

func (h *UserHanlder) UpdateSetting(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	var payload *dto.UpdateUserSetting
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			SendString("Bad payload provided")
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			SendString("Validation error")
	}
	result, err := h.svc.UpdateSetting(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"settings": result}})
}

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

func (h *UserHanlder) CreateStorage(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	var payload *dto.UserStorage
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			SendString("Bad payload provided")
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			SendString(errs.Error())
	}
	result, err := h.svc.CreateStorage(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"storages": result}})
}

func (h *UserHanlder) DeleteStorage(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	var payload *dto.UserStorage
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			SendString("Bad payload provided")
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			SendString(errs.Error())
	}
	err := h.svc.DeleteStorage(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

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

func (h *UserHanlder) CreateShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	var payload *dto.CreateShelfLife
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			SendString("Bad body provided")
	}
	payload.UserID = id
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			SendString("Validation error")
	}
	result, err := h.svc.CreateShelfLife(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf-life": result}})
}

func (h *UserHanlder) UpdateShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	var payload *dto.UserShelfLife
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			SendString("Bad body provided")
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			SendString("Validation error")
	}
	result, err := h.svc.UpdateShelfLife(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf-life": result}})
}

func (h *UserHanlder) RestoreShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	var payload *dto.UserShelfLife
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			SendString("Bad body provided")
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			SendString("Validation error")
	}
	result, err := h.svc.RestoreShelfLife(ctx.Context(), id, payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true, Data: handlers.Data{"shelf-life": result}})
}

func (h *UserHanlder) DeleteShelfLife(ctx *fiber.Ctx) error {
	id := ctx.Locals(context.UserID).(int)
	var payload *dto.UserShelfLife
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return ctx.Status(http.StatusBadRequest).
			SendString("Bad body provided")
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return ctx.Status(http.StatusBadRequest).
			SendString("Validation error")
	}
	if err := h.svc.DeleteShelfLife(ctx.Context(), id, payload); err != nil {
		h.log.ServerError(ctx, err)
		return ctx.Status(http.StatusBadGateway).
			SendString("Bad Gateway")
	}
	return ctx.JSON(handlers.HTTPSuccess{Success: true})
}

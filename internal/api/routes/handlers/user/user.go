package user

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	service "github.com/romankravchuk/muerta/internal/services/user"
)

type UserHanlder struct {
	svc service.UserServicer
}

func New(svc service.UserServicer) *UserHanlder {
	return &UserHanlder{svc: svc}
}

func (h *UserHanlder) FindByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	user, err := h.svc.FindByID(ctx.Context(), id)
	if err != nil {
		return fiber.ErrNotFound
	}
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"data": fiber.Map{
			"user": user,
		}})
}

func (h *UserHanlder) FindMany(ctx *fiber.Ctx) error {
	users, err := h.svc.FindMany(ctx.Context(), "")
	if err != nil {
		return fiber.ErrNotFound
	}
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"data": fiber.Map{
			"users": users,
		}})
}

func (h *UserHanlder) FindByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	user, err := h.svc.FindByName(ctx.Context(), name)
	if err != nil {
		return fiber.ErrNotFound
	}
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"data": fiber.Map{
			"user": user,
		}})
}

func (h *UserHanlder) Create(ctx *fiber.Ctx) error {
	user, err := h.svc.Create(ctx.Context(), "")
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"data": fiber.Map{
			"user": user,
		}})
}

type UpdateUserPayload struct {
	ID int `json:"id,omitempty"`
}

func (h *UserHanlder) Update(ctx *fiber.Ctx) error {
	var payload *UpdateUserPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.ErrBadRequest
	}
	user, err := h.svc.Update(ctx.Context(), payload.ID, "")
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"data": fiber.Map{
			"user": user,
		}})
}

type DeleteUserPayload struct {
	ID int `json:"id,omitempty"`
}

func (h *UserHanlder) Delete(ctx *fiber.Ctx) error {
	var payload *DeleteUserPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.ErrBadRequest
	}
	err := h.svc.Delete(ctx.Context(), payload.ID)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"data": fiber.Map{
			"success": true,
		}})
}

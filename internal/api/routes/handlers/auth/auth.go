package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/services/auth"
	"github.com/romankravchuk/muerta/internal/services/jwt"
)

type AuthHandler struct {
	jsvc jwt.JWTServicer
	asvc auth.AuthServicer
	log  *log.Logger
}

func New(jsvc jwt.JWTServicer, asvc auth.AuthServicer, log *log.Logger) *AuthHandler {
	return &AuthHandler{jsvc: jsvc, asvc: asvc, log: log}
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var payload *dto.LoginUserPayload
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	user, err := h.asvc.FindUser(ctx.Context(), payload)
	if err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	tokenPayload := dto.TokenPayload{
		Name:  user.Name,
		Roles: []string{"user"},
	}
	token, err := h.jsvc.CreateToken(tokenPayload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"success": true})
}

func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	var payload *dto.SignUpUserPayload
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	hash := h.asvc.HashPassword(payload.Password)
	if len(hash) == 0 {
		h.log.ServerError(ctx, errors.New("hash password length is zero"))
		return fiber.ErrInternalServerError
	}
	_ = dto.CreateUserPayload{
		SignUpUserPayload: *payload,
		Hash:              hash,
	}
	return ctx.Status(fiber.StatusOK).
		JSON(fiber.Map{"success": true})
}

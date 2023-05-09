package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/dto"
	"github.com/romankravchuk/muerta/internal/api/validator"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/auth"
)

type AuthHandler struct {
	svc         service.AuthServicer
	log         *log.Logger
	tokenMaxAge int
}

func New(cfg *config.Config, svc service.AuthServicer, log *log.Logger) *AuthHandler {
	return &AuthHandler{svc: svc, log: log, tokenMaxAge: cfg.TokenMaxAge}
}

func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	var payload *dto.SignUpDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrBadRequest
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrBadRequest
	}
	if payload.Password != payload.PasswordConfirm {
		h.log.ClientError(ctx, fmt.Errorf("passwords do not match"))
		return fiber.ErrBadRequest
	}
	if err := h.svc.SignUpUser(ctx.Context(), payload); err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			h.log.ClientError(ctx, err)
			return fiber.ErrBadRequest
		}
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var payload *dto.LoginDTO
	if err := ctx.BodyParser(&payload); err != nil {
		h.log.ClientError(ctx, err)
		return fiber.ErrUnauthorized
	}
	if errs := validator.Validate(payload); errs != nil {
		h.log.ValidationError(ctx, errs)
		return fiber.ErrUnauthorized
	}
	token, err := h.svc.LoginUser(ctx.Context(), payload)
	if err != nil {
		h.log.ServerError(ctx, err)
		return fiber.ErrBadGateway
	}
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   h.tokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   h.tokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
	})
	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token": token,
		},
	})
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	ctx.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expired,
	})
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

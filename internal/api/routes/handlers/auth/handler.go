package auth

import (
	"github.com/gofiber/fiber/v2"
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
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (h *AuthHandler) Refresh(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

package jwt

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/services/jwt"
)

type jwtMiddleware struct {
	svc jwt.JWTServicer
	log *log.Logger
}

func New(svc jwt.JWTServicer, log *log.Logger) *jwtMiddleware {
	return &jwtMiddleware{
		svc: svc,
		log: log,
	}
}

func (m *jwtMiddleware) DeserializeUser(ctx *fiber.Ctx) error {
	var token string
	authorization := ctx.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		token = strings.TrimPrefix(authorization, "Bearer ")
	} else {
		token = ctx.Cookies("access_token")
	}
	if token == "" {
		m.log.ClientError(ctx, fmt.Errorf("unauthorized request"))
		return fiber.ErrUnauthorized
	}
	payload, err := m.svc.ValidateToken(token)
	if err != nil {
		m.log.ClientError(ctx, err)
		return fiber.ErrForbidden
	}
	ctx.Locals("user", payload)
	return ctx.Next()
}

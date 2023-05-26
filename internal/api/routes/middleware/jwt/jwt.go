package jwt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/handlers"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
)

type JWTMiddleware struct {
	log          logger.Logger
	accessPubKey []byte
}

func New(cfg *config.Config, log logger.Logger) *JWTMiddleware {
	return &JWTMiddleware{
		log:          log,
		accessPubKey: cfg.AccessTokenPublicKey,
	}
}

func (m *JWTMiddleware) DeserializeUser(ctx *fiber.Ctx) error {
	var token string
	authorization := ctx.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		token = strings.TrimPrefix(authorization, "Bearer ")
	} else {
		token = ctx.Cookies("access_token")
	}
	if token == "" {
		m.log.Error(ctx, logger.Client, fmt.Errorf("unauthorized request"))
		return ctx.Status(http.StatusUnauthorized).
			JSON(handlers.HTTPError{Error: fiber.ErrUnauthorized.Error()})
	}
	payload, err := jwt.ValidateToken(token, m.accessPubKey)
	if err != nil {
		m.log.Error(ctx, logger.Client, err)
		return ctx.Status(http.StatusForbidden).
			JSON(handlers.HTTPError{Error: fiber.ErrForbidden.Error()})
	}
	ctx.Locals("access_token_uuid", payload.UUID)
	ctx.Locals("user", payload)
	return ctx.Next()
}

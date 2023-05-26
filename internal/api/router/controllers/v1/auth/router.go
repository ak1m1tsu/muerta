package auth

import (
	"github.com/gofiber/fiber/v2"

	jware "github.com/romankravchuk/muerta/internal/api/router/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	"github.com/romankravchuk/muerta/internal/services/auth"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/postgres/role"
	"github.com/romankravchuk/muerta/internal/storage/postgres/user"
	"github.com/romankravchuk/muerta/internal/storage/redis"
)

func NewRouter(
	cfg *config.Config,
	db postgres.Client,
	logger logger.Logger,
	redis redis.Client,
	jware *jware.JWTMiddleware,
) *fiber.App {
	userRepo := user.New(db)
	roleRepo := role.New(db)
	svc := auth.New(cfg, userRepo, roleRepo, redis)
	r := fiber.New()
	h := New(cfg, svc, logger)
	r.Post("/sign-up", h.SignUp)
	r.Post("/login", h.Login)
	r.Post("/logout", jware.DeserializeUser, h.Logout)
	r.Post("/refresh", h.RefreshAccessToken)
	return r
}

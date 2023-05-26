package auth

import (
	"github.com/gofiber/fiber/v2"

	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/services/auth"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	"github.com/romankravchuk/muerta/internal/storage/postgres/role"
	"github.com/romankravchuk/muerta/internal/storage/postgres/user"
	"github.com/romankravchuk/muerta/internal/storage/redis"
)

func NewRouter(
	cfg *config.Config,
	client postgres.Client,
	logger *log.Logger,
	redis redis.Client,
	jware *jware.JWTMiddleware,
) *fiber.App {
	userRepo := user.New(client)
	roleRepo := role.New(client)
	svc := auth.New(cfg, userRepo, roleRepo, redis)
	r := fiber.New()
	h := New(cfg, svc, logger)
	r.Post("/sign-up", h.SignUp)
	r.Post("/login", h.Login)
	r.Post("/logout", jware.DeserializeUser, h.Logout)
	r.Post("/refresh", h.RefreshAccessToken)
	return r
}

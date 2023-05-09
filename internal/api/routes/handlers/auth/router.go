package auth

import (
	"github.com/gofiber/fiber/v2"

	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/user"
	"github.com/romankravchuk/muerta/internal/services/auth"
	"github.com/romankravchuk/muerta/internal/services/jwt"
)

func NewRouter(cfg *config.Config, client repositories.PostgresClient, logger *log.Logger) *fiber.App {
	repo := user.New(client)
	jsvc := jwt.New(cfg)
	jware := jware.New(jsvc, logger)
	svc := auth.New(jsvc, repo)
	r := fiber.New()
	h := New(cfg, svc, logger)
	r.Post("/sign-up", h.SignUp)
	r.Post("/login", h.Login)
	r.Post("/logout", jware.DeserializeUser, h.Logout)
	return r
}

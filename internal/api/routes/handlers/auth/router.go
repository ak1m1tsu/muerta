package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories/user"
	"github.com/romankravchuk/muerta/internal/services/auth"
	"github.com/romankravchuk/muerta/internal/services/jwt"
)

func NewRouter(cfg *config.Config, db *sqlx.DB, logger *log.Logger) *fiber.App {
	repo := user.New(db)
	jsvc := jwt.New(cfg)
	asvc := auth.New(jsvc, repo)
	r := fiber.New()
	h := New(jsvc, asvc, logger)
	r.Post("/sign-up", h.SignUp)
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	return r
}

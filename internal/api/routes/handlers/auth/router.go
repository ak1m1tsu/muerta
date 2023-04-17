package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *fiber.App {
	r := fiber.New()
	h := New()
	r.Post("/sign-up", h.SignUp)
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)
	return r
}

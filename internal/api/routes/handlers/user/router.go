package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
	svc "github.com/romankravchuk/muerta/internal/services/user"
)

func NewRouter(db *sqlx.DB, log *log.Logger) *fiber.App {
	r := fiber.New()
	repo := repo.New(db)
	svc := svc.New(repo)
	h := New(svc, log)
	r.Get("/", h.FindMany)
	r.Get("/:id<int>", h.FindByID)
	r.Get("/:name<alpha>", h.FindByName)
	r.Post("/", h.Create)
	r.Put("/", h.Update)
	r.Delete("/", h.Delete)
	return r
}

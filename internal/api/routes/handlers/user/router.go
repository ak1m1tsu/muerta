package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
	svc "github.com/romankravchuk/muerta/internal/services/user"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	r := fiber.New()
	repo := repo.New(client)
	svc := svc.New(repo)
	h := New(svc, log)
	r.Get("/", h.FindMany)
	r.Post("/", h.Create)
	r.Route("/:id<int>", func(r fiber.Router) {
		r.Get("/", h.FindByID)
		r.Put("/", h.Update)
		r.Patch("/", h.Restore)
		r.Delete("/", h.Delete)
	})
	return r
}

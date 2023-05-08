package storage

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/storage"
	svc "github.com/romankravchuk/muerta/internal/services/storage"
)

func NewRouter(client repositories.PostgresClient, logger *log.Logger) *fiber.App {
	router := fiber.New()
	repo := repo.New(client)
	svc := svc.New(repo)
	handler := New(svc, logger)
	router.Get("/", handler.FindMany)
	router.Post("/", handler.Create)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindOne)
		router.Delete("/", handler.Delete)
		router.Put("/", handler.Update)
		router.Patch("/", handler.Restore)
	})
	return router
}

package shelflife

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/shelf-life"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindShelfLives)
	router.Post("/", handler.CreateShelfLife)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindShelfLifeByID)
		router.Put("/", handler.UpdateShelfLife)
		router.Delete("/", handler.DeleteShelfLife)
		router.Patch("/", handler.RestoreShelfLife)
	})
	return router
}

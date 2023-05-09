package storagetype

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/storage-type"
	service "github.com/romankravchuk/muerta/internal/services/storage-type"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindStorageTypes)
	router.Post("/", handler.CreateStorageType)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindStorageTypeByID)
		router.Delete("/", handler.DeleteStorageType)
		router.Put("/", handler.UpdateStorageType)
	})
	return router
}

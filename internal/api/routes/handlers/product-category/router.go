package productcategory

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/category"
	service "github.com/romankravchuk/muerta/internal/services/category"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", handler.Create)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindOne)
		router.Put("/", handler.Update)
		router.Patch("/", handler.Restore)
		router.Delete("/", handler.Delete)
	})
	return router
}

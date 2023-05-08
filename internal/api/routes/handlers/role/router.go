package role

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/role"
	service "github.com/romankravchuk/muerta/internal/services/role"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindRoles)
	router.Post("/", handler.CreateRole)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindRole)
		router.Put("/", handler.UpdateRole)
		router.Patch("/", handler.RestoreRole)
		router.Delete("/", handler.DeleteRole)
	})
	return router
}

package usersetting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/setting"
	usersetting "github.com/romankravchuk/muerta/internal/services/user-setting"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repo := setting.New(client)
	svc := usersetting.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", handler.Create)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindByID)
		router.Put("/", handler.Update)
		router.Patch("/", handler.Restore)
		router.Delete("/", handler.Delete)
	})
	return router
}

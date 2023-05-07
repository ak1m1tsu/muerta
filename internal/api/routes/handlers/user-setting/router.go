package usersetting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/settings"
	usersettings "github.com/romankravchuk/muerta/internal/services/user-settings"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repo := settings.New(client)
	svc := usersettings.New(repo)
	hanlder := New(svc, log)
	router.Get("/", hanlder.FindMany)
	router.Post("/", hanlder.Create)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", hanlder.FindByID)
	})
	return router
}

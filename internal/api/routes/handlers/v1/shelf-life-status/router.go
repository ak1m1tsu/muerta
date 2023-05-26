package shelflifestatus

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/access"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life-status"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/shelf-life-status"
)

func NewRouter(
	client postgres.Client,
	log *log.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.Create)
	router.Route("/:id", func(router fiber.Router) {
		router.Get("/", handler.FindOne)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
	})
	return router
}

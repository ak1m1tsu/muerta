package shelflifestatus

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/shelf-life-status"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life-status"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindShelfLifeStatuses)
	router.Post("/", jware.DeserializeUser, handler.CreateShelfLifeStatus)
	router.Route("/:id", func(router fiber.Router) {
		router.Get("/", handler.FindShelfLifeStatusByID)
		router.Put("/", jware.DeserializeUser, handler.UpdateShelfLifeStatus)
		router.Delete("/", jware.DeserializeUser, handler.DeleteShelfLifeStatus)
	})
	return router
}

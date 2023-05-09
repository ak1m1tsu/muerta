package shelflife

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/shelf-life"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindShelfLives)
	router.Post("/", jware.DeserializeUser, handler.CreateShelfLife)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindShelfLifeByID)
		router.Put("/", jware.DeserializeUser, handler.UpdateShelfLife)
		router.Delete("/", jware.DeserializeUser, handler.DeleteShelfLife)
		router.Patch("/", jware.DeserializeUser, handler.RestoreShelfLife)
	})
	return router
}

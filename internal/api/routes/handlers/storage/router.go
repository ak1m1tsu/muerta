package storage

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/storage"
	svc "github.com/romankravchuk/muerta/internal/services/storage"
)

func NewRouter(client repositories.PostgresClient, logger *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repo := repo.New(client)
	svc := svc.New(repo)
	handler := New(svc, logger)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, handler.Create)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindOne)
		router.Delete("/", jware.DeserializeUser, handler.Delete)
		router.Put("/", jware.DeserializeUser, handler.Update)
		router.Patch("/", jware.DeserializeUser, handler.Restore)
	})
	return router
}

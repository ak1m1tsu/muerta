package storagetype

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/storage-type"
	service "github.com/romankravchuk/muerta/internal/services/storage-type"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindStorageTypes)
	router.Post("/", jware.DeserializeUser, handler.CreateStorageType)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindStorageTypeByID)
		router.Delete("/", jware.DeserializeUser, handler.DeleteStorageType)
		router.Put("/", jware.DeserializeUser, handler.UpdateStorageType)
	})
	return router
}

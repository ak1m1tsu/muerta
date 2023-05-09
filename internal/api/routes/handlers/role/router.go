package role

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/role"
	service "github.com/romankravchuk/muerta/internal/services/role"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindRoles)
	router.Post("/", jware.DeserializeUser, handler.CreateRole)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindRole)
		router.Put("/", jware.DeserializeUser, handler.UpdateRole)
		router.Patch("/", jware.DeserializeUser, handler.RestoreRole)
		router.Delete("/", jware.DeserializeUser, handler.DeleteRole)
	})
	return router
}

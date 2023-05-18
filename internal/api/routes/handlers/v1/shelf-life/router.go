package shelflife

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/shelf-life"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life"
)

func NewRouter(
	client repositories.PostgresClient,
	log *log.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, handler.Create)
	router.Route(context.ShelfLifeID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.ShelfLifeID))
		router.Get("/", handler.FindOne)
		router.Put("/", jware.DeserializeUser, handler.Update)
		router.Delete("/", jware.DeserializeUser, handler.Delete)
		router.Patch("/", jware.DeserializeUser, handler.Restore)
		router.Route("/statuses", func(router fiber.Router) {
			router.Get("/", handler.FindStatuses)
			router.Route(context.StatusID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.StatusID))
				router.Post("/", jware.DeserializeUser, handler.AddStatus)
				router.Delete("/", jware.DeserializeUser, handler.RemoveStatus)
			})
		})
	})
	return router
}

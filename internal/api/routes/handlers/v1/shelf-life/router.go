package shelflife

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/shelf-life"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/shelf-life"
)

func NewRouter(
	client postgres.Client,
	log logger.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.Create)
	router.Route(context.ShelfLifeID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.ShelfLifeID))
		router.Get("/", handler.FindOne)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
		router.Patch("/", jware.DeserializeUser, access.AdminOnly(log), handler.Restore)
		router.Route("/statuses", func(router fiber.Router) {
			router.Get("/", handler.FindStatuses)
			router.Route(context.StatusID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.StatusID))
				router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.AddStatus)
				router.Delete(
					"/",
					jware.DeserializeUser,
					access.AdminOnly(log),
					handler.RemoveStatus,
				)
			})
		})
	})
	return router
}

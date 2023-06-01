package measure

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/router/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/measure"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/measure"
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
	router.Route(context.MeasureID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.MeasureID))
		router.Get("/", handler.FindOne)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
	})
	return router
}

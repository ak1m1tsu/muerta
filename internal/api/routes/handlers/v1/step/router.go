package step

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/step"
	service "github.com/romankravchuk/muerta/internal/services/step"
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
	router.Get("/", handler.FindSteps)
	router.Post("/", jware.DeserializeUser, handler.CreateStep)
	router.Route(context.StepID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.StepID))
		router.Get("/", handler.FindStep)
		router.Put("/", jware.DeserializeUser, handler.UpdateStep)
		router.Delete("/", jware.DeserializeUser, handler.DeleteStep)
		router.Patch("/", jware.DeserializeUser, handler.RestoreStep)
	})
	return router
}

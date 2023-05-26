package productcategory

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	service "github.com/romankravchuk/muerta/internal/services/category"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/category"
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
	router.Route(context.CategoryID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.CategoryID))
		router.Get("/", handler.FindOne)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Patch("/", jware.DeserializeUser, access.AdminOnly(log), handler.Restore)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
	})
	return router
}

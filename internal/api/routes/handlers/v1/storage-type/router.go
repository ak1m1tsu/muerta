package storagetype

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/storage-type"
	service "github.com/romankravchuk/muerta/internal/services/storage-type"
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
	router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.Create)
	router.Route(context.TypeID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.TypeID))
		router.Get("/", handler.FindOne)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Route("/storages", func(router fiber.Router) {
			router.Get("/", handler.FindStorages)
		})
		router.Route("/tips", func(router fiber.Router) {
			router.Get("/", handler.FindTips)
			router.Route(context.TipID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.TipID))
				router.Use(jware.DeserializeUser)
				router.Use(access.AdminOnly(log))
				router.Post("/", handler.AddTip)
				router.Delete("/", handler.RemoveTip)
			})
		})
	})
	return router
}

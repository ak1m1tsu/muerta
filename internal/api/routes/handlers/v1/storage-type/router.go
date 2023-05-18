package storagetype

import (
	"github.com/gofiber/fiber/v2"
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
	router.Get("/", handler.FindStorageTypes)
	router.Post("/", jware.DeserializeUser, handler.CreateStorageType)
	router.Route(context.TypeID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.TypeID))
		router.Get("/", handler.FindStorageTypeByID)
		router.Delete("/", jware.DeserializeUser, handler.DeleteStorageType)
		router.Put("/", jware.DeserializeUser, handler.UpdateStorageType)
		router.Route("/storages", func(router fiber.Router) {
			router.Get("/", handler.FindStorages)
		})
		router.Route("/tips", func(router fiber.Router) {
			router.Get("/", handler.FindTips)
			router.Route(context.TipID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.TipID))
				router.Post("/", jware.DeserializeUser, handler.CreateTip)
				router.Delete("/", jware.DeserializeUser, handler.DeleteTip)
			})
		})
	})
	return router
}

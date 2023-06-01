package tip

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/router/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/tip"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/tip"
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
	router.Route(context.TipID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.TipID))
		router.Get("/", handler.FindOne)
		router.Route("/products", func(router fiber.Router) {
			router.Get("/", handler.FindProducts)
			router.Route(context.ProductID.Path(), func(router fiber.Router) {
				router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.AddProduct)
				router.Delete(
					"/",
					jware.DeserializeUser,
					access.AdminOnly(log),
					handler.RemoveProduct,
				)
			})
		})
		router.Route("/storages", func(router fiber.Router) {
			router.Get("/", handler.FindStorages)
			router.Route(context.StorageID.Path(), func(router fiber.Router) {
				router.Post("/", jware.DeserializeUser, handler.AddStorage)
				router.Delete(
					"/",
					jware.DeserializeUser,
					access.AdminOnly(log),
					handler.RemoveStorage,
				)
			})
		})
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
		router.Patch("/", jware.DeserializeUser, access.AdminOnly(log), handler.Restore)
	})
	return router
}

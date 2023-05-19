package tip

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/tip"
	service "github.com/romankravchuk/muerta/internal/services/tip"
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
	router.Route(context.TipID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.TipID))
		router.Get("/", handler.FindOne)
		router.Route("/products", func(router fiber.Router) {
			router.Get("/", handler.FindProducts)
			router.Route(context.ProductID.Path(), func(router fiber.Router) {
				router.Post("/", jware.DeserializeUser, handler.AddProduct)
				router.Delete("/", jware.DeserializeUser, handler.RemoveProduct)
			})
		})
		router.Route("/storages", func(router fiber.Router) {
			router.Get("/", handler.FindStorages)
			router.Route(context.StorageID.Path(), func(router fiber.Router) {
				router.Post("/", jware.DeserializeUser, handler.AddStorage)
				router.Delete("/", jware.DeserializeUser, handler.RemoveStorage)
			})
		})
		router.Put("/", jware.DeserializeUser, handler.Update)
		router.Delete("/", jware.DeserializeUser, handler.Delete)
		router.Patch("/", jware.DeserializeUser, handler.Restore)
	})
	return router
}

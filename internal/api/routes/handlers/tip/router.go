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

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindTips)
	router.Post("/", jware.DeserializeUser, handler.CreateTip)
	router.Route(context.TipID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.TipID))
		router.Get("/", handler.FindTipByID)
		router.Route("/products", func(router fiber.Router) {
			router.Get("/", handler.FindTipProducts)
			router.Route(context.ProductID.Path(), func(router fiber.Router) {
				router.Post("/", jware.DeserializeUser, handler.AddProductToTip)
				router.Delete("/", jware.DeserializeUser, handler.RemoveProductFromTip)
			})
		})
		router.Route("/storages", func(router fiber.Router) {
			router.Get("/", handler.FindTipStorages)
			router.Route(context.StorageID.Path(), func(router fiber.Router) {
				router.Post("/", jware.DeserializeUser, handler.AddStorageToTip)
				router.Delete("/", jware.DeserializeUser, handler.RemoveStorageFromTip)
			})
		})
		router.Put("/", jware.DeserializeUser, handler.UpdateTip)
		router.Delete("/", jware.DeserializeUser, handler.DeleteTip)
		router.Patch("/", jware.DeserializeUser, handler.RestoreTip)
	})
	return router
}

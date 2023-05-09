package tip

import (
	"github.com/gofiber/fiber/v2"
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
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindTipByID)
		router.Route("/products", func(router fiber.Router) {
			router.Get("/", handler.FindTipProducts)
		})
		router.Route("/storages", func(router fiber.Router) {
			router.Get("/", handler.FindTipStorages)
		})
		router.Put("/", jware.DeserializeUser, handler.UpdateTip)
		router.Delete("/", jware.DeserializeUser, handler.DeleteTip)
		router.Patch("/", jware.DeserializeUser, handler.RestoreTip)
	})
	return router
}

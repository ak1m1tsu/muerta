package tip

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repository "github.com/romankravchuk/muerta/internal/repositories/tip"
	service "github.com/romankravchuk/muerta/internal/services/tip"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindTips)
	router.Post("/", handler.CreateTip)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindTipByID)
		router.Put("/", handler.UpdateTip)
		router.Delete("/", handler.DeleteTip)
		router.Patch("/", handler.RestoreTip)
	})
	return router
}

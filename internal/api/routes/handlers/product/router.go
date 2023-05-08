package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/product"
	svc "github.com/romankravchuk/muerta/internal/services/product"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repository := repo.New(client)
	service := svc.New(repository)
	handler := New(service, log)
	router.Post("/", handler.CreateProduct)
	router.Get("/", handler.FindProducts)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindProductByID)
		router.Put("/", handler.UpdateProduct)
		router.Delete("/", handler.DeleteProduct)
		router.Patch("/", handler.RestoreProduct)
	})
	return router
}

package product

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/product"
	svc "github.com/romankravchuk/muerta/internal/services/product"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repository := repo.New(client)
	service := svc.New(repository)
	handler := New(service, log)
	router.Post("/", jware.DeserializeUser, handler.CreateProduct)
	router.Get("/", handler.FindProducts)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindProductByID)
		router.Route("/categories", func(router fiber.Router) {
			router.Get("/", handler.FindProductCategories)
		})
		router.Route("/recipes", func(router fiber.Router) {
			router.Get("/", handler.FindProductRecipes)
		})
		router.Put("/", jware.DeserializeUser, handler.UpdateProduct)
		router.Delete("/", jware.DeserializeUser, handler.DeleteProduct)
		router.Patch("/", jware.DeserializeUser, handler.RestoreProduct)
	})
	return router
}

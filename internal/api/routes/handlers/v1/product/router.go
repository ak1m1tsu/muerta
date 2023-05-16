package product

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
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
	router.Get("/", handler.FindProducts)
	router.Post("/", jware.DeserializeUser, handler.CreateProduct)
	router.Route(context.ProductID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.ProductID))
		router.Get("/", handler.FindProductByID)
		router.Route("/categories", func(router fiber.Router) {
			router.Get("/", handler.FindProductCategories)
			router.Route(context.CategoryID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.CategoryID))
				router.Post("/", handler.CreateCategory)
				router.Delete("/", handler.DeleteCategory)
			})
		})
		router.Route("/recipes", func(router fiber.Router) {
			router.Get("/", handler.FindProductRecipes)
		})
		router.Route("/tips", func(router fiber.Router) {
			router.Get("/", handler.FindProductTips)
			router.Route(context.TipID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.TipID))
				router.Post("/", handler.CreateProductTip)
				router.Delete("/", handler.DeleteProductTip)
			})
		})
		router.Put("/", jware.DeserializeUser, handler.UpdateProduct)
		router.Delete("/", jware.DeserializeUser, handler.DeleteProduct)
		router.Patch("/", jware.DeserializeUser, handler.RestoreProduct)
	})
	return router
}

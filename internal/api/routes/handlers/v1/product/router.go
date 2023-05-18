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

func NewRouter(
	client repositories.PostgresClient,
	log *log.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repository := repo.New(client)
	service := svc.New(repository)
	handler := New(service, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, handler.Create)
	router.Route(context.ProductID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.ProductID))
		router.Get("/", handler.FindOne)
		router.Route("/categories", func(router fiber.Router) {
			router.Get("/", handler.FindCategories)
			router.Route(context.CategoryID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.CategoryID))
				router.Post("/", handler.AddCategory)
				router.Delete("/", handler.RemoveCategory)
			})
		})
		router.Route("/recipes", func(router fiber.Router) {
			router.Get("/", handler.FindRecipes)
		})
		router.Route("/tips", func(router fiber.Router) {
			router.Get("/", handler.FindTips)
			router.Route(context.TipID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.TipID))
				router.Post("/", handler.AddTip)
				router.Delete("/", handler.RemoveTip)
			})
		})
		router.Put("/", jware.DeserializeUser, handler.Update)
		router.Delete("/", jware.DeserializeUser, handler.Delete)
		router.Patch("/", jware.DeserializeUser, handler.Restore)
	})
	return router
}

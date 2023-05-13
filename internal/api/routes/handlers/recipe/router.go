package recipe

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/recipe"
	svc "github.com/romankravchuk/muerta/internal/services/recipe"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repository := repo.New(client)
	service := svc.New(repository)
	handler := New(service, log)
	router.Get("/", handler.FindRecipes)
	router.Post("/", jware.DeserializeUser, handler.CreateRecipe)
	router.Route(context.RecipeID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.RecipeID))
		router.Get("/", handler.FindRecipeByID)
		router.Put("/", jware.DeserializeUser, handler.UpdateRecipe)
		router.Delete("/", jware.DeserializeUser, handler.DeleteRecipe)
		router.Patch("/", jware.DeserializeUser, handler.RestoreRecipe)
		router.Route("/ingredients", func(router fiber.Router) {
			router.Get("/", handler.FindRecipeIngredients)
			router.Post("/", jware.DeserializeUser, handler.CreateRecipeIngredient)
			router.Put("/", jware.DeserializeUser, handler.UpdateRecipeIngredient)
			router.Delete("/", jware.DeserializeUser, handler.DeleteRecipeIngredient)
		})
		router.Route("/steps", func(router fiber.Router) {
			router.Get("/", handler.FindRecipeSteps)
			router.Route(context.StepID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.StepID))
				router.Post("/", jware.DeserializeUser, handler.CreateRecipeStep)
				router.Delete("/", jware.DeserializeUser, handler.DeleteRecipeStep)
			})
		})
	})
	return router
}

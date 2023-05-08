package recipe

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/recipes"
	svc "github.com/romankravchuk/muerta/internal/services/recipe"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repository := repo.New(client)
	service := svc.New(repository)
	handler := New(service, log)
	router.Post("/", handler.CreateRecipe)
	router.Get("/", handler.FindRecipes)
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindRecipeByID)
		router.Put("/", handler.UpdateRecipe)
		router.Delete("/", handler.DeleteRecipe)
		router.Patch("/", handler.RestoreRecipe)
	})
	return router
}

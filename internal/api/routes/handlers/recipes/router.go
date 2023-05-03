package recipes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/recipes"
	svc "github.com/romankravchuk/muerta/internal/services/recipes"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger) *fiber.App {
	router := fiber.New()
	repository := repo.New(client)
	service := svc.New(repository)
	handler := New(service, log)
	router.Post("/", handler.CreateRecipe)
	router.Get("/", handler.FindRecipes)
	router.Get("/:id<int>", handler.FindRecipeByID)
	router.Put("/:id<int>", handler.UpdateRecipe)
	// router.Delete("/:id<int>", handler.DeleteRecipe)
	router.Get("/:name<alpha>", handler.FindRecipeByName)
	return router
}

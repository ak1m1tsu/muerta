package recipe

import (
	"github.com/gofiber/fiber/v2"
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
	router.Route("/:id<int>", func(router fiber.Router) {
		router.Get("/", handler.FindRecipeByID)
		router.Put("/", jware.DeserializeUser, handler.UpdateRecipe)
		router.Delete("/", jware.DeserializeUser, handler.DeleteRecipe)
		router.Patch("/", jware.DeserializeUser, handler.RestoreRecipe)
	})
	return router
}

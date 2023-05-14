package storage

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/storage"
	svc "github.com/romankravchuk/muerta/internal/services/storage"
)

func NewRouter(client repositories.PostgresClient, logger *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repo := repo.New(client)
	svc := svc.New(repo)
	handler := New(svc, logger)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, handler.Create)
	router.Route(context.StorageID.Path(), func(router fiber.Router) {
		router.Use(context.New(logger, context.StorageID))
		router.Get("/", handler.FindOne)
		router.Delete("/", jware.DeserializeUser, handler.Delete)
		router.Put("/", jware.DeserializeUser, handler.Update)
		router.Patch("/", jware.DeserializeUser, handler.Restore)
		router.Route("/tips", func(router fiber.Router) {
			router.Get("/", handler.FindTips)
			router.Route(context.TipID.Path(), func(router fiber.Router) {
				router.Use(context.New(logger, context.TipID))
				router.Post("/", jware.DeserializeUser, handler.CreateTip)
				router.Delete("/", jware.DeserializeUser, handler.DeleteTip)
			})
		})
		router.Route("/shelf-lives", func(router fiber.Router) {
			router.Get("/", handler.FindShelfLives)
		})
	})
	return router
}

package storage

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	svc "github.com/romankravchuk/muerta/internal/services/storage"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repo "github.com/romankravchuk/muerta/internal/storage/postgres/storage"
)

func NewRouter(
	client postgres.Client,
	log logger.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repo := repo.New(client)
	svc := svc.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.Create)
	router.Route(context.StorageID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.StorageID))
		router.Get("/", handler.FindOne)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Patch("/", jware.DeserializeUser, access.AdminOnly(log), handler.Restore)
		router.Route("/tips", func(router fiber.Router) {
			router.Get("/", handler.FindTips)
			router.Route(context.TipID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.TipID))
				router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.AddTip)
				router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.RemoveTip)
			})
		})
		router.Route("/shelf-lives", func(router fiber.Router) {
			router.Get("/", handler.FindShelfLives)
		})
	})
	return router
}

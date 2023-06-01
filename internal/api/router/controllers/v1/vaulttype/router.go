package vaulttype

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/router/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/router/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	service "github.com/romankravchuk/muerta/internal/services/storage-type"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repository "github.com/romankravchuk/muerta/internal/storage/postgres/storage-type"
)

func NewRouter(
	client postgres.Client,
	log logger.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	router := fiber.New()
	repo := repository.New(client)
	svc := service.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, access.AdminOnly(log), handler.Create)
	router.Route(context.TypeID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.TypeID))
		router.Get("/", handler.FindOne)
		router.Delete("/", jware.DeserializeUser, access.AdminOnly(log), handler.Delete)
		router.Put("/", jware.DeserializeUser, access.AdminOnly(log), handler.Update)
		router.Route("/storages", func(router fiber.Router) {
			router.Get("/", handler.FindStorages)
		})
		router.Route("/tips", func(router fiber.Router) {
			router.Get("/", handler.FindTips)
			router.Route(context.TipID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.TipID))
				router.Use(jware.DeserializeUser)
				router.Use(access.AdminOnly(log))
				router.Post("/", handler.AddTip)
				router.Delete("/", handler.RemoveTip)
			})
		})
	})
	return router
}

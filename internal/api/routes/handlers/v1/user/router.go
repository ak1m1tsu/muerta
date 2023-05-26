package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/access"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	svc "github.com/romankravchuk/muerta/internal/services/user"
	"github.com/romankravchuk/muerta/internal/storage/postgres"
	repo "github.com/romankravchuk/muerta/internal/storage/postgres/user"
)

func NewRouter(
	client postgres.Client,
	log logger.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	r := fiber.New()
	repo := repo.New(client)
	svc := svc.New(repo)
	h := New(svc, log)
	r.Get("/", h.FindMany)
	r.Post("/", jware.DeserializeUser, access.AdminOnly(log), h.Create)
	r.Route(context.UserID.Path(), func(r fiber.Router) {
		r.Use(context.New(log, context.UserID))
		r.Get("/", h.FindOne)
		r.Put("/", jware.DeserializeUser, access.OwnerOnly(log), h.Update)
		r.Patch("/", jware.DeserializeUser, access.OwnerOnly(log), h.Restore)
		r.Delete("/", jware.DeserializeUser, access.OwnerOnly(log), h.Delete)
		r.Route("/shelf-lives", func(router fiber.Router) {
			router.Get("/", h.FindShelfLives)
			router.Post("/", jware.DeserializeUser, access.OwnerOnly(log), h.CreateShelfLife)
			router.Route(context.ShelfLifeID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.ShelfLifeID))
				router.Put("/", jware.DeserializeUser, access.OwnerOnly(log), h.UpdateShelfLife)
				router.Patch("/", jware.DeserializeUser, access.OwnerOnly(log), h.RestoreShelfLife)
				router.Delete("/", jware.DeserializeUser, access.OwnerOnly(log), h.DeleteShelfLife)
			})
		})
		r.Route("/settings", func(router fiber.Router) {
			router.Get("/", jware.DeserializeUser, access.OwnerOnly(log), h.FindSettings)
			router.Route(context.SettingID.Path(), func(router fiber.Router) {
				router.Put("/", jware.DeserializeUser, access.OwnerOnly(log), h.UpdateSetting)
			})
		})
		r.Get("/roles", h.FindRoles)
		r.Route("/storages", func(router fiber.Router) {
			router.Get("/", h.FindStorages)
			router.Route(context.StorageID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.StorageID))
				router.Use(jware.DeserializeUser)
				router.Use(access.OwnerOnly(log))
				router.Post("/", h.AddStorage)
				router.Delete("/", h.RemoveStorage)
			})
		})
	})
	return r
}

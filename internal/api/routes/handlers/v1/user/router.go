package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
	svc "github.com/romankravchuk/muerta/internal/services/user"
)

func NewRouter(
	client repositories.PostgresClient,
	log *log.Logger,
	jware *jware.JWTMiddleware,
) *fiber.App {
	r := fiber.New()
	repo := repo.New(client)
	svc := svc.New(repo)
	h := New(svc, log)
	r.Get("/", h.FindMany)
	r.Post("/", jware.DeserializeUser, h.Create)
	r.Route(context.UserID.Path(), func(r fiber.Router) {
		r.Use(context.New(log, context.UserID))
		r.Get("/", h.FindOne)
		r.Put("/", jware.DeserializeUser, h.Update)
		r.Patch("/", jware.DeserializeUser, h.Restore)
		r.Delete("/", jware.DeserializeUser, h.Delete)
		r.Route("/shelf-lives", func(router fiber.Router) {
			router.Get("/", h.FindShelfLives)
			router.Post("/", jware.DeserializeUser, h.CreateShelfLife)
			router.Route(context.ShelfLifeID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.ShelfLifeID))
				router.Put("/", jware.DeserializeUser, h.UpdateShelfLife)
				router.Patch("/", jware.DeserializeUser, h.RestoreShelfLife)
				router.Delete("/", jware.DeserializeUser, h.DeleteShelfLife)
			})
		})
		r.Route("/settings", func(router fiber.Router) {
			router.Get("/", h.FindSettings)
			router.Route(context.SettingID.Path(), func(router fiber.Router) {
				router.Put("/", jware.DeserializeUser, h.UpdateSetting)
			})
		})
		r.Get("/roles", h.FindRoles)
		r.Route("/storages", func(router fiber.Router) {
			router.Get("/", h.FindStorages)
			router.Route(context.StorageID.Path(), func(router fiber.Router) {
				router.Use(context.New(log, context.StorageID))
				router.Post("/", jware.DeserializeUser, h.AddStorage)
				router.Delete("/", jware.DeserializeUser, h.RemoveStorage)
			})
		})
	})
	return r
}

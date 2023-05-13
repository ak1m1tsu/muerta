package user

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	repo "github.com/romankravchuk/muerta/internal/repositories/user"
	svc "github.com/romankravchuk/muerta/internal/services/user"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	r := fiber.New()
	repo := repo.New(client)
	svc := svc.New(repo)
	h := New(svc, log)
	r.Get("/", h.FindMany)
	r.Post("/", jware.DeserializeUser, h.Create)
	r.Route("/:id<int>", func(r fiber.Router) {
		r.Get("/", h.FindByID)
		r.Put("/", jware.DeserializeUser, h.Update)
		r.Patch("/", jware.DeserializeUser, h.Restore)
		r.Delete("/", jware.DeserializeUser, h.Delete)
		r.Route("/shelf-lives", func(router fiber.Router) {
			router.Get("/", h.FindShelfLives)
			router.Post("/", jware.DeserializeUser, h.CreateShelfLife)
			router.Put("/", jware.DeserializeUser, h.UpdateShelfLife)
			router.Patch("/", jware.DeserializeUser, h.RestoreShelfLife)
			router.Delete("/", jware.DeserializeUser, h.DeleteShelfLife)
		})
		r.Route("/settings", func(router fiber.Router) {
			router.Get("/", h.FindSettings)
			router.Put("/", jware.DeserializeUser, h.UpdateSetting)
		})
		r.Get("/roles", h.FindRoles)
		r.Route("/storages", func(router fiber.Router) {
			router.Get("/", h.FindStorages)
			router.Post("/", jware.DeserializeUser, h.CreateStorage)
			router.Delete("/", jware.DeserializeUser, h.DeleteStorage)
		})
	})
	return r
}

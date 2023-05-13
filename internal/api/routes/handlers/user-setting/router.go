package usersetting

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/api/routes/middleware/context"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	"github.com/romankravchuk/muerta/internal/repositories"
	"github.com/romankravchuk/muerta/internal/repositories/setting"
	usersetting "github.com/romankravchuk/muerta/internal/services/user-setting"
)

func NewRouter(client repositories.PostgresClient, log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	repo := setting.New(client)
	svc := usersetting.New(repo)
	handler := New(svc, log)
	router.Get("/", handler.FindMany)
	router.Post("/", jware.DeserializeUser, handler.Create)
	router.Route(context.SettingID.Path(), func(router fiber.Router) {
		router.Use(context.New(log, context.SettingID))
		router.Get("/", handler.FindByID)
		router.Put("/", jware.DeserializeUser, handler.Update)
		router.Patch("/", jware.DeserializeUser, handler.Restore)
		router.Delete("/", jware.DeserializeUser, handler.Delete)
	})
	return router
}

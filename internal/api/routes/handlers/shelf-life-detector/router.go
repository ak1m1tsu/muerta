package shelflifedetector

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/routes/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	sldetector "github.com/romankravchuk/muerta/internal/services/shelf-life-detector"
)

func NewRouter(log *log.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	service := sldetector.New()
	handler := New(service, log)
	router.Post("/", jware.DeserializeUser, handler.DetectDates)
	return router
}

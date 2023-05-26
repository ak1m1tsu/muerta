package shelflifedetector

import (
	"github.com/gofiber/fiber/v2"
	jware "github.com/romankravchuk/muerta/internal/api/router/middleware/jwt"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/logger"
	sldetector "github.com/romankravchuk/muerta/internal/services/shelf-life-detector"
)

func NewRouter(cfg *config.Config, log logger.Logger, jware *jware.JWTMiddleware) *fiber.App {
	router := fiber.New()
	service := sldetector.New(cfg.ShutdownShelfDetectorChan)
	handler := New(service, log)
	router.Post("/", jware.DeserializeUser, handler.DetectDates)
	return router
}

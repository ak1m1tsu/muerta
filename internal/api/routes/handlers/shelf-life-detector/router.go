package shelflifedetector

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romankravchuk/muerta/internal/pkg/log"
	sldetector "github.com/romankravchuk/muerta/internal/services/shelf-life-detector"
)

func NewRouter(log *log.Logger) *fiber.App {
	router := fiber.New()
	service := sldetector.New()
	handler := New(service, log)
	router.Post("/", handler.DetectDates)
	return router
}

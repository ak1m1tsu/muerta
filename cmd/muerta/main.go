package main

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:               "Muerta API v1.0",
		DisableStartupMessage: true,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
	})
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "welcome!"}) })
	log.Fatal(app.Listen(":3000"))
}

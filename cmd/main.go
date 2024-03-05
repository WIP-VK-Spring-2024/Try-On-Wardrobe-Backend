package main

import (
	"log"

	"try-on/internal/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	api := fiber.New()

	recover := recover.New(recover.Config{
		EnableStackTrace: true,
	})

	logger := logger.New(logger.Config{
		Format: config.JsonLogFormat,
	})

	cors := cors.New(cors.Config{
		AllowOrigins:     "localhost:80",
		AllowCredentials: true,
		MaxAge:           -1,
	})

	api.Use(recover, logger, cors)

	api.Get("/", func(c *fiber.Ctx) error {
		_, err := c.WriteString("Hello, world!\n")
		return err
	})

	log.Fatal(api.Listen(":8000"))
}

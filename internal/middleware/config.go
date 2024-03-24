package middleware

import (
	"context"

	"try-on/internal/pkg/config"

	"github.com/gofiber/fiber/v2"
)

type configKeyType struct{}

var configKey configKeyType

func AddConfig(cfg *config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.UserContext(), configKey, cfg)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func Config(ctx *fiber.Ctx) *config.Config {
	cfg, ok := ctx.UserContext().Value(configKey).(*config.Config)
	if !ok {
		return nil
	}
	return cfg
}

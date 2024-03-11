package middleware

import (
	"context"
	"errors"

	"try-on/internal/pkg/app_errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type loggerKeyType string

const loggerKey loggerKeyType = "logger"

func AddLogger(logger *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(context.WithValue(ctx.UserContext(), loggerKey, logger))
		return ctx.Next()
	}
}

func LogError(ctx *fiber.Ctx, err error) {
	logger, ok := ctx.UserContext().Value(loggerKey).(*zap.SugaredLogger)
	if !ok {
		logger = zap.S()
	}

	var e *app_errors.Error
	if errors.As(err, &e) {
		logger.Errorw(err.Error(), "method", ctx.Method(), "path", ctx.Path(), "file", e.File, "line", e.Line)
	} else {
		logger.Errorw(err.Error(), "method", ctx.Method(), "path", ctx.Path())
	}
}

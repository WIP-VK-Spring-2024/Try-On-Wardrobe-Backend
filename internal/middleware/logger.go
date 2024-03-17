package middleware

import (
	"context"
	"errors"

	"try-on/internal/pkg/app_errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type loggerKeyType struct{}

var loggerKey loggerKeyType = struct{}{}

func AddLogger(logger *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(context.WithValue(ctx.UserContext(), loggerKey, logger))
		return ctx.Next()
	}
}

func GetLogger(ctx *fiber.Ctx) *zap.SugaredLogger {
	logger, ok := ctx.UserContext().Value(loggerKey).(*zap.SugaredLogger)
	if !ok {
		logger = zap.S()
	}
	return logger
}

func LogError(ctx *fiber.Ctx, err error) {
	logger := GetLogger(ctx)

	values := []interface{}{"method", ctx.Method(), "path", ctx.Path(), "ip", ctx.IP()}

	var e *app_errors.InternalError
	if errors.As(err, &e) {
		values = append(values, "file", e.File, "line", e.Line)
	}

	logger.Errorw(err.Error(), values...)
}

package middleware

import (
	"context"
	"errors"

	"try-on/internal/pkg/app_errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type loggerKeyType struct{}

var loggerKey loggerKeyType

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func AddLogger(logger *zap.SugaredLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(WithLogger(ctx.UserContext(), logger))
		return ctx.Next()
	}
}

func GetLogger(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger)
	if !ok {
		logger = zap.S()
	}
	return logger
}

func log(ctx *fiber.Ctx, err error, logfunc func(string, ...interface{}), fields ...interface{}) {
	if err == nil {
		return
	}

	values := []interface{}{"method", ctx.Method(), "path", ctx.Path(), "ip", ctx.IP()}

	var e *app_errors.InternalError
	if errors.As(err, &e) {
		values = append(values, "file", e.File, "line", e.Line)
	}

	values = append(values, fields...)
	logfunc(err.Error(), values...)
}

func LogError(ctx *fiber.Ctx, err error, fields ...interface{}) {
	logger := GetLogger(ctx.UserContext())
	log(ctx, err, logger.Errorw, fields...)
}

func LogWarning(ctx *fiber.Ctx, err error, fields ...interface{}) {
	logger := GetLogger(ctx.UserContext())
	log(ctx, err, logger.Warnw, fields...)
}

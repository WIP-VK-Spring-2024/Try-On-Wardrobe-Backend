package middleware

import (
	"context"
	"slices"
	"strings"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"

	"github.com/gofiber/fiber/v2"
)

type sessionKeyType struct{}

var sessionKey sessionKeyType

type SessionConfig struct {
	TokenName    string
	Sessions     domain.SessionUsecase
	NoAuthRoutes []string
	SecureRoutes []string
}

func CheckSession(cfg SessionConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if slices.Contains(cfg.NoAuthRoutes, ctx.Path()) {
			return ctx.Next()
		}

		session := domain.Session{
			ID: ctx.Get(cfg.TokenName),
		}

		ok, err := cfg.Sessions.IsLoggedIn(&session)
		if err != nil && err != app_errors.ErrInvalidCredentials {
			return err
		}

		if ok {
			context := context.WithValue(ctx.UserContext(), sessionKey, &session)
			ctx.SetUserContext(context)
		} else if slices.ContainsFunc(cfg.SecureRoutes, func(prefix string) bool {
			return strings.HasPrefix(ctx.Path(), prefix)
		}) {
			return app_errors.ErrUnauthorized
		}

		return ctx.Next()
	}
}

func Session(ctx *fiber.Ctx) *domain.Session {
	// userID, _ := utils.ParseUUID("2a78df8a-0277-4c72-a2d9-43fb8fef1d2c") // first account
	// userID, _ := utils.ParseUUID("7eef49f3-9b52-4dc8-bcd6-c40c6eb966bc") // ux account

	// return &domain.Session{
	// 	UserID: userID,
	// }

	value := ctx.UserContext().Value(sessionKey)
	session, ok := value.(*domain.Session)
	if !ok {
		return nil
	}
	return session
}

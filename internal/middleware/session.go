package middleware

import (
	"context"

	"try-on/internal/pkg/api_errors"
	"try-on/internal/pkg/domain"

	"github.com/gofiber/fiber/v2"
)

type sessionKeyType string

const sessionKey sessionKeyType = "session"

type SessionConfig struct {
	CookieName string
	Sessions   domain.SessionRepository
}

func AddSession(cfg SessionConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		cookie := ctx.Cookies(cfg.CookieName)

		session, err := cfg.Sessions.Get(cookie)
		if err != nil && err != api_errors.ErrInvalidCredentials {
			return err
		}

		context := context.WithValue(ctx.UserContext(), sessionKey, session)
		ctx.SetUserContext(context)

		return ctx.Next()
	}
}

func GetSession(ctx *fiber.Ctx) *domain.Session {
	value := ctx.UserContext().Value(sessionKey)
	session, ok := value.(*domain.Session)
	if !ok {
		return nil
	}
	return session
}

func CheckSession(ctx *fiber.Ctx) error {
	session := GetSession(ctx)
	if session == nil {
		return fiber.ErrUnauthorized
	}
	return ctx.Next()
}

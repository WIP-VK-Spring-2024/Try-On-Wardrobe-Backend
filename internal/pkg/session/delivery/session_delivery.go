package delivery

import (
	"encoding/json"
	"errors"

	"try-on/internal/pkg/api_errors"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	sessionRepo "try-on/internal/pkg/session/repository"
	sessionUsecase "try-on/internal/pkg/session/usecase"
	userRepo "try-on/internal/pkg/users/repository"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type SessionHandler struct {
	sessions domain.SessionUsecase
	cfg      config.Session
}

type Config struct {
	config.Session
	config.Redis
}

func NewSessionHandler(db *gorm.DB, cfg Config) *SessionHandler {
	return &SessionHandler{
		sessions: sessionUsecase.NewSessionUsecase(
			userRepo.NewUserRepository(db),
			sessionRepo.NewRedisSessionStorage(sessionRepo.Config{
				Namespace:     cfg.KeyNamespace,
				MaxConn:       cfg.MaxConn,
				RedisAddr:     cfg.Addr,
				ExpireSeconds: int64(cfg.MaxAge),
			}),
		),
		cfg: cfg.Session,
	}
}

func (h *SessionHandler) Register(ctx *fiber.Ctx) error {
	var credentials domain.Credentials

	err := json.Unmarshal(ctx.Body(), &credentials)
	if err != nil {
		return err
	}

	user := domain.User{
		Name:     credentials.Name,
		Password: credentials.Password,
	}

	session, err := h.sessions.Register(&user)
	switch {
	case err == nil:
		ctx.Cookie(getCookie(h.cfg.CookieName, session.ID, h.cfg.MaxAge))
		return ctx.SendString(utils.EmptyJson)

	case errors.Is(err, api_errors.ErrAlreadyExists):
		return fiber.ErrConflict

	case errors.Is(err, api_errors.ErrSessionNotInitialized):
		log.Warnw("user", credentials.Name, "error", err)
		return nil

	default:
		return err
	}
}

func (h *SessionHandler) Login(ctx *fiber.Ctx) error {
	var credentials domain.Credentials

	err := json.Unmarshal(ctx.Body(), &credentials)
	if err != nil {
		return err
	}

	session, err := h.sessions.Login(credentials)
	switch {
	case err == nil:
		ctx.Cookie(getCookie(h.cfg.CookieName, session.ID, h.cfg.MaxAge))
		return ctx.SendString(utils.EmptyJson)

	case errors.Is(err, api_errors.ErrSessionNotInitialized):
		log.Warnw("user", credentials.Name, "error", err)
		return nil

	case errors.Is(err, api_errors.ErrInvalidCredentials):
		return fiber.ErrForbidden

	default:
		return err
	}
}

func (h *SessionHandler) Logout(ctx *fiber.Ctx) error {
	sessionID := ctx.Cookies(h.cfg.CookieName)
	if sessionID == "" {
		return fiber.ErrUnauthorized
	}

	err := h.sessions.Logout(sessionID)
	if err != nil {
		return err
	}

	return ctx.SendString(utils.EmptyJson)
}

func getCookie(name, value string, maxAge int) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		SameSite: "strict",
		HTTPOnly: true,
		MaxAge:   maxAge,
	}
}

package delivery

import (
	"errors"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	sessionUsecase "try-on/internal/pkg/session/usecase"
	userRepo "try-on/internal/pkg/users/repository"
	userUsecase "try-on/internal/pkg/users/usecase"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SessionHandler struct {
	Sessions domain.SessionUsecase
	users    domain.UserUsecase
	cfg      *config.Session
}

//easyjson:json
type tokenResponse struct {
	Token string
}

func New(db *gorm.DB, cfg *config.Session) *SessionHandler {
	userRepo := userRepo.New(db)

	return &SessionHandler{
		Sessions: sessionUsecase.New(
			userRepo,
			cfg,
		),
		users: userUsecase.New(userRepo),
		cfg:   cfg,
	}
}

func (h *SessionHandler) Register(ctx *fiber.Ctx) error {
	var credentials domain.Credentials
	if err := ctx.BodyParser(&credentials); err != nil {
		return app_errors.New(err)
	}

	user, err := h.users.Create(credentials)
	switch {
	case err == nil:
		break

	case errors.Is(err, app_errors.ErrAlreadyExists):
		return fiber.ErrConflict

	default:
		return err
	}

	token, err := h.Sessions.IssueToken(user.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(tokenResponse{
		Token: token,
	})
}

func (h *SessionHandler) Login(ctx *fiber.Ctx) error {
	var credentials domain.Credentials
	if err := ctx.BodyParser(&credentials); err != nil {
		return app_errors.New(err)
	}

	session, err := h.Sessions.Login(credentials)
	switch {
	case err == nil:
		return ctx.JSON(tokenResponse{
			Token: session.ID,
		})

	case errors.Is(err, app_errors.ErrInvalidCredentials):
		return fiber.ErrForbidden

	default:
		return err
	}
}

func (h *SessionHandler) Renew(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return fiber.ErrUnauthorized
	}

	token, err := h.Sessions.IssueToken(session.UserID)
	if err != nil {
		return err
	}

	return ctx.JSON(tokenResponse{
		Token: token,
	})
}

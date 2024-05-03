package session

import (
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	userRepo "try-on/internal/pkg/repository/sqlc/users"
	sessionUsecase "try-on/internal/pkg/usecase/session"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
)

type SessionHandler struct {
	Sessions domain.SessionUsecase
	users    domain.UserRepository
	cfg      *config.Session
}

//easyjson:json
type loginResponse struct {
	Token    string
	UserName string
	UserID   utils.UUID
	Email    string
	Gender   domain.Gender
	Privacy  domain.Privacy
	Avatar   string
}

func New(db *pgxpool.Pool, cfg *config.Session) *SessionHandler {
	userRepo := userRepo.New(db)

	return &SessionHandler{
		Sessions: sessionUsecase.New(
			userRepo,
			cfg,
		),
		users: userRepo,
		cfg:   cfg,
	}
}

func (h *SessionHandler) Login(ctx *fiber.Ctx) error {
	var credentials domain.Credentials
	if err := easyjson.Unmarshal(ctx.Body(), &credentials); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	session, err := h.Sessions.Login(credentials)
	if err != nil {
		return app_errors.New(err)
	}

	user, err := h.users.GetByID(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(loginResponse{
		Token:    session.ID,
		UserID:   session.UserID,
		UserName: user.Name,
		Email:    user.Email,
		Gender:   user.Gender,
		Privacy:  user.Privacy,
		Avatar:   user.Avatar,
	})
}

func (h *SessionHandler) Renew(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	token, err := h.Sessions.IssueToken(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	user, err := h.users.GetByID(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(loginResponse{
		Token:    token,
		UserID:   session.UserID,
		UserName: user.Name,
		Email:    user.Email,
		Gender:   user.Gender,
		Privacy:  user.Privacy,
		Avatar:   user.Avatar,
	})
}

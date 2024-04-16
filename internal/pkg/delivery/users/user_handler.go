package users

import (
	"net/http"
	"strings"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/users"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	users domain.UserRepository
}

func New(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		users: users.New(db),
	}
}

func (h UserHandler) SearchUsers(ctx *fiber.Ctx) error {
	name := strings.TrimSpace(ctx.Query("name"))
	if name == "" {
		return app_errors.ResponseError{
			Code: http.StatusBadRequest,
			Msg:  "query param 'name' should be non-empty",
		}
	}

	users, err := h.users.SearchUsers(name)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(users)
}

func (h UserHandler) GetSubscriptions(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	users, err := h.users.GetSubscriptions(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(users)
}

package tags

import (
	"strconv"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/tags"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagsHandler struct {
	tags domain.TagRepository
}

func New(db *pgxpool.Pool) *TagsHandler {
	return &TagsHandler{
		tags: tags.New(db),
	}
}

const defaultTagLimit = 5

func (h TagsHandler) GetFavourite(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = defaultTagLimit
	}

	tags, err := h.tags.GetUserFavourite(session.UserID, limit)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(tags)
}

func (h TagsHandler) Get(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	from, _ := strconv.Atoi(ctx.Query("from"))

	tags, err := h.tags.Get(limit, from)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(tags)
}

package tags

import (
	"strconv"

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

func (h *TagsHandler) Get(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	from, _ := strconv.Atoi(ctx.Query("from"))

	tags, err := h.tags.Get(limit, from)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(tags)
}

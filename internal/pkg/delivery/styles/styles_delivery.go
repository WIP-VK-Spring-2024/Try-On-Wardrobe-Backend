package styles

import (
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/styles"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StylesHandler struct {
	styles domain.StylesRepository
}

func New(db *pgxpool.Pool) *StylesHandler {
	return &StylesHandler{
		styles: styles.New(db),
	}
}

func (h *StylesHandler) GetAll(ctx *fiber.Ctx) error {
	styles, err := h.styles.GetAll()
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(styles)
}

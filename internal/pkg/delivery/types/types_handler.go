package types

import (
	"database/sql"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/subtypes"
	"try-on/internal/pkg/repository/sqlc/types"

	"github.com/gofiber/fiber/v2"
)

type TypesHandler struct {
	types    domain.TypeRepository
	subtypes domain.SubtypeRepository
}

func New(db *sql.DB) *TypesHandler {
	return &TypesHandler{
		types:    types.New(db),
		subtypes: subtypes.New(db),
	}
}

func (h *TypesHandler) GetTypes(ctx *fiber.Ctx) error {
	types, err := h.types.GetAll()
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(types)
}

func (h *TypesHandler) GetSubtypes(ctx *fiber.Ctx) error {
	subtypes, err := h.subtypes.GetAll()
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(subtypes)
}

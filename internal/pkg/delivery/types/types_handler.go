package types

import (
	"database/sql"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"

	"github.com/gofiber/fiber/v2"
)

type TypesHandler struct {
	types    domain.TypeRepository
	subtypes domain.SubtypeRepository
}

func New(db *sql.DB) *TypesHandler {
	return &TypesHandler{
		types:    nil,
		subtypes: nil,
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

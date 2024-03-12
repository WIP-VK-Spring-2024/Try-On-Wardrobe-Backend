package delivery

import (
	"try-on/internal/pkg/app_errors"
	clothesRepo "try-on/internal/pkg/clothes/repository"
	clothesUsecase "try-on/internal/pkg/clothes/usecase"
	"try-on/internal/pkg/domain"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ClothesHandler struct {
	clothes domain.ClothesUsecase
}

func New(db *gorm.DB) *ClothesHandler {
	return &ClothesHandler{
		clothes: clothesUsecase.New(clothesRepo.New(db)),
	}
}

func (h *ClothesHandler) GetByID(ctx *fiber.Ctx) error {
	return app_errors.ErrUnimplemented
}

func (h *ClothesHandler) Upload(ctx *fiber.Ctx) error {
	return app_errors.ErrUnimplemented
}

func (h *ClothesHandler) GetByUser(ctx *fiber.Ctx) error {
	return app_errors.ErrUnimplemented
}

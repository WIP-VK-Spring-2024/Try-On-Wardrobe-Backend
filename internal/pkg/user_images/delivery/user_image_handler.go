package delivery

import (
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/file_manager/filesystem"
	"try-on/internal/pkg/user_images/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserImageHandler struct {
	userImages domain.UserImageRepository
	file       domain.FileManager
	model      domain.ClothesProcessingModel
	cfg        *config.Static
}

func New(
	db *gorm.DB,
	model domain.ClothesProcessingModel,
	cfg *config.Static,
) *UserImageHandler {
	return &UserImageHandler{
		userImages: repository.New(db),
		file:       filesystem.New(cfg.Dir),
		model:      model,
		cfg:        cfg,
	}
}

func (h *UserImageHandler) GetByID(ctx *fiber.Ctx) error {
	userImageID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	userImage, err := h.userImages.Get(userImageID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(userImage)
}

func (h *UserImageHandler) Delete(ctx *fiber.Ctx) error {
	userImageID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserImageIdInvalid
	}

	err = h.userImages.Delete(userImageID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *UserImageHandler) Upload(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var userImage domain.UserImage
	if err := ctx.BodyParser(&userImage); err != nil {
		return app_errors.ErrBadRequest
	}

	// TODO

	return ctx.SendString(common.EmptyJson)
}

func (h *UserImageHandler) GetByUser(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	userImages, err := h.userImages.GetByUser(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(userImages)
}

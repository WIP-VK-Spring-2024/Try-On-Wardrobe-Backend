package delivery

import (
	"fmt"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	clothesRepo "try-on/internal/pkg/clothes/repository"
	clothesUsecase "try-on/internal/pkg/clothes/usecase"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ClothesHandler struct {
	clothes domain.ClothesUsecase
	file    domain.FileManager
	model   domain.ClothesProcessingModel
}

func New(db *gorm.DB, file domain.FileManager) *ClothesHandler {
	return &ClothesHandler{
		clothes: clothesUsecase.New(clothesRepo.New(db)),
		file:    file,
	}
}

func (h *ClothesHandler) GetByID(ctx *fiber.Ctx) error {
	return app_errors.ErrUnimplemented
}

type clothesUploadArgs struct {
	Type string
}

func (h *ClothesHandler) Upload(ctx *fiber.Ctx) error {
	var args clothesUploadArgs
	if err := ctx.BodyParser(&args); err != nil {
		return app_errors.ErrBadRequest
	}

	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	clothes := &domain.Clothes{
		Type:   args.Type,
		Name:   args.Type,
		UserID: session.UserID,
	}

	err := h.clothes.Create(clothes)
	if err != nil {
		return app_errors.New(err)
	}

	fileHeader, err := ctx.FormFile("img")
	if err != nil {
		return app_errors.ErrBadRequest
	}

	file, err := fileHeader.Open()
	if err != nil {
		return app_errors.New(err)
	}
	defer file.Close()

	imagePath := fmt.Sprintf("%s/clothes/raw", session.UserID.String())

	err = h.file.Save(
		ctx.UserContext(),
		imagePath,
		clothes.ID.String(),
		file,
	)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.model.Process(domain.ClothesProcessingOpts{
		UserID:        session.UserID,
		ClothesID:     clothes.ID,
		ImagePath:     imagePath,
		CutBackground: true,
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(utils.EmptyJson)
}

func (h *ClothesHandler) GetByUser(ctx *fiber.Ctx) error {
	return app_errors.ErrUnimplemented
}

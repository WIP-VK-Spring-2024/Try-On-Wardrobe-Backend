package delivery

import (
	"fmt"
	"net/http"
	"os"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	clothesRepo "try-on/internal/pkg/clothes/repository"
	clothesUsecase "try-on/internal/pkg/clothes/usecase"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClothesHandler struct {
	clothes domain.ClothesUsecase
	file    domain.FileManager
	model   domain.ClothesProcessingModel
}

func New(db *gorm.DB, file domain.FileManager, model domain.ClothesProcessingModel) *ClothesHandler {
	return &ClothesHandler{
		clothes: clothesUsecase.New(clothesRepo.New(db)),
		file:    file,
		model:   model,
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
		middleware.LogError(ctx, err)
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
		middleware.LogError(ctx, err)
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

	err = h.model.Process(ctx.UserContext(), domain.ClothesProcessingOpts{
		UserID:    session.UserID,
		ImageID:   clothes.ID,
		FileName:  imagePath,
		ImageType: domain.ImageTypeCloth,
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *ClothesHandler) GetByUser(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		middleware.LogError(ctx, err)
		return &app_errors.ErrorMsg{
			Code: http.StatusBadRequest,
			Msg:  "userID should be a valid uuid",
		}
	}

	var filters domain.ClothesFilters

	if err := ctx.QueryParser(&filters); err != nil {
		middleware.LogError(ctx, err)
		return &app_errors.ErrorMsg{
			Code: http.StatusBadRequest,
			Msg:  "invalid filters passed",
		}
	}

	clothes, err := h.clothes.GetByUser(userID, &filters)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(clothes) // TODO: Make normal
}

func (h *ClothesHandler) ListenTryOnResults() {
}

func (h *ClothesHandler) TryOn(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		middleware.LogError(ctx, err)
		return &app_errors.ErrorMsg{
			Code: http.StatusBadRequest,
			Msg:  "userID should be a valid uuid",
		}
	}

	clothingID, err := uuid.Parse(ctx.Query("clothing_id"))
	if err != nil {
		middleware.LogError(ctx, err)
		return &app_errors.ErrorMsg{
			Code: http.StatusBadRequest,
			Msg:  "clothingID should be a valid uuid",
		}
	}

	clothing, err := h.clothes.Get(clothingID)
	if err != nil {
		return app_errors.New(err)
	}

	curPath, err := os.Getwd()
	if err != nil {
		return app_errors.New(err)
	}

	personFileName := "person.jpg"

	err = h.model.TryOn(ctx.UserContext(), domain.TryOnOpts{
		UserID:          userID,
		ClothesFileName: clothing.Image,
		ClothesFilePath: curPath + "/stubs/clothes/" + clothing.Image,
		PersonFileName:  personFileName,
		PersonFilePath:  curPath + "/stubs/people/" + personFileName,
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

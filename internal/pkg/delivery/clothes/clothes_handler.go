package clothes

import (
	"strconv"
	"strings"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/file_manager/filesystem"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClothesHandler struct {
	clothes domain.ClothesUsecase
	file    domain.FileManager
	model   domain.ClothesProcessingModel
	cfg     *config.Static
}

func New(
	clothes domain.ClothesUsecase,
	model domain.ClothesProcessingModel,
	cfg *config.Static,
) *ClothesHandler {
	return &ClothesHandler{
		clothes: clothes,
		file:    filesystem.New(cfg.Dir),
		model:   model,
		cfg:     cfg,
	}
}

func (h *ClothesHandler) GetByID(ctx *fiber.Ctx) error {
	clothesID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	clothes, err := h.clothes.Get(clothesID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(clothes)
}

func (h *ClothesHandler) Delete(ctx *fiber.Ctx) error {
	clothesID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	err = h.clothes.Delete(clothesID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *ClothesHandler) Update(ctx *fiber.Ctx) error {
	clothesID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	clothes := &domain.Clothes{ID: clothesID}
	if err := ctx.BodyParser(clothes); err != nil {
		middleware.LogError(ctx, err)
		return app_errors.ErrBadRequest
	}

	err = h.clothes.Update(clothes)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *ClothesHandler) Upload(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var clothes domain.Clothes
	if err := ctx.BodyParser(&clothes); err != nil {
		middleware.LogError(ctx, err)
		return app_errors.ErrBadRequest
	}

	if len(clothes.Tags) == 1 {
		clothes.Tags = strings.Split(clothes.Tags[0], ",")
	}

	clothes.UserID = session.UserID

	err := h.clothes.Create(&clothes)
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

	err = h.file.Save(
		ctx.UserContext(),
		h.cfg.Clothes,
		clothes.ID.String(),
		file,
	)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.model.Process(ctx.UserContext(), domain.ClothesProcessingOpts{
		UserID:    session.UserID,
		ImageID:   clothes.ID,
		FileName:  clothes.ID.String(),
		ImageType: domain.ImageTypeCloth,
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *ClothesHandler) getClothes(userID uuid.UUID, ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	clothes, err := h.clothes.GetByUser(userID, limit)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(clothes)
}

func (h *ClothesHandler) GetByUser(ctx *fiber.Ctx) error {
	userID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserIdInvalid
	}

	return h.getClothes(userID, ctx)
}

func (h *ClothesHandler) GetOwn(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	return h.getClothes(session.UserID, ctx)
}

package user_images

import (
	"errors"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/usecase/user_images"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserImageHandler struct {
	userImages domain.UserImageUsecase
	file       domain.FileManager
	cfg        *config.Static
}

func New(
	db *pgxpool.Pool,
	fileManager domain.FileManager,
	cfg *config.Static,
) *UserImageHandler {
	return &UserImageHandler{
		userImages: user_images.New(db),
		file:       fileManager,
		cfg:        cfg,
	}
}

//easyjson:json
type imageUploadedResponse struct {
	Uuid  utils.UUID
	Image string
}

func (h *UserImageHandler) GetByID(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	userImageID, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserImageIdInvalid
	}

	userImage, err := h.userImages.Get(userImageID)
	if err != nil {
		return app_errors.New(err)
	}

	if userImage.UserID != session.UserID {
		return app_errors.ErrNotOwner
	}

	return ctx.JSON(userImage)
}

func (h *UserImageHandler) Delete(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	userImageID, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserImageIdInvalid
	}

	userImage, err := h.userImages.Get(userImageID)
	if err != nil {
		return app_errors.New(err)
	}

	if userImage.UserID != session.UserID {
		return app_errors.ErrNotOwner
	}

	err = h.userImages.Delete(userImageID)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.file.Delete(ctx.UserContext(), h.cfg.Type, userImageID.String())
	if err != nil {
		middleware.LogWarning(ctx, errors.Join(err, errors.New("user_image image deletion error")))
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *UserImageHandler) Upload(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	userImage := domain.UserImage{
		UserID: session.UserID,
		Image:  h.cfg.FullBody,
	}

	fileHeader, err := ctx.FormFile("img")
	if err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	file, err := fileHeader.Open()
	if err != nil {
		return app_errors.New(err)
	}
	defer file.Close()

	err = h.userImages.Create(&userImage)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.file.Save(
		ctx.UserContext(),
		h.cfg.FullBody,
		userImage.ID.String(),
		file,
	)
	if err != nil {
		deleteErr := h.userImages.Delete(userImage.ID)
		middleware.LogError(ctx, deleteErr)
		return app_errors.New(err)
	}

	return ctx.JSON(&imageUploadedResponse{
		Uuid:  userImage.ID,
		Image: userImage.Image,
	})
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

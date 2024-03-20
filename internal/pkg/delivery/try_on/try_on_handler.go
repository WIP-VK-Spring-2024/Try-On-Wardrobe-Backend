package try_on

import (
	"os"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/try_on"
	"try-on/internal/pkg/repository/sqlc/user_images"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TryOnHandler struct {
	model      domain.ClothesProcessingModel
	clothes    domain.ClothesUsecase
	userImages domain.UserImageRepository
	results    domain.TryOnResultRepository
	logger     *zap.SugaredLogger
	cfg        *config.Static
}

func New(
	db *pgxpool.Pool,
	model domain.ClothesProcessingModel,
	clothes domain.ClothesUsecase,
	logger *zap.SugaredLogger,
	cfg *config.Static,
) *TryOnHandler {
	return &TryOnHandler{
		model:      model,
		clothes:    clothes,
		userImages: user_images.New(db),
		results:    try_on.New(db),
		logger:     logger,
		cfg:        cfg,
	}
}

func (h *TryOnHandler) ListenTryOnResults() {
	go func() {
		err := h.model.GetTryOnResults(h.logger, h.handleResult)
		if err != nil {
			h.logger.Errorw(err.Error())
		}
	}()
}

func (h *TryOnHandler) handleResult(resp *domain.TryOnResponse) domain.Result {
	tryOnRes := &domain.TryOnResult{
		UserImageID: resp.UserImageID,
		ClothesID:   resp.ClothesID,
		Image:       resp.ResFilePath,
	}

	err := h.results.Create(tryOnRes)
	if err != nil {
		h.logger.Errorw(err.Error())
		return domain.ResultRetry
	}

	return domain.ResultOk
}

//easyjson:json
type tryOnRequest struct {
	ClothesID   uuid.UUID
	UserImageID uuid.UUID
}

func (h *TryOnHandler) TryOn(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var req tryOnRequest
	if err := ctx.BodyParser(&req); err != nil {
		middleware.LogError(ctx, err)
		return app_errors.ErrClothesIdInvalid
	}

	clothes, err := h.clothes.Get(req.ClothesID)
	if err != nil {
		return app_errors.New(err)
	}

	userImage, err := h.userImages.Get(req.UserImageID)
	if err != nil {
		return app_errors.New(err)
	}

	curPath, err := os.Getwd()
	if err != nil {
		return app_errors.New(err)
	}

	err = h.model.TryOn(ctx.UserContext(), domain.TryOnOpts{
		UserID:          session.UserID,
		ClothesID:       req.ClothesID,
		ClothesFileName: clothes.Image,
		ClothesFilePath: curPath + "/stubs/clothes/" + clothes.Image,
		PersonFileName:  userImage.Image,
		PersonFilePath:  curPath + "/stubs/people/" + userImage.Image,
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (c *TryOnHandler) GetTryOnResult(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	result, err := c.results.GetLast(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(result)
}

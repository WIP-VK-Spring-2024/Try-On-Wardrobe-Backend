package clothes

import (
	"context"
	"strconv"
	"strings"

	"try-on/internal/generated/proto/centrifugo"
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ClothesHandler struct {
	clothes domain.ClothesUsecase
	file    domain.FileManager
	model   domain.ClothesProcessingModel
	cfg     *config.Static

	logger     *zap.SugaredLogger
	centrifugo centrifugo.CentrifugoApiClient
}

func New(
	clothes domain.ClothesUsecase,
	model domain.ClothesProcessingModel,
	fileManager domain.FileManager,
	cfg *config.Static,
) *ClothesHandler {
	return &ClothesHandler{
		clothes: clothes,
		file:    fileManager,
		model:   model,
		cfg:     cfg,
	}
}

func (h *ClothesHandler) GetByID(ctx *fiber.Ctx) error {
	clothesID, err := utils.ParseUUID(ctx.Params("id"))
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
	clothesID, err := utils.ParseUUID(ctx.Params("id"))
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
	clothesID, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	clothes := &domain.Clothes{}
	if err := ctx.BodyParser(clothes); err != nil {
		middleware.LogError(ctx, err)
		return app_errors.ErrBadRequest
	}
	clothes.ID = clothesID

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

	var clothes domain.Clothes
	if err := ctx.BodyParser(&clothes); err != nil {
		middleware.LogError(ctx, err)
		return app_errors.ErrBadRequest
	}

	if len(clothes.Tags) == 1 {
		clothes.Tags = strings.Split(clothes.Tags[0], ",")
	}

	clothes.UserID = session.UserID

	err = h.clothes.Create(&clothes)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.file.Save(
		ctx.UserContext(),
		h.cfg.Clothes,
		clothes.ID.String(),
		file,
	)
	if err != nil {
		deleteErr := h.clothes.Delete(clothes.ID)
		middleware.LogError(ctx, deleteErr)
		return app_errors.New(err)
	}

	err = h.model.Process(ctx.UserContext(), domain.ClothesProcessingRequest{
		UserID:     session.UserID,
		ClothesID:  clothes.ID,
		ClothesDir: h.cfg.Clothes,
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *ClothesHandler) getClothes(userID utils.UUID, ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	clothes, err := h.clothes.GetByUser(userID, limit)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(clothes)
}

func (h *ClothesHandler) GetByUser(ctx *fiber.Ctx) error {
	userID, err := utils.ParseUUID(ctx.Params("id"))
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

func (h *ClothesHandler) ListenProcessingResults(cfg *config.Centrifugo) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Error(err)
			}
		}()

		err := h.model.GetProcessingResults(h.logger, h.handleQueueResponse(cfg))
		if err != nil {
			h.logger.Errorw(err.Error())
		}
	}()
}

func (h *ClothesHandler) handleQueueResponse(cfg *config.Centrifugo) func(resp *domain.ClothesProcessingResponse) domain.Result {
	return func(resp *domain.ClothesProcessingResponse) domain.Result {
		userChannel := cfg.ProcessingChannel + resp.UserID.String()
		h.logger.Infow("centrifugo", "channel", userChannel)

		centrifugoResp, err := h.centrifugo.Publish(
			context.Background(),
			&centrifugo.PublishRequest{
				Channel: userChannel,
				Data:    []byte(common.EmptyJson),
			},
		)

		switch {
		case err != nil:
			h.logger.Errorw(err.Error())
			return domain.ResultRetry
		case centrifugoResp.Error != nil:
			h.logger.Errorw(centrifugoResp.Error.Message)
		}

		return domain.ResultOk
	}
}

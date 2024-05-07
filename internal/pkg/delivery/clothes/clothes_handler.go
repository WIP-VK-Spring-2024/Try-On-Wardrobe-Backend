package clothes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/validate"

	"github.com/gofiber/fiber/v2"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type ClothesHandler struct {
	clothes domain.ClothesUsecase
	tags    domain.TagUsecase
	file    domain.FileManager
	model   domain.ClothesProcessingModel
	cfg     *config.Static

	logger    *zap.SugaredLogger
	publisher domain.ChannelPublisher[easyjson.Marshaler]
}

func New(
	clothes domain.ClothesUsecase,
	tags domain.TagUsecase,
	model domain.ClothesProcessingModel,
	fileManager domain.FileManager,
	cfg *config.Static,
	logger *zap.SugaredLogger,
	publisher domain.ChannelPublisher[easyjson.Marshaler],
) *ClothesHandler {
	return &ClothesHandler{
		clothes:   clothes,
		tags:      tags,
		file:      fileManager,
		model:     model,
		cfg:       cfg,
		logger:    logger,
		publisher: publisher,
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
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	clothesID, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	clothes, err := h.clothes.Get(clothesID)
	if err != nil {
		return app_errors.New(err)
	}

	if clothes.UserID != session.UserID {
		return app_errors.ErrNotOwner
	}

	err = h.clothes.Delete(session.UserID, clothesID)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.file.Delete(ctx.UserContext(), h.cfg.Clothes, clothesID.String())
	if err != nil {
		middleware.LogWarning(ctx, err, "clothes_id", clothesID)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *ClothesHandler) Update(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	clothesID, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	clothesUpdate := &domain.Clothes{}
	if err := easyjson.Unmarshal(ctx.Body(), clothesUpdate); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	err = validate.Struct(clothesUpdate)
	if err != nil {
		return app_errors.ValidationError(err)
	}

	clothesUpdate.ID = clothesID
	clothesUpdate.UserID = session.UserID

	clothes, err := h.clothes.Get(clothesID)
	if err != nil {
		return app_errors.New(err)
	}

	if clothes.UserID != session.UserID {
		return app_errors.New(app_errors.ErrNotOwner)
	}

	if len(clothesUpdate.Tags) > 0 {
		err = h.tags.Create(clothesUpdate.Tags)
		if err != nil {
			return app_errors.New(err)
		}
	}

	err = h.clothes.Update(clothesUpdate)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

//easyjson:json
type uploadResponse struct {
	Uuid  utils.UUID
	Msg   string
	Image string
}

func (h *ClothesHandler) Upload(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	fileHeader, err := ctx.FormFile("img")
	if err != nil {
		middleware.LogWarning(ctx, err)
		middleware.LogError(ctx, err)
		return app_errors.ErrBadRequest
	}

	file, err := fileHeader.Open()
	if err != nil {
		return app_errors.New(err)
	}
	defer file.Close()

	var clothes domain.Clothes

	clothes.UserID = session.UserID
	clothes.Image = h.cfg.Clothes

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
		deleteErr := h.clothes.Delete(session.UserID, clothes.ID)
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

	return ctx.JSON(&uploadResponse{
		Uuid:  clothes.ID,
		Msg:   domain.ClothesStatusCreated,
		Image: h.cfg.Clothes + "/" + clothes.ID.String(),
	})
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

//easyjson:json
type processingResponse struct {
	uploadResponse
	Tryonable      bool
	Classification domain.ClothesClassificationResponse
}

func (h *ClothesHandler) handleQueueResponse(cfg *config.Centrifugo) func(resp *domain.ClothesProcessingResponse) domain.Result {
	ctx := middleware.WithLogger(context.Background(), h.logger)

	return func(resp *domain.ClothesProcessingResponse) domain.Result {
		userChannel := cfg.ProcessingChannel + resp.UserID.String()

		fmt.Printf("Resp in handler after post processing: %+v\n", resp)

		if !utils.HttpOk(resp.StatusCode) {
			fmt.Println("Resp code is", resp.StatusCode)

			h.publisher.Publish(
				ctx,
				userChannel,
				&app_errors.ResponseError{
					Code: http.StatusInternalServerError,
					Msg:  resp.Message,
				},
			)

			return domain.ResultOk
		}

		cutImageUrl := h.cfg.Cut + "/" + resp.ClothesID.String()
		err := h.clothes.SetImage(resp.ClothesID, cutImageUrl)
		if err != nil {
			h.logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		clothesUpdate := domain.Clothes{
			Model: domain.Model{
				ID: resp.ClothesID,
			},
			Seasons:   resp.Classification.Seasons,
			Tags:      resp.Classification.Tags,
			StyleID:   resp.Classification.Style,
			TypeID:    resp.Classification.Type,
			SubtypeID: resp.Classification.Subtype,
		}

		err = h.clothes.Update(&clothesUpdate)
		if err != nil {
			h.logger.Errorw(err.Error())
		}

		payload := &processingResponse{
			uploadResponse: uploadResponse{
				Uuid:  resp.ClothesID,
				Msg:   domain.ClothesStatusProcessed,
				Image: cutImageUrl,
			},
			Tryonable:      resp.Tryonable,
			Classification: resp.Classification,
		}

		fmt.Printf("Sending to centrifugo: %+v\n", payload)

		h.publisher.Publish(
			ctx,
			userChannel,
			payload,
		)

		return domain.ResultOk
	}
}

package try_on

import (
	"context"
	"net/http"

	"try-on/internal/generated/proto/centrifugo"
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/try_on"
	"try-on/internal/pkg/repository/sqlc/user_images"
	"try-on/internal/pkg/usecase/outfits"
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/translate"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type TryOnHandler struct {
	model domain.ClothesProcessingModel

	clothes    domain.ClothesUsecase
	outfits    domain.OutfitUsecase
	userImages domain.UserImageRepository
	results    domain.TryOnResultRepository

	centrifugo centrifugo.CentrifugoApiClient

	logger *zap.SugaredLogger
}

func New(
	db *pgxpool.Pool,
	model domain.ClothesProcessingModel,
	clothes domain.ClothesUsecase,
	logger *zap.SugaredLogger,
	centrifugoConn grpc.ClientConnInterface,
) *TryOnHandler {
	return &TryOnHandler{
		model:      model,
		clothes:    clothes,
		userImages: user_images.New(db),
		results:    try_on.New(db),
		outfits:    outfits.NewWithSqlcRepo(db),
		logger:     logger,
		centrifugo: centrifugo.NewCentrifugoApiClient(centrifugoConn),
	}
}

func (h *TryOnHandler) ListenTryOnResults(cfg *config.Centrifugo) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Error(err)
			}
		}()

		err := h.model.GetTryOnResults(h.logger, h.handleQueueResponse(cfg))
		if err != nil {
			h.logger.Errorw(err.Error())
		}
	}()
}

func (h *TryOnHandler) handleQueueResponse(cfg *config.Centrifugo) func(resp *domain.TryOnResponse) domain.Result {
	return func(resp *domain.TryOnResponse) domain.Result {
		tryOnRes := &domain.TryOnResult{
			UserImageID: resp.UserImageID,
			ClothesID:   resp.ClothesID,
			Image:       "/" + resp.TryOnResultDir + "/" + resp.TryOnResultID,
		}

		handleResult := domain.ResultOk

		err := h.results.Create(tryOnRes)
		if err != nil {
			h.logger.Errorw(err.Error())
			handleResult = domain.ResultDiscard
		}

		var payload []byte
		if handleResult == domain.ResultDiscard {
			payload, _ = easyjson.Marshal(app_errors.ResponseError{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			})
		} else {
			payload, _ = easyjson.Marshal(tryOnRes)
		}

		userChannel := cfg.TryOnChannel + resp.UserID.String()
		h.logger.Infow("centrifugo", "channel", userChannel, "payload", payload)

		centrifugoResp, err := h.centrifugo.Publish(
			context.Background(),
			&centrifugo.PublishRequest{
				Channel: userChannel,
				Data:    payload,
			},
		)

		switch {
		case err != nil:
			h.logger.Errorw(err.Error())
		case centrifugoResp.Error != nil:
			h.logger.Errorw(centrifugoResp.Error.Message)
		}

		return domain.ResultOk
	}
}

//easyjson:json
type tryOnRequest struct {
	ClothesID   utils.UUID
	UserImageID utils.UUID
}

func (h *TryOnHandler) TryOn(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var req tryOnRequest
	err := easyjson.Unmarshal(ctx.Body(), &req)
	if err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	clothes, err := h.clothes.Get(req.ClothesID)
	if err != nil {
		return app_errors.New(err)
	}

	_, err = h.userImages.Get(req.UserImageID)
	if err != nil {
		return app_errors.New(err)
	}

	cfg := middleware.Config(ctx)

	err = h.model.TryOn(ctx.UserContext(), domain.TryOnRequest{
		UserID:       session.UserID,
		UserImageID:  req.UserImageID,
		UserImageDir: cfg.Static.FullBody,
		ClothesDir:   cfg.Static.Cut,
		Clothes: map[utils.UUID]string{
			req.ClothesID: translate.ClothesTypeToTryOnCategory(clothes.Type),
		},
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

//easyjson:json
type tryOnOutfitRequest struct {
	OutfitID    utils.UUID
	UserImageID utils.UUID
}

func (h *TryOnHandler) TryOnOutfit(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var req tryOnOutfitRequest
	err := easyjson.Unmarshal(ctx.Body(), &req)
	if err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	clothesInfo, err := h.outfits.GetClothesInfo(req.OutfitID)
	if err != nil {
		return app_errors.New(err)
	}

	_, err = h.userImages.Get(req.UserImageID)
	if err != nil {
		return app_errors.New(err)
	}

	cfg := middleware.Config(ctx)

	err = h.model.TryOn(ctx.UserContext(), domain.TryOnRequest{
		UserID:       session.UserID,
		UserImageID:  req.UserImageID,
		UserImageDir: cfg.Static.FullBody,
		ClothesDir:   cfg.Static.Cut,
		Clothes:      clothesInfo,
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *TryOnHandler) GetTryOnResult(ctx *fiber.Ctx) error {
	userImageID, err := utils.ParseUUID(ctx.Query("photo_id"))
	if err != nil {
		return app_errors.ErrUserImageIdInvalid
	}

	clothesID, err := utils.ParseUUID(ctx.Query("clothes_id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	result, err := h.results.Get(userImageID, clothesID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(result)
}

func (h *TryOnHandler) GetByUser(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	results, err := h.results.GetByUser(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(results)
}

//easyjson:json
type ratingRequest struct {
	Rating int
}

func (h *TryOnHandler) Rate(ctx *fiber.Ctx) error {
	tryOnResultId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrTryOnIdInvalid
	}

	var req ratingRequest
	if err := easyjson.Unmarshal(ctx.Body(), &req); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	err = h.results.Rate(tryOnResultId, req.Rating)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

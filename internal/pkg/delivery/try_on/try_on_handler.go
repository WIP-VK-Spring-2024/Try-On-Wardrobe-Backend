package try_on

import (
	"context"
	"fmt"
	"net/http"

	"try-on/internal/generated/proto/centrifugo"
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/try_on"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type TryOnHandler struct {
	tryOnModel domain.TryOnUsecase
	results    domain.TryOnResultRepository

	centrifugo centrifugo.CentrifugoApiClient
	logger     *zap.SugaredLogger
}

func New(
	db *pgxpool.Pool,
	tryOnModel domain.TryOnUsecase,
	logger *zap.SugaredLogger,
	centrifugoConn grpc.ClientConnInterface,
) *TryOnHandler {
	return &TryOnHandler{
		tryOnModel: tryOnModel,
		results:    try_on.New(db),
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

		err := h.tryOnModel.GetTryOnResults(h.logger, h.handleQueueResponse(cfg))
		if err != nil {
			h.logger.Errorw(err.Error())
		}
	}()
}

func (h *TryOnHandler) centrifugoPublish(message easyjson.Marshaler, channel string) domain.Result {
	payload, _ := easyjson.Marshal(message)

	h.logger.Infow("centrifugo", "channel", channel, "payload", string(payload))

	centrifugoResp, err := h.centrifugo.Publish(
		context.Background(),
		&centrifugo.PublishRequest{
			Channel: channel,
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

func (h *TryOnHandler) handleQueueResponse(cfg *config.Centrifugo) func(resp *domain.TryOnResponse) domain.Result {
	return func(resp *domain.TryOnResponse) domain.Result {
		userChannel := cfg.TryOnChannel + resp.UserID.String()

		if resp.StatusCode != http.StatusOK {
			return h.centrifugoPublish(app_errors.ResponseError{
				Code: http.StatusInternalServerError,
				Msg:  resp.Message,
			}, userChannel)
		}

		fmt.Println("Got clothes from rabbit try on", resp.Clothes)

		clothesIds := make([]utils.UUID, 0, len(resp.Clothes))
		for _, clothes := range resp.Clothes {
			clothesIds = append(clothesIds, clothes.ClothesID)
		}

		fmt.Println("Clothes IDs from try on", clothesIds)

		tryOnRes := &domain.TryOnResult{
			UserImageID: resp.UserImageID,
			ClothesID:   clothesIds,
			Image:       "/" + resp.TryOnDir + "/" + resp.TryOnID,
		}

		handleResult := domain.ResultOk

		err := h.results.Create(tryOnRes)
		if err != nil {
			h.logger.Errorw(err.Error())
			handleResult = domain.ResultDiscard
		}

		var payload easyjson.Marshaler
		if handleResult == domain.ResultDiscard {
			payload = app_errors.ResponseError{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}
		} else {
			payload = tryOnRes
		}

		return h.centrifugoPublish(payload, userChannel)
	}
}

//easyjson:json
type tryOnRequest struct {
	ClothesID   []utils.UUID
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

	cfg := middleware.Config(ctx)

	err = h.tryOnModel.TryOn(ctx.UserContext(), req.ClothesID, domain.TryOnOpts{
		UserID:       session.UserID,
		UserImageID:  req.UserImageID,
		UserImageDir: cfg.Static.FullBody,
		ClothesDir:   cfg.Static.Cut,
	})
	if err != nil {
		middleware.LogWarning(ctx, err)
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

	cfg := middleware.Config(ctx)

	err = h.tryOnModel.TryOnOutfit(ctx.UserContext(), req.OutfitID, domain.TryOnOpts{
		UserID:       session.UserID,
		UserImageID:  req.UserImageID,
		UserImageDir: cfg.Static.FullBody,
		ClothesDir:   cfg.Static.Cut,
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *TryOnHandler) GetTryOnResult(ctx *fiber.Ctx) error {
	id, err := utils.ParseUUID(ctx.Query("photo_id"))
	if err != nil {
		return app_errors.ErrTryOnIdInvalid
	}

	result, err := h.results.Get(id)
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

	if req.Rating > 1 {
		req.Rating = 1
	}
	if req.Rating < -1 {
		req.Rating = -1
	}

	err = h.results.Rate(tryOnResultId, req.Rating)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

package try_on

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"try-on/internal/generated/proto/centrifugo"
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/try_on"
	"try-on/internal/pkg/repository/sqlc/user_images"
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/translate"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type TryOnHandler struct {
	model      domain.ClothesProcessingModel
	clothes    domain.ClothesUsecase
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

func (h *TryOnHandler) handleQueueResponse(cfg *config.Centrifugo) func(resp interface{}) domain.Result {
	return func(response interface{}) domain.Result {
		resp := response.(*domain.TryOnResponse) // BAD CODE

		tryOnRes := &domain.TryOnResult{
			UserImageID: resp.UserImageID,
			ClothesID:   resp.ClothesID,
			Image:       "/" + resp.TryOnResultDir + "/" + resp.TryOnResultID,
		}

		fmt.Println("Path to image", tryOnRes.Image)

		handleResult := domain.ResultOk

		err := h.results.Create(tryOnRes)
		switch {
		case err == nil:
			break
		case errors.Is(err, app_errors.ErrAlreadyExists) || errors.Is(err, app_errors.ErrNoRelatedEntity):
			h.logger.Errorw(err.Error())
			handleResult = domain.ResultDiscard
		default:
			h.logger.Errorw(err.Error())
			time.Sleep(time.Second) // BAD CODE
			return domain.ResultRetry
		}

		var payload []byte
		if handleResult == domain.ResultDiscard {
			payload, _ = easyjson.Marshal(app_errors.ResponseError{
				Code: http.StatusConflict,
				Msg:  err.Error(),
			})
		} else {
			payload, _ = easyjson.Marshal(tryOnRes)
		}

		userChannel := cfg.TryOnChannel + resp.UserID.String()
		h.logger.Infow("centrifugo", "channel", userChannel)

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
		middleware.LogError(ctx, err)
		return app_errors.ErrBadRequest
	}

	fmt.Printf("%+v\n", req)

	clothes, err := h.clothes.Get(req.ClothesID)
	if err != nil {
		middleware.LogError(ctx, err)
		return app_errors.New(err)
	}

	_, err = h.userImages.Get(req.UserImageID)
	if err != nil {
		middleware.LogError(ctx, err)
		return app_errors.New(err)
	}

	cfg := middleware.Config(ctx)

	err = h.model.TryOn(ctx.UserContext(), domain.TryOnRequest{
		UserID:       session.UserID,
		ClothesID:    req.ClothesID,
		UserImageID:  req.UserImageID,
		UserImageDir: cfg.Static.FullBody,
		ClothesDir:   cfg.Static.Cut,
		Category:     translate.ClothesTypeToTryOnCategory(clothes.Type),
	})
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (c *TryOnHandler) GetTryOnResult(ctx *fiber.Ctx) error {
	userImageID, err := utils.ParseUUID(ctx.Query("photo_id"))
	if err != nil {
		return app_errors.ErrUserImageIdInvalid
	}

	clothesID, err := utils.ParseUUID(ctx.Query("clothes_id"))
	if err != nil {
		return app_errors.ErrClothesIdInvalid
	}

	result, err := c.results.Get(userImageID, clothesID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(result)
}

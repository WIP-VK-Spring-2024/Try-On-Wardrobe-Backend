package try_on

import (
	"context"
	"errors"
	"net/http"
	"os"

	"try-on/internal/generated/proto/centrifugo"
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/try_on"
	"try-on/internal/pkg/repository/sqlc/user_images"
	"try-on/internal/pkg/utils"

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
	logger     *zap.SugaredLogger
	cfg        *config.Centrifugo
}

func New(
	db *pgxpool.Pool,
	model domain.ClothesProcessingModel,
	clothes domain.ClothesUsecase,
	logger *zap.SugaredLogger,
	centrifugoConn grpc.ClientConnInterface,
	cfg *config.Centrifugo,
) *TryOnHandler {
	return &TryOnHandler{
		model:      model,
		clothes:    clothes,
		userImages: user_images.New(db),
		results:    try_on.New(db),
		logger:     logger,
		centrifugo: centrifugo.NewCentrifugoApiClient(centrifugoConn),
		cfg:        cfg,
	}
}

func (h *TryOnHandler) ListenTryOnResults() {
	go func() {
		err := h.model.GetTryOnResults(h.logger, h.handleQueueResponse)
		if err != nil {
			h.logger.Errorw(err.Error())
		}
	}()
}

func (h *TryOnHandler) handleQueueResponse(resp *domain.TryOnResponse) domain.Result {
	tryOnRes := &domain.TryOnResult{
		UserImageID: resp.UserImageID,
		ClothesID:   resp.ClothesID,
		Image:       resp.ResFilePath,
	}

	handleResult := domain.ResultOk

	err := h.results.Create(tryOnRes)
	switch {
	case err == nil:
		break
	case errors.Is(err, app_errors.ErrAlreadyExists):
		h.logger.Errorw(err.Error())
		handleResult = domain.ResultDiscard
	default:
		h.logger.Errorw(err.Error())
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

	userChannel := h.cfg.TryOnChannel + resp.UserID.String()
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

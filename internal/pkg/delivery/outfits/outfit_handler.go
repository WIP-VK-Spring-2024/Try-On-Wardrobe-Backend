package outfits

import (
	"context"
	"fmt"
	"time"

	"try-on/internal/generated/proto/centrifugo"
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	outfitRepo "try-on/internal/pkg/repository/sqlc/outfits"
	outfitUsecase "try-on/internal/pkg/usecase/outfits"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	easyjson "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type OutfitHandler struct {
	outfits   domain.OutfitUsecase
	generator domain.OutfitGenerator

	file domain.FileManager
	cfg  *config.Static

	logger     *zap.SugaredLogger
	centrifugo centrifugo.CentrifugoApiClient
}

func New(
	db *pgxpool.Pool,
	generator domain.OutfitGenerator,
	file domain.FileManager,
	cfg *config.Static,
	logger *zap.SugaredLogger,
	centrifugoConn grpc.ClientConnInterface,
) *OutfitHandler {
	return &OutfitHandler{
		outfits:    outfitUsecase.New(outfitRepo.New(db)),
		generator:  generator,
		file:       file,
		cfg:        cfg,
		logger:     logger,
		centrifugo: centrifugo.NewCentrifugoApiClient(centrifugoConn),
	}
}

func (h *OutfitHandler) GetById(ctx *fiber.Ctx) error {
	id, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrOutfitIdInvalid
	}

	outfit, err := h.outfits.GetById(id)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(outfit)
}

//easyjson:json
type getOutfitsParams struct {
	Limit int
	Since time.Time
}

func (h *OutfitHandler) Get(ctx *fiber.Ctx) error {
	var params getOutfitsParams

	if err := ctx.QueryParser(&params); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	outfits, err := h.outfits.Get(params.Since, params.Limit)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(outfits)
}

func (h *OutfitHandler) GetByUser(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	outfits, err := h.outfits.GetByUser(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(outfits)
}

//easyjson:json
type createdResponse struct {
	Uuid  utils.UUID
	Image string
}

func (h *OutfitHandler) Create(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
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

	var outfit domain.Outfit

	transforms := ctx.FormValue("transforms")

	if err := easyjson.Unmarshal([]byte(transforms), &outfit.Transforms); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}
	outfit.UserID = session.UserID
	outfit.Image = h.cfg.Outfits

	err = h.outfits.Create(&outfit)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.file.Save(ctx.UserContext(), h.cfg.Outfits, outfit.ID.String(), file)
	if err != nil {
		if deleteErr := h.outfits.Delete(session.UserID, outfit.ID); deleteErr != nil {
			middleware.LogError(ctx, err)
		}
		return app_errors.New(err)
	}

	return ctx.JSON(&createdResponse{
		Uuid:  outfit.ID,
		Image: outfit.Image,
	})
}

func (h *OutfitHandler) Update(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	id, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrOutfitIdInvalid
	}

	var outfit domain.Outfit

	if err := ctx.BodyParser(&outfit); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}
	outfit.UserID = session.UserID
	outfit.ID = id

	transforms := ctx.FormValue("transforms")

	if err := easyjson.Unmarshal([]byte(transforms), &outfit.Transforms); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	err = h.outfits.Update(&outfit)
	if err != nil {
		return app_errors.New(err)
	}

	fileHeader, err := ctx.FormFile("img")
	switch {
	case err == nil && fileHeader != nil:
		break
	case err != fasthttp.ErrMissingFile:
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	default:
		return ctx.SendString(common.EmptyJson)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return app_errors.New(err)
	}
	defer file.Close()

	err = h.file.Save(ctx.UserContext(), h.cfg.Outfits, outfit.ID.String(), file)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *OutfitHandler) Delete(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	outfitId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrOutfitIdInvalid
	}

	err = h.outfits.Delete(session.UserID, outfitId)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.file.Delete(ctx.UserContext(), h.cfg.Outfits, outfitId.String())
	if err != nil {
		middleware.LogWarning(ctx, err, "outfit_id", outfitId)
	}
	return ctx.SendString(common.EmptyJson)
}

func (h *OutfitHandler) Generate(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var req domain.OutfitGenerationRequest
	req = domain.OutfitGenerationRequest{
		Amount:   4,
		Prompt:   "something something",
		Purposes: []string{"clothes for outdoor"},
	}

	if err := easyjson.Unmarshal(ctx.Body(), &req); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	req.UserID = session.UserID
	req.Pos.IP = ctx.IP()
	fmt.Println("Generating outfit for: ", req.Pos.IP)

	err := h.generator.Generate(ctx.UserContext(), req)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(common.EmptyJson)
}

func (h *OutfitHandler) GetPurposes(ctx *fiber.Ctx) error {
	purposes, err := h.outfits.GetOutfitPurposes()
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(purposes)
}

func (h *OutfitHandler) GetGenerationResults(cfg *config.Centrifugo) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Error(err)
			}
		}()

		err := h.generator.ListenGenerationResults(h.logger, h.handleGenResults(cfg))
		if err != nil {
			h.logger.Errorw(err.Error())
		}
	}()
}

func (h *OutfitHandler) handleGenResults(cfg *config.Centrifugo) func(resp *domain.OutfitGenerationResponse) domain.Result {
	return func(resp *domain.OutfitGenerationResponse) domain.Result {
		userChannel := cfg.OutfitGenChannel + resp.UserID.String()

		bytes, err := easyjson.Marshal(resp)
		if err != nil {
			h.logger.Errorw(err.Error())
			return domain.ResultDiscard
		}

		h.logger.Infow("centrifugo", "channel", userChannel, "payload", string(bytes))

		centrifugoResp, err := h.centrifugo.Publish(
			context.Background(),
			&centrifugo.PublishRequest{
				Channel: userChannel,
				Data:    bytes,
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

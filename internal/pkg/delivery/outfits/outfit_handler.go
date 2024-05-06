package outfits

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	outfitRepo "try-on/internal/pkg/repository/sqlc/outfits"
	outfitUsecase "try-on/internal/pkg/usecase/outfits"
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/validate"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	easyjson "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type OutfitHandler struct {
	outfits   domain.OutfitUsecase
	generator domain.OutfitGenerator

	file domain.FileManager
	cfg  *config.Static

	logger    *zap.SugaredLogger
	publisher domain.ChannelPublisher[easyjson.Marshaler]
}

func New(
	db *pgxpool.Pool,
	generator domain.OutfitGenerator,
	file domain.FileManager,
	cfg *config.Static,
	logger *zap.SugaredLogger,
	publisher domain.ChannelPublisher[easyjson.Marshaler],
) *OutfitHandler {
	return &OutfitHandler{
		outfits:   outfitUsecase.New(outfitRepo.New(db)),
		generator: generator,
		file:      file,
		cfg:       cfg,
		logger:    logger,
		publisher: publisher,
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

func (h *OutfitHandler) GetOwn(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	outfits, err := h.outfits.GetByUser(session.UserID, false)
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

	userId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserIdInvalid
	}

	outfits, err := h.outfits.GetByUser(userId, true)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(outfits)
}

//easyjson:json
type createdResponse struct {
	Uuid  utils.UUID
	Image string
	domain.Timestamp
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
		Timestamp: domain.Timestamp{
			CreatedAt: outfit.CreatedAt,
			UpdatedAt: outfit.UpdatedAt,
		},
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

	fmt.Printf("outfit: %+v\n", outfit)

	err = validate.Struct(&outfit)
	if err != nil {
		return app_errors.ValidationError(err)
	}

	outfit.UserID = session.UserID
	outfit.ID = id

	transforms := ctx.FormValue("transforms")

	if err := easyjson.Unmarshal([]byte(transforms), &outfit.Transforms); err != nil {
		outfit.Transforms = nil
	}

	err = h.outfits.Update(&outfit)
	if err != nil {
		return app_errors.New(err)
	}

	fileHeader, err := ctx.FormFile("img")
	switch {
	case fileHeader == nil || err == fasthttp.ErrMissingFile || err == io.EOF:
		return ctx.SendString(common.EmptyJson)
	case err != nil:
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	default:
		break
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

	return ctx.JSON(createdResponse{
		Timestamp: domain.Timestamp{
			CreatedAt: outfit.CreatedAt,
			UpdatedAt: outfit.UpdatedAt,
		},
	})
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
	if err := ctx.QueryParser(&req); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	fmt.Printf("Got from query: %+v\n", req)

	req.UserID = session.UserID
	req.Pos.IP = ctx.IP()

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
	ctx := middleware.WithLogger(context.Background(), h.logger)

	return func(resp *domain.OutfitGenerationResponse) domain.Result {
		userChannel := cfg.OutfitGenChannel + resp.UserID.String()
		if !utils.HttpOk(resp.StatusCode) {
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

		h.publisher.Publish(
			ctx,
			userChannel,
			resp,
		)

		return domain.ResultOk
	}
}

package outfits

import (
	"time"

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
)

type OutfitHandler struct {
	outfits domain.OutfitUsecase
	file    domain.FileManager
	cfg     *config.Static
}

func New(db *pgxpool.Pool, file domain.FileManager, cfg *config.Static) *OutfitHandler {
	return &OutfitHandler{
		outfits: outfitUsecase.New(outfitRepo.New(db)),
		file:    file,
		cfg:     cfg,
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

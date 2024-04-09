package tryon

import (
	"context"
	"slices"

	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/clothes"
	"try-on/internal/pkg/repository/sqlc/outfits"
	"try-on/internal/pkg/repository/sqlc/user_images"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TryOnUsecase struct {
	clothes    domain.ClothesRepository
	outfits    domain.OutfitRepository
	userImages domain.UserImageRepository

	subscriber domain.Subscriber[domain.TryOnResponse]
	publisher  domain.Publisher[domain.TryOnRequest]
}

func New(
	db *pgxpool.Pool,
	pub domain.Publisher[domain.TryOnRequest],
	sub domain.Subscriber[domain.TryOnResponse],
) domain.TryOnUsecase {
	return &TryOnUsecase{
		clothes:    clothes.New(db),
		outfits:    outfits.New(db),
		userImages: user_images.New(db),
		publisher:  pub,
		subscriber: sub,
	}
}

func (u *TryOnUsecase) Close() {
	u.publisher.Close()
}

func (u *TryOnUsecase) TryOn(ctx context.Context, clothesIds []utils.UUID, opts domain.TryOnOpts) error {
	_, err := u.userImages.Get(opts.UserImageID)
	if err != nil {
		return err
	}

	clothes, err := u.clothes.GetTryOnInfo(clothesIds)
	if err != nil {
		return err
	}

	if err := validateTryOnCategories(clothes); err != nil {
		return err
	}

	return u.publisher.Publish(ctx, domain.TryOnRequest{
		TryOnOpts: opts,
		Clothes:   clothes,
	})
}

func (u *TryOnUsecase) TryOnOutfit(ctx context.Context, outfit utils.UUID, opts domain.TryOnOpts) error {
	_, err := u.userImages.Get(opts.UserImageID)
	if err != nil {
		return app_errors.New(err)
	}

	clothes, err := u.outfits.GetClothesInfo(outfit)
	if err != nil {
		return app_errors.New(err)
	}

	if err := validateTryOnCategories(clothes); err != nil {
		return err
	}

	return u.publisher.Publish(ctx, domain.TryOnRequest{
		TryOnOpts: opts,
		Clothes:   clothes,
	})
}

func (u *TryOnUsecase) GetTryOnResults(logger *zap.SugaredLogger, handler func(*domain.TryOnResponse) domain.Result) error {
	return nil
}

func filterClothesForTryOn(clothes []domain.TryOnClothesInfo) []domain.TryOnClothesInfo {
	dressIdx := slices.IndexFunc(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryDress
	})

	if dressIdx != -1 {
		return []domain.TryOnClothesInfo{clothes[dressIdx]}
	}

	result := make([]domain.TryOnClothesInfo, 0, 2)

	upperBodyIdx := slices.IndexFunc(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryUpper
	})

	if upperBodyIdx != -1 {
		result = append(result, clothes[upperBodyIdx])
	}

	lowerBodyIdx := slices.IndexFunc(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryLower
	})

	if lowerBodyIdx != -1 {
		result = append(result, clothes[lowerBodyIdx])
	}

	return result
}

func validateTryOnCategories(clothes []domain.TryOnClothesInfo) error {
	if len(clothes) < 1 || len(clothes) > 2 {
		return app_errors.ErrTryOnInvalidClothesNum
	}

	dressIdx := slices.IndexFunc(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryDress
	})
	if dressIdx != -1 && len(clothes) != 1 {
		return app_errors.ErrTryOnInvalidClothesType
	}

	upperCount := utils.Count(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryUpper
	})
	if upperCount > 1 {
		return app_errors.ErrTryOnInvalidClothesType
	}

	lowerCount := utils.Count(clothes, func(c domain.TryOnClothesInfo) bool {
		return c.Category == domain.TryOnCategoryLower
	})
	if lowerCount > 1 {
		return app_errors.ErrTryOnInvalidClothesType
	}

	return nil
}

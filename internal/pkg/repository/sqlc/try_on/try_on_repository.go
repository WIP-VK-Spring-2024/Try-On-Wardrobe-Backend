package try_on

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TryOnResultRepository struct {
	queries *sqlc.Queries
	db      *pgxpool.Pool
}

func New(db *pgxpool.Pool) domain.TryOnResultRepository {
	return &TryOnResultRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

func (repo TryOnResultRepository) Create(res *domain.TryOnResult) error {
	id, err := repo.queries.CreateTryOnResult(
		context.Background(),
		sqlc.CreateTryOnResultParams{
			UserImageID: res.UserImageID,
			ClothesID:   res.ClothesID,
			Image:       res.Image,
		},
	)
	if err != nil {
		return utils.PgxError(err)
	}

	res.ID = id
	return nil
}

func (repo TryOnResultRepository) Delete(id utils.UUID) error {
	err := repo.queries.DeleteTryOnResult(context.Background(), id)
	return utils.PgxError(err)
}

func (repo TryOnResultRepository) SetTryOnResultID(outfitId, tryOnResultId utils.UUID) error {
	err := repo.queries.SetOutfitTryOnResult(context.Background(), outfitId, tryOnResultId)
	return utils.PgxError(err)
}

func (repo TryOnResultRepository) GetByUser(userID utils.UUID) ([]domain.TryOnResult, error) {
	results, err := repo.queries.GetTryOnResultsByUser(context.Background(), userID)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(results, fromSqlc), nil
}

func (repo TryOnResultRepository) GetByClothes(userImageId utils.UUID, clothesID []utils.UUID) (*domain.TryOnResult, error) {
	result, err := repo.queries.GetTryOnResultByClothes(context.Background(), userImageId, clothesID)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return fromSqlc(&result), nil
}

func (repo TryOnResultRepository) GetByOutfit(userImageId, outfitId utils.UUID, updateOutfit bool) (*domain.TryOnResult, error) {
	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	queries := repo.queries.WithTx(tx)

	clothesIds, err := queries.GetClothesIdByOutfit(ctx, outfitId)
	if err != nil {
		return nil, utils.PgxError(err)
	}

	result, err := queries.GetTryOnResultByClothes(context.Background(), userImageId, clothesIds)
	if err != nil {
		return nil, utils.PgxError(err)
	}

	if updateOutfit {
		err = queries.SetOutfitTryOnResult(ctx, outfitId, result.ID)
		if err != nil {
			return nil, utils.PgxError(err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return fromSqlc(&result), nil
}

func (repo TryOnResultRepository) Get(id utils.UUID) (*domain.TryOnResult, error) {
	result, err := repo.queries.GetTryOnResult(context.Background(), id)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return fromSqlc(&result), nil
}

func (repo TryOnResultRepository) Rate(id utils.UUID, rating int) error {
	err := repo.queries.RateTryOnResult(context.Background(), id, int32(rating))
	return utils.PgxError(err)
}

func fromSqlc(model *sqlc.TryOnResult) *domain.TryOnResult {
	return &domain.TryOnResult{
		Model: domain.Model{
			ID: model.ID,
			Timestamp: domain.Timestamp{
				CreatedAt: utils.Time{Time: model.CreatedAt.Time},
				UpdatedAt: utils.Time{Time: model.UpdatedAt.Time},
			},
		},
		Image:       model.Image,
		UserImageID: model.UserImageID,
		ClothesID:   model.ClothesID,
		Rating:      int(model.Rating.Int32),
	}
}

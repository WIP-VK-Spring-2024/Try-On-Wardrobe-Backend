package try_on

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TryOnResultRepository struct {
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.TryOnResultRepository {
	return &TryOnResultRepository{
		queries: sqlc.New(db),
	}
}

func (repo *TryOnResultRepository) Create(res *domain.TryOnResult) error {
	id, err := repo.queries.CreateTryOnResult(
		context.Background(),
		res.ClothesID,
		res.UserImageID,
	)
	if err != nil {
		return utils.PgxError(err)
	}

	res.ID = id
	return nil
}

func (repo *TryOnResultRepository) Delete(id uuid.UUID) error {
	err := repo.queries.DeleteTryOnResult(context.Background(), id)
	return utils.PgxError(err)
}

func (repo *TryOnResultRepository) GetByUser(userID uuid.UUID) ([]domain.TryOnResult, error) {
	results, err := repo.queries.GetTryOnResultsByUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return utils.Map(results, fromSqlc), nil
}

func (repo *TryOnResultRepository) GetByClothes(clothesID uuid.UUID) ([]domain.TryOnResult, error) {
	results, err := repo.queries.GetTryOnResultsByClothes(context.Background(), clothesID)
	if err != nil {
		return nil, err
	}
	return utils.Map(results, fromSqlc), nil
}

func (repo *TryOnResultRepository) GetLast(userID uuid.UUID) (*domain.TryOnResult, error) {
	result, err := repo.queries.GetLastTryOnResult(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return fromSqlc(&result), nil
}

func (repo *TryOnResultRepository) Rate(id uuid.UUID, rating int) error {
	err := repo.queries.RateTryOnResult(context.Background(), id, int32(rating))
	return utils.PgxError(err)
}

func fromSqlc(model *sqlc.TryOnResult) *domain.TryOnResult {
	return &domain.TryOnResult{
		Model: domain.Model{
			ID: model.ID,
			AutoTimestamp: domain.AutoTimestamp{
				CreatedAt: model.CreatedAt.Time,
				UpdatedAt: model.UpdatedAt.Time,
			},
		},
		Image:       model.Image.String,
		UserImageID: model.UserImageID,
		ClothesID:   model.ClothesID,
		Rating:      int(model.Rating.Int32),
	}
}

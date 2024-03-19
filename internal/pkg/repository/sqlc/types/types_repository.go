package types

import (
	"context"
	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TypeRepository struct {
	db *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.TypeRepository {
	return &TypeRepository{
		db: sqlc.New(db),
	}
}

func (repo *TypeRepository) GetAll() ([]domain.Type, error) {
	types, err := repo.db.GetTypes(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(types, func(t *sqlc.Type) *domain.Type {
		return &domain.Type{
			Model: domain.Model{ID: t.ID},
			Name:  t.Name,
		}
	}), nil
}

package types

import (
	"context"
	"database/sql"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
)

type TypeRepository struct {
	db *sqlc.Queries
}

func New(db *sql.DB) domain.TypeRepository {
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
			Name:  t.Name.String,
		}
	}), nil
}

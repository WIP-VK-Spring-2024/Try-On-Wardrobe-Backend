package types

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TypeRepository struct {
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.TypeRepository {
	return &TypeRepository{
		queries: sqlc.New(db),
	}
}

func (repo *TypeRepository) GetAll() ([]domain.Type, error) {
	types, err := repo.queries.GetTypes(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(types, func(t *sqlc.GetTypesRow) *domain.Type {
		return &domain.Type{
			Model: domain.Model{
				ID: t.ID,
				AutoTimestamp: domain.AutoTimestamp{
					CreatedAt: t.CreatedAt.Time,
					UpdatedAt: t.UpdatedAt.Time,
				},
			},
			Name: t.Name,
			Subtypes: utils.Zip(t.SubtypeIds, t.SubtypeNames,
				func(id utils.UUID, name string) domain.Subtype {
					return domain.Subtype{
						Model: domain.Model{ID: id},
						Name:  name,
					}
				}),
		}
	}), nil
}

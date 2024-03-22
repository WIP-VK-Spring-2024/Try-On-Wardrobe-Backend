package subtypes

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubtypeRepository struct {
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.SubtypeRepository {
	return &SubtypeRepository{
		queries: sqlc.New(db),
	}
}

func (repo *SubtypeRepository) GetAll() ([]domain.Subtype, error) {
	types, err := repo.queries.GetSubtypes(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(types, func(t *sqlc.Subtype) *domain.Subtype {
		return &domain.Subtype{
			Model: domain.Model{
				ID: t.ID,
				AutoTimestamp: domain.AutoTimestamp{
					CreatedAt: t.CreatedAt.Time,
					UpdatedAt: t.UpdatedAt.Time,
				},
			},
			Name:   t.Name,
			TypeID: t.TypeID,
		}
	}), nil
}

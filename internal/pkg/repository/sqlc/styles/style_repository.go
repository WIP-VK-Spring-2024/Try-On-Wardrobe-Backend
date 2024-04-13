package styles

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StyleRepository struct {
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.StyleRepository {
	return &StyleRepository{
		queries: sqlc.New(db),
	}
}

func (repo StyleRepository) GetAll() ([]domain.Style, error) {
	types, err := repo.queries.GetStyles(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(types, func(t *sqlc.Style) *domain.Style {
		return &domain.Style{
			Model: domain.Model{
				ID: t.ID,
				Timestamp: domain.Timestamp{
					CreatedAt: utils.Time{Time: t.CreatedAt.Time},
					UpdatedAt: utils.Time{Time: t.UpdatedAt.Time},
				},
			},
			Name: t.Name,
		}
	}), nil
}

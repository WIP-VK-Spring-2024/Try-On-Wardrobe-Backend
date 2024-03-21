package styles

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StylesRepository struct {
	db *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.StylesRepository {
	return &StylesRepository{
		db: sqlc.New(db),
	}
}

func (repo *StylesRepository) GetAll() ([]domain.Style, error) {
	types, err := repo.db.GetStyles(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(types, func(t *sqlc.Style) *domain.Style {
		return &domain.Style{
			Model: domain.Model{
				ID: t.ID,
				AutoTimestamp: domain.AutoTimestamp{
					CreatedAt: t.CreatedAt.Time,
					UpdatedAt: t.UpdatedAt.Time,
				},
			},
			Name: t.Name,
		}
	}), nil
}

package tags

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepository struct {
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.TagRepository {
	return &TagRepository{
		queries: sqlc.New(db),
	}
}

func (repo TagRepository) Get(limit, offset int) ([]domain.Tag, error) {
	tags, err := repo.queries.GetTags(context.Background(), int32(limit), int32(offset))
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(tags, func(t *sqlc.Tag) *domain.Tag {
		return &domain.Tag{
			Model:    domain.Model{ID: t.ID},
			Name:     t.Name,
			UseCount: t.UseCount,
		}
	}), nil
}

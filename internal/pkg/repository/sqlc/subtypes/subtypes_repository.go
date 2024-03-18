package subtypes

import (
	"context"
	"database/sql"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
)

type SubtypeRepository struct {
	db *sqlc.Queries
}

func New(db *sql.DB) domain.SubtypeRepository {
	return &SubtypeRepository{
		db: sqlc.New(db),
	}
}

func (repo *SubtypeRepository) GetAll() ([]domain.Subtype, error) {
	types, err := repo.db.GetSubtypes(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(types, func(t *sqlc.Subtype) *domain.Subtype {
		return &domain.Subtype{
			Model:  domain.Model{ID: t.ID},
			Name:   t.Name.String,
			TypeID: t.TypeID,
		}
	}), nil
}

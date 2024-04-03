package classification

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClothesClassificationRepository struct {
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.ClothesClassificationRepository {
	return &ClothesClassificationRepository{
		queries: sqlc.New(db),
	}
}

func (c ClothesClassificationRepository) GetClassifications(tagLimit int32) (*domain.ClothesClassificationRequest, error) {
	tagNames, err := c.queries.GetTagEngNames(context.Background(), tagLimit, 0)
	if err != nil {
		return nil, utils.PgxError(err)
	}

	styleNames, err := c.queries.GetStyleEngNames(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	typeNames, err := c.queries.GetTypeEngNames(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	subtypeNames, err := c.queries.GetSubtypeEngNames(context.Background())
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return &domain.ClothesClassificationRequest{
		Styles:        styleNames,
		Categories:    typeNames,
		Subcategories: subtypeNames,
		Seasons:       domain.Seasons,
		Tags: utils.Map(tagNames, func(t *pgtype.Text) *string {
			return &t.String
		}),
	}, nil
}

func (c ClothesClassificationRepository) GetTypeId(engName string) (utils.UUID, error) {
	return c.queries.GetTypeIdByEngName(context.Background(), engName)
}

func (c ClothesClassificationRepository) GetSubtypeIds(engNames []string) ([]utils.UUID, error) {
	return c.queries.GetSubtypeIdsByEngName(context.Background(), engNames)
}

func (c ClothesClassificationRepository) GetStyleId(engName string) (utils.UUID, error) {
	return c.queries.GetStyleIdByEngName(context.Background(), engName)
}

func (c ClothesClassificationRepository) GetTags(engNames []string) ([]string, error) {
	return c.queries.GetTagsByEngName(context.Background(), engNames)
}

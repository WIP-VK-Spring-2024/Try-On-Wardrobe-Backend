package clothes

import (
	"context"
	"database/sql"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/translate"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClothesRepository struct {
	queries *sqlc.Queries
	db      *pgxpool.Pool
}

func New(db *pgxpool.Pool) domain.ClothesRepository {
	return &ClothesRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

func (c *ClothesRepository) Create(clothes *domain.ClothesModel) error {
	ctx := context.Background()

	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := c.queries.WithTx(tx)
	createParams := sqlc.CreateClothesParams{
		UserID:    clothes.UserID,
		Name:      clothes.Name,
		TypeID:    clothes.TypeID,
		SubtypeID: clothes.SubtypeID,
		Color:     pgtype.Text(clothes.Color),
	}

	tags := translate.TagsToString(clothes.Tags)

	err = queries.CreateTags(ctx, tags)
	if err != nil {
		return utils.PgxError(err)
	}

	clothesId, err := queries.CreateClothes(ctx, createParams)
	if err != nil {
		return utils.PgxError(err)
	}

	clothes.ID = clothesId

	err = queries.CreateClothesTagLinks(ctx, clothes.ID,
		tags,
	)
	if err != nil {
		return utils.PgxError(err)
	}

	return tx.Commit(ctx)
}

func (c *ClothesRepository) Update(clothes *domain.ClothesModel) error {
	ctx := context.Background()

	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := c.queries.WithTx(tx)
	updateParams := sqlc.UpdateClothesParams{
		ID:        clothes.ID,
		Name:      clothes.Name,
		Note:      pgtype.Text(clothes.Note),
		TypeID:    clothes.TypeID,
		SubtypeID: clothes.SubtypeID,
		Color:     pgtype.Text(clothes.Color),
		Seasons: utils.Map(clothes.Seasons, func(t *domain.Season) *sqlc.Season {
			tmp := sqlc.Season(*t)
			return &tmp
		}),
	}

	if clothes.Style != nil {
		styleId, err := queries.CreateStyle(ctx, clothes.Style.Name)
		if err != nil {
			return utils.PgxError(err)
		}
		updateParams.StyleID = styleId
	}

	err = c.queries.UpdateClothes(ctx, updateParams)
	if err != nil {
		return utils.PgxError(err)
	}

	tags := translate.TagsToString(clothes.Tags)

	err = queries.CreateTags(ctx, tags)
	if err != nil {
		return utils.PgxError(err)
	}

	err = queries.CreateClothesTagLinks(ctx,
		clothes.ID,
		tags,
	)
	if err != nil {
		return utils.PgxError(err)
	}

	return tx.Commit(ctx)
}

func (c *ClothesRepository) Get(id utils.UUID) (*domain.ClothesModel, error) {
	clothes, err := c.queries.GetClothesById(context.Background(), id)
	if err != nil {
		return nil, utils.PgxError(err)
	}

	clothesCompat := sqlc.GetClothesByUserRow(clothes) // костыль, но в гошке иначе нельзя
	return fromSqlc(&clothesCompat), nil
}

func (c *ClothesRepository) Delete(id utils.UUID) error {
	return utils.PgxError(c.queries.DeleteClothes(context.Background(), id))
}

func (c *ClothesRepository) GetByUser(userID utils.UUID, _ int) ([]domain.ClothesModel, error) {
	clothes, err := c.queries.GetClothesByUser(context.Background(), userID)
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(clothes, fromSqlc), nil
}

func fromSqlc(model *sqlc.GetClothesByUserRow) *domain.ClothesModel {
	result := &domain.ClothesModel{
		Model: domain.Model{
			ID: model.ID,
			AutoTimestamp: domain.AutoTimestamp{
				CreatedAt: model.CreatedAt.Time,
				UpdatedAt: model.UpdatedAt.Time,
			},
		},
		TypeID:    model.TypeID,
		SubtypeID: model.SubtypeID,
		UserID:    model.UserID,
		StyleID:   model.StyleID,
		Color:     sql.NullString(model.Color),
		Name:      model.Name,
		Note:      sql.NullString(model.Note),
		Type: domain.Type{
			Name: model.Type,
		},
		Subtype: domain.Subtype{
			Name: model.Subtype,
		},
		Tags: translate.TagsFromString(model.Tags),
	}

	if model.Style.Valid {
		result.Style = &domain.Style{
			Name: model.Style.String,
		}
	}

	return result
}

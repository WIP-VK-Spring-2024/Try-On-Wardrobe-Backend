package clothes

import (
	"context"
	"database/sql"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/optional"

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

func (c ClothesRepository) Create(clothes *domain.Clothes) error {
	ctx := context.Background()

	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := c.queries.WithTx(tx)
	createParams := sqlc.CreateClothesParams{
		UserID: clothes.UserID,
		Name:   clothes.Name,
		Image:  clothes.Image,
	}

	clothesId, err := queries.CreateClothes(ctx, createParams)
	if err != nil {
		return utils.PgxError(err)
	}

	clothes.ID = clothesId

	err = queries.SetClothesImage(ctx, clothesId, clothes.Image+"/"+clothesId.String())
	if err != nil {
		return utils.PgxError(err)
	}

	return tx.Commit(ctx)
}

func (c ClothesRepository) SetImage(id utils.UUID, path string) error {
	return utils.PgxError(c.queries.SetClothesImage(context.Background(), id, path))
}

func (c ClothesRepository) GetTryOnInfo(ids []utils.UUID) ([]domain.TryOnClothesInfo, error) {
	clothesInfo, err := c.queries.GetClothesTryOnInfo(context.Background(), ids)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(clothesInfo, func(t *sqlc.GetClothesTryOnInfoRow) *domain.TryOnClothesInfo {
		return &domain.TryOnClothesInfo{
			ClothesID: t.ID,
			Category:  t.Category,
			Layer:     int(t.Layer),
		}
	}), nil
}

func (c ClothesRepository) GetByWeather(userId utils.UUID, temp *int) ([]domain.GenClothesInfo, error) {
	tempParam := pgtype.Int4{}
	if temp != nil {
		tempParam.Int32 = int32(*temp)
		tempParam.Valid = true
	}

	clothesInfo, err := c.queries.GetClothesInfoByWeather(context.Background(), userId, tempParam)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(clothesInfo, func(t *sqlc.GetClothesInfoByWeatherRow) *domain.GenClothesInfo {
		return &domain.GenClothesInfo{
			ClothesID: t.ID,
			Category:  t.Category,
		}
	}), nil
}

func (c ClothesRepository) Update(clothes *domain.Clothes) error {
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
		Note:      pgtype.Text(clothes.Note.NullString),
		TypeID:    clothes.TypeID,
		SubtypeID: clothes.SubtypeID,
		StyleID:   clothes.StyleID,
		Color:     pgtype.Text(clothes.Color.NullString),
		Seasons:   clothes.Seasons,
		Privacy:   clothes.Privacy,
	}

	err = c.queries.UpdateClothes(ctx, updateParams)
	if err != nil {
		return utils.PgxError(err)
	}

	err = queries.DeleteClothesTagLinks(ctx, clothes.ID, clothes.Tags)
	if err != nil {
		return utils.PgxError(err)
	}

	err = queries.CreateClothesTagLinks(ctx, clothes.ID, clothes.Tags)
	if err != nil {
		return utils.PgxError(err)
	}

	return tx.Commit(ctx)
}

func (c ClothesRepository) Get(id utils.UUID) (*domain.Clothes, error) {
	clothes, err := c.queries.GetClothesById(context.Background(), id)
	if err != nil {
		return nil, utils.PgxError(err)
	}

	clothesCompat := sqlc.GetClothesByUserRow(clothes) // костыль, но в гошке иначе нельзя
	return fromSqlc(&clothesCompat), nil
}

func (c ClothesRepository) Delete(id utils.UUID) error {
	return utils.PgxError(c.queries.DeleteClothes(context.Background(), id))
}

func (c ClothesRepository) GetByUser(userID utils.UUID, _ int) ([]domain.Clothes, error) {
	clothes, err := c.queries.GetClothesByUser(context.Background(), userID)
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(clothes, fromSqlc), nil
}

func fromSqlc(model *sqlc.GetClothesByUserRow) *domain.Clothes {
	result := &domain.Clothes{
		Model: domain.Model{
			ID: model.ID,
			Timestamp: domain.Timestamp{
				CreatedAt: utils.Time{Time: model.CreatedAt.Time},
				UpdatedAt: utils.Time{Time: model.UpdatedAt.Time},
			},
		},
		Image:     model.Image,
		TypeID:    model.TypeID,
		SubtypeID: model.SubtypeID,
		UserID:    model.UserID,
		StyleID:   model.StyleID,
		Color:     optional.String{NullString: sql.NullString(model.Color)},
		Name:      model.Name,
		Note:      optional.String{NullString: sql.NullString(model.Note)},
		Type:      model.Type.String,
		Subtype:   model.Subtype.String,
		Tags:      model.Tags,
		Seasons:   model.Seasons,
		Tryonable: model.Tryonable,
	}

	if model.Style.Valid {
		result.Style = model.Style.String
	}

	return result
}

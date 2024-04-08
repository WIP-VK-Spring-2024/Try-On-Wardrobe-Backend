package outfits

import (
	"context"
	"database/sql"
	"time"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/optional"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type OutfitRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.OutfitRepository {
	return &OutfitRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (repo OutfitRepository) Create(outfit *domain.Outfit) error {
	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := repo.queries.WithTx(tx)

	transforms, err := easyjson.Marshal(outfit.Transforms)
	if err != nil {
		return err
	}

	id, err := queries.CreateOutfit(context.Background(), outfit.UserID, transforms)
	if err != nil {
		return utils.PgxError(err)
	}

	err = queries.SetOutfitImage(ctx, id, outfit.Image+"/"+id.String())
	if err != nil {
		return utils.PgxError(err)
	}

	outfit.ID = id
	return tx.Commit(ctx)
}

func (repo OutfitRepository) Update(outfit *domain.Outfit) (err error) {
	var transforms []byte
	if outfit.Transforms != nil {
		transforms, err = easyjson.Marshal(outfit.Transforms)
		if err != nil {
			return err
		}
	}

	updateParams := sqlc.UpdateOutfitParams{
		ID:         outfit.ID,
		Note:       pgtype.Text(outfit.Note.NullString),
		StyleID:    outfit.StyleID,
		Seasons:    outfit.Seasons,
		Transforms: transforms,
	}

	if outfit.Name != "" {
		updateParams.Name.String = outfit.Name
		updateParams.Name.Valid = true
	}

	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := repo.queries.WithTx(tx)

	err = queries.UpdateOutfit(context.Background(), updateParams)
	if err != nil {
		return utils.PgxError(err)
	}

	err = queries.CreateTags(ctx, outfit.Tags)
	if err != nil {
		return utils.PgxError(err)
	}

	err = queries.DeleteOutfitTagLinks(ctx, outfit.ID, outfit.Tags)
	if err != nil {
		return utils.PgxError(err)
	}

	err = queries.CreateOutfitTagLinks(ctx, outfit.ID, outfit.Tags)
	if err != nil {
		return utils.PgxError(err)
	}

	return tx.Commit(ctx)
}

func (repo OutfitRepository) Delete(id utils.UUID) error {
	return utils.PgxError(repo.queries.DeleteOutfit(context.Background(), id))
}

func (repo OutfitRepository) GetById(id utils.UUID) (*domain.Outfit, error) {
	outfit, err := repo.queries.GetOutfit(context.Background(), id)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return fromSqlc(&outfit), nil
}

func (repo OutfitRepository) GetByUser(userId utils.UUID) ([]domain.Outfit, error) {
	outfits, err := repo.queries.GetOutfitsByUser(context.Background(), userId)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(outfits, fromGetOutfitsByUser), nil
}

func (repo OutfitRepository) GetClothesInfo(outfitId utils.UUID) ([]domain.TryOnClothesInfo, error) {
	clothesInfoSlice, err := repo.queries.GetOutfitClothesInfo(context.Background(), outfitId)
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(clothesInfoSlice, func(t *sqlc.GetOutfitClothesInfoRow) *domain.TryOnClothesInfo {
		return &domain.TryOnClothesInfo{
			ClothesID: t.ID,
			Category:  t.Category,
		}
	}), nil
}

func (repo OutfitRepository) Get(since time.Time, limit int) ([]domain.Outfit, error) {
	outfits, err := repo.queries.GetOutfits(
		context.Background(),
		pgtype.Timestamptz{Time: since, Valid: true},
		int32(limit),
	)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(outfits, fromGetOutfits), nil
}

func fromGetOutfits(value *sqlc.GetOutfitsRow) *domain.Outfit {
	model := sqlc.GetOutfitRow(*value)
	return fromSqlc(&model)
}

func fromGetOutfitsByUser(value *sqlc.GetOutfitsByUserRow) *domain.Outfit {
	model := sqlc.GetOutfitRow(*value)
	return fromSqlc(&model)
}

func fromSqlc(model *sqlc.GetOutfitRow) *domain.Outfit {
	result := &domain.Outfit{
		Model: domain.Model{
			ID: model.ID,
			AutoTimestamp: domain.AutoTimestamp{
				CreatedAt: utils.Time{Time: model.CreatedAt.Time},
				UpdatedAt: utils.Time{Time: model.UpdatedAt.Time},
			},
		},
		Public:  model.Public,
		UserID:  model.UserID,
		StyleID: model.StyleID,
		Name:    model.Name.String,
		Note:    optional.String{NullString: sql.NullString(model.Note)},
		Image:   model.Image.String,
		Seasons: model.Seasons,
		Tags:    model.Tags,
	}

	err := easyjson.Unmarshal(model.Transforms, &result.Transforms)
	if err != nil {
		zap.S().Errorw("outfit transform map unmarshalling" + err.Error())
	}

	return result
}

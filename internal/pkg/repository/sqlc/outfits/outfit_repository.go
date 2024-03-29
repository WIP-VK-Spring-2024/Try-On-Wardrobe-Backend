package outfits

import (
	"context"
	"database/sql"

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

func (repo *OutfitRepository) Create(outfit *domain.Outfit) error {
	transforms, err := easyjson.Marshal(outfit.Transforms)
	if err != nil {
		return err
	}

	id, err := repo.queries.CreateOutfit(context.Background(), outfit.UserID, transforms)
	if err != nil {
		return utils.PgxError(err)
	}

	outfit.ID = id
	return nil
}

func (repo *OutfitRepository) Update(outfit *domain.Outfit) (err error) {
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

	err = repo.queries.UpdateOutfit(context.Background(), updateParams)

	// TODO: Add tags
	return utils.PgxError(err)
}

func (repo *OutfitRepository) Delete(id utils.UUID) error {
	return utils.PgxError(repo.queries.DeleteOutfit(context.Background(), id))
}

func (repo *OutfitRepository) Get(id utils.UUID) (*domain.Outfit, error) {
	outfit, err := repo.queries.GetOutfit(context.Background(), id)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return fromSqlc(&outfit), nil
}

func (repo *OutfitRepository) GetByUser(userId utils.UUID) ([]domain.Outfit, error) {
	outfits, err := repo.queries.GetOutfitsByUser(context.Background(), userId)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(outfits, fromSqlc), nil
}

func fromSqlc(model *sqlc.Outfit) *domain.Outfit {
	result := &domain.Outfit{
		Model: domain.Model{
			ID: model.ID,
			AutoTimestamp: domain.AutoTimestamp{
				CreatedAt: utils.Time{Time: model.CreatedAt.Time},
				UpdatedAt: utils.Time{Time: model.UpdatedAt.Time},
			},
		},
		UserID:  model.UserID,
		StyleID: model.StyleID,
		Name:    model.Name.String,
		Note:    optional.String{NullString: sql.NullString(model.Note)},
		Image:   model.Image.String,
		Seasons: model.Seasons,
	}

	err := easyjson.Unmarshal(model.Transforms, &result.Transforms)
	if err != nil {
		zap.S().Errorw("outfit transform map unmarshalling" + err.Error())
	}

	return result
}

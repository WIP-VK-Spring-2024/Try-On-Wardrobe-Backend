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
	db      *pgxpool.Pool
}

func New(db *pgxpool.Pool) domain.TagRepository {
	return &TagRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

func (repo TagRepository) GetUserFavourite(userId utils.UUID, limit int) ([]domain.Tag, error) {
	tags, err := repo.queries.GetUserFavouriteTags(context.Background(), userId, int32(limit))
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(tags, fromSqlc), nil
}

func (repo TagRepository) GetNotCreated(tags []string) ([]string, error) {
	return repo.queries.GetNotCreatedTags(context.Background(), tags)
}

// Should've done this using a single query, but sqlc has a bug I don't want to bypass
func (repo TagRepository) SetEngNames(tags, engNames []string) error {
	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := repo.queries.WithTx(tx)

	for i := range min(len(tags), len(engNames)) {
		err = queries.SetTagEngName(ctx, tags[i], engNames[i])
		if err != nil {
			return utils.PgxError(err)
		}
	}

	return tx.Commit(ctx)
}

// Should've done this using a single query, but sqlc has ANOTHER bug I don't want to bypass
func (repo TagRepository) Create(tags []domain.Tag) error {
	ctx := context.Background()
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := repo.queries.WithTx(tx)
	for _, tag := range tags {
		err = queries.CreateTagsWithEng(context.Background(), tag.Name, tag.EngName)
		if err != nil {
			return utils.PgxError(err)
		}
	}
	return tx.Commit(ctx)
}

func (repo TagRepository) Get(limit, offset int) ([]domain.Tag, error) {
	tags, err := repo.queries.GetTags(context.Background(), int32(limit), int32(offset))
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(tags, fromSqlc), nil
}

func fromSqlc(t *sqlc.Tag) *domain.Tag {
	return &domain.Tag{
		Model:    domain.Model{ID: t.ID},
		Name:     t.Name,
		UseCount: t.UseCount,
	}
}

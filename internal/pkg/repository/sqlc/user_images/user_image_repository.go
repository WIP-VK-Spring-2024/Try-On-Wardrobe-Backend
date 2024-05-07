package user_images

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserImageRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.UserImageRepository {
	return &UserImageRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

func (repo UserImageRepository) Create(userImage *domain.UserImage) error {
	id, err := repo.queries.CreateUserImage(context.Background(), userImage.UserID, userImage.Image)
	if err != nil {
		return utils.PgxError(err)
	}
	userImage.ID = id

	return nil
}

func (repo UserImageRepository) SetUserImageUrl(id utils.UUID, url string) error {
	err := repo.queries.SetUserImageUrl(context.Background(), id, url)
	return utils.PgxError(err)
}

func (repo UserImageRepository) Delete(id utils.UUID) error {
	err := repo.queries.DeleteUserImage(context.Background(), id)
	return utils.PgxError(err)
}

func (repo UserImageRepository) GetByUser(userID utils.UUID) ([]domain.UserImage, error) {
	userImages, err := repo.queries.GetUserImageByUser(context.Background(), userID)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(userImages, fromSqlc), nil
}

func (repo UserImageRepository) Get(id utils.UUID) (*domain.UserImage, error) {
	userImage, err := repo.queries.GetUserImageByID(context.Background(), id)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return fromSqlc(&userImage), nil
}

func fromSqlc(model *sqlc.UserImage) *domain.UserImage {
	return &domain.UserImage{
		Model: domain.Model{
			ID: model.ID,
			Timestamp: domain.Timestamp{
				CreatedAt: utils.Time{Time: model.CreatedAt.Time},
				UpdatedAt: utils.Time{Time: model.UpdatedAt.Time},
			},
		},
		UserID: model.UserID,
		Image:  model.Image,
	}
}

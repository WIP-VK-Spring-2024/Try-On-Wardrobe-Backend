package users

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.UserRepository {
	return UserRepository{
		queries: sqlc.New(db),
	}
}

func (repo UserRepository) Create(user *domain.User) error {
	id, err := repo.queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(user.Password),
		Gender:   user.Gender,
	})
	if err != nil {
		return utils.PgxError(err)
	}

	user.ID = id
	return nil
}

func (repo UserRepository) Update(user domain.User) error {
	params := sqlc.UpdateUserParams{
		ID:     user.ID,
		Name:   user.Name,
		Avatar: user.Avatar,
	}

	if user.Gender == domain.Male || user.Gender == domain.Female {
		params.Gender.Gender = sqlc.Gender(user.Gender)
		params.Gender.Valid = true
	}

	if user.Privacy == domain.PrivacyPrivate || user.Privacy == domain.PrivacyPublic {
		params.Privacy.Privacy = sqlc.Privacy(user.Privacy)
		params.Privacy.Valid = true
	}

	err := repo.queries.UpdateUser(context.Background(), params)
	return utils.PgxError(err)
}

func (repo UserRepository) GetByName(name string) (*domain.User, error) {
	user, err := repo.queries.GetUserByName(context.Background(), name)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return fromSqlc(&user), nil
}

func (repo UserRepository) GetByEmail(email string) (*domain.User, error) {
	user, err := repo.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return fromSqlc(&user), nil
}

func (repo UserRepository) GetByID(id utils.UUID) (*domain.User, error) {
	user, err := repo.queries.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return fromSqlc(&user), nil
}

func (repo UserRepository) GetSubscriptions(userId utils.UUID) ([]domain.User, error) {
	users, err := repo.queries.GetSubscribedToUsers(context.Background(), userId)
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(users, fromSqlc), nil
}

func (repo UserRepository) SearchUsers(opts domain.SearchUserOpts) ([]domain.User, error) {
	users, err := repo.queries.SearchUsers(context.Background(), sqlc.SearchUsersParams{
		SubscriberID: opts.UserID,
		Name:         "%" + opts.Name + "%",
		Since:        opts.Since,
		Limit:        int32(opts.Limit),
	})
	if err == pgx.ErrNoRows {
		return []domain.User{}, nil
	}
	if err != nil {
		return nil, utils.PgxError(err)
	}
	return utils.Map(users, fromSqlc), nil
}

func fromSqlc(model *sqlc.User) *domain.User {
	return &domain.User{
		Model: domain.Model{
			ID: model.ID,
			Timestamp: domain.Timestamp{
				CreatedAt: utils.Time{Time: model.CreatedAt.Time},
				UpdatedAt: utils.Time{Time: model.UpdatedAt.Time},
			},
		},
		Name:     model.Name,
		Password: model.Password,
		Email:    model.Email,
		Gender:   model.Gender,
		Privacy:  model.Privacy,
		Avatar:   model.Avatar,
	}
}

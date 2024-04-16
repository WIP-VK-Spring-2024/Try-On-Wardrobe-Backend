package users

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgtype"
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
		Email:    pgtype.Text{String: user.Email, Valid: true},
		Password: string(user.Password),
	})
	if err != nil {
		return utils.PgxError(err)
	}

	user.ID = id
	return nil
}

func (repo UserRepository) Update(user domain.User) error {
	params := sqlc.UpdateUserParams{
		ID:      user.ID,
		Name:    user.Name,
		Privacy: user.Privacy,
		Avatar:  user.Avatar,
	}

	if user.Gender != "" {
		params.Gender = sqlc.NullGender{
			Gender: sqlc.Gender(user.Gender),
			Valid:  true,
		}
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

func (repo UserRepository) SearchUsers(name string) ([]domain.User, error) {
	users, err := repo.queries.SearchUsers(context.Background(), "%"+name+"%")
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
		Password: []byte(model.Password),
		Email:    model.Email.String,
	}
}

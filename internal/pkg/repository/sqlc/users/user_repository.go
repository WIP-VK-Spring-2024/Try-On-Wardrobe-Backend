package users

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
		Email:    pgtype.Text(user.Email.NullString),
		Password: string(user.Password),
	})
	if err != nil {
		return utils.PgxError(err)
	}

	user.ID = id
	return nil
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
		Email:    optional.String{NullString: sql.NullString(model.Email)},
	}
}

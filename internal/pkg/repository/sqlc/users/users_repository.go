package users

import (
	"context"
	"database/sql"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"

	"github.com/google/uuid"
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
		Email:    pgtype.Text(user.Email),
		Password: string(user.Password),
	})
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (repo UserRepository) GetByName(name string) (*domain.User, error) {
	user, err := repo.queries.GetUserByName(context.Background(), name)
	if err != nil {
		return nil, err
	}
	return fromSqlc(&user), nil
}

func (repo UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	user, err := repo.queries.GetUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return fromSqlc(&user), nil
}

func fromSqlc(model *sqlc.User) *domain.User {
	return &domain.User{
		Model: domain.Model{
			ID: model.ID,
			AutoTimestamp: domain.AutoTimestamp{
				CreatedAt: model.CreatedAt.Time,
				UpdatedAt: model.UpdatedAt.Time,
			},
		},
		Name:     model.Name,
		Password: []byte(model.Password),
		Email:    sql.NullString(model.Email),
	}
}

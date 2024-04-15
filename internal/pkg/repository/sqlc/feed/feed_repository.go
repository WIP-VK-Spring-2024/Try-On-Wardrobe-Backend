package feed

import (
	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedRepository struct {
	queries *sqlc.Queries
}

func New(db *pgxpool.Pool) domain.FeedRepository {
	return &FeedRepository{
		queries: sqlc.New(db),
	}
}

func (f *FeedRepository) GetPosts(opts domain.GetPostsOpts) ([]domain.Post, error) {
	return nil, nil
}

func (f *FeedRepository) GetPost(postId utils.UUID) (*domain.Post, error) {
	return nil, nil
}

func (f *FeedRepository) GetComments(postId utils.UUID) ([]domain.Comment, error) {
	return nil, nil
}

func (f *FeedRepository) RatePost(postId utils.UUID, rating int) error {
	return nil
}

func (f *FeedRepository) RateComment(commentId utils.UUID, rating int) error {
	return nil
}

func (f *FeedRepository) Comment(postId utils.UUID, comment domain.CommentModel) error {
	return nil
}

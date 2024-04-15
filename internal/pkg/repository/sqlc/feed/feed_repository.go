package feed

import (
	"context"

	"try-on/internal/generated/sqlc"
	"try-on/internal/pkg/app_errors"
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
	posts, err := f.queries.GetPosts(context.Background(), sqlc.GetPostsParams{
		UserID: opts.RequestingUserID,
		Limit:  opts.Limit,
		Since:  opts.Since,
	})
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(posts, postsFromSqlc), nil
}

func (f *FeedRepository) GetPost(postId utils.UUID) (*domain.Post, error) {
	return nil, app_errors.ErrUnimplemented
}

func (f *FeedRepository) GetComments(opts domain.GetCommentsOpts) ([]domain.Comment, error) {
	comments, err := f.queries.GetComments(context.Background(), sqlc.GetCommentsParams{
		UserID: opts.RequestingUserID,
		PostID: opts.PostID,
		Limit:  opts.Limit,
		Since:  opts.Since,
	})
	if err != nil {
		return nil, utils.PgxError(err)
	}

	return utils.Map(comments, commentsFromSqlc), nil
}

func (f *FeedRepository) RatePost(userId, postId utils.UUID, rating int) error {
	err := f.queries.RatePost(context.Background(), sqlc.RatePostParams{
		UserID: userId,
		PostID: postId,
		Value:  int32(rating),
	})
	return utils.PgxError(err)
}

func (f *FeedRepository) RateComment(userId, commentId utils.UUID, rating int) error {
	err := f.queries.RateComment(context.Background(), sqlc.RateCommentParams{
		UserID:    userId,
		CommentID: commentId,
		Value:     int32(rating),
	})
	return utils.PgxError(err)
}

func (f *FeedRepository) Comment(postId utils.UUID, comment domain.CommentModel) error {
	_, err := f.queries.CreateComment(context.Background(), sqlc.CreateCommentParams{
		PostID: postId,
		UserID: comment.UserID,
		Body:   comment.Body,
	})
	return utils.PgxError(err)
}

func postsFromSqlc(model *sqlc.GetPostsRow) *domain.Post {
	return &domain.Post{
		Model: domain.Model{
			ID: model.ID,
			Timestamp: domain.Timestamp{
				CreatedAt: model.CreatedAt,
				UpdatedAt: model.UpdatedAt,
			},
		},
		OutfitID:    model.OutfitID,
		OutfitImage: model.OutfitImage.String,
		UserID:      model.UserID,
		UserImage:   model.UserImage,
		Rating:      int(model.Rating),
		Liked:       model.Liked,
	}
}

func commentsFromSqlc(model *sqlc.GetCommentsRow) *domain.Comment {
	return &domain.Comment{
		Model: domain.Model{
			ID: model.ID,
			Timestamp: domain.Timestamp{
				CreatedAt: model.CreatedAt,
				UpdatedAt: model.UpdatedAt,
			},
		},
		CommentModel: domain.CommentModel{
			UserID: model.UserID,
			Body:   model.Body,
		},
		UserImage: model.UserImage,
		Rating:    int(model.Rating),
		Liked:     model.Liked,
	}
}

package domain

import (
	"try-on/internal/pkg/utils"
)

//easyjson:json
type Post struct {
	Model

	OutfitID utils.UUID
	Rating   int

	OutfitImage string
	TryOnImage  string

	UserID    utils.UUID
	UserImage string

	Liked bool
}

//easyjson:json
type Comment struct {
	Model
	CommentModel

	UserImage string

	Rating int
	Liked  bool
}

//easyjson:json
type CommentModel struct {
	UserID utils.UUID
	Body   string
}

//easyjson:json
type GetPostsOpts struct {
	RequestingUserID utils.UUID `json:"-"`
	Limit            int32      `query:"limit"`
	Since            utils.Time `query:"since"`
}

//easyjson:json
type GetCommentsOpts struct {
	PostID           utils.UUID
	RequestingUserID utils.UUID `json:"-"`
	Limit            int32      `query:"limit"`
	Since            utils.Time `query:"since"`
}

type FeedRepository interface {
	GetPosts(opts GetPostsOpts) ([]Post, error)
	GetPost(postId utils.UUID) (*Post, error)
	GetComments(opts GetCommentsOpts) ([]Comment, error)

	RatePost(userId, postId utils.UUID, rating int) error
	RateComment(userId, commentId utils.UUID, rating int) error

	Comment(postId utils.UUID, comment CommentModel) error
}

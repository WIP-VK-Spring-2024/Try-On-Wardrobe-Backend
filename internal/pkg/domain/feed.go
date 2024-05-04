package domain

import (
	"try-on/internal/pkg/utils"
)

//easyjson:json
type Post struct {
	Model

	OutfitID    utils.UUID
	OutfitImage string

	TryOnID    utils.UUID
	TryOnImage string

	UserID    utils.UUID
	UserName  string
	UserImage string
	IsSubbed  bool `json:"is_subbed,!omitempty"` //lint:ignore SA5008 easyjson custom tags

	Rating     int `json:"rating,!omitempty"`      //lint:ignore SA5008 easyjson custom tags
	UserRating int `json:"user_rating,!omitempty"` //lint:ignore SA5008 easyjson custom tags
}

//easyjson:json
type Comment struct {
	Model
	CommentModel

	UserName  string
	UserImage string

	Rating     int `json:"rating,!omitempty"`      //lint:ignore SA5008 easyjson custom tags
	UserRating int `json:"user_rating,!omitempty"` //lint:ignore SA5008 easyjson custom tags

	Level int `json:"level,!omitempty"` //lint:ignore SA5008 easyjson custom tags

	Replies []Comment
}

//easyjson:json
type CommentModel struct {
	UserID   utils.UUID
	Body     string
	ParentID utils.UUID
}

type GetPostsOpts struct {
	RequestingUserID utils.UUID
	Limit            int32      `query:"limit"`
	Since            utils.Time `query:"since"`
	Genders          []Gender   `query:"gender"`
	Query            string     `query:"query"`
	Tags             []string   `query:"tags"`
	Style            utils.UUID `query:"style"`
}

type GetCommentsOpts struct {
	PostID           utils.UUID
	RequestingUserID utils.UUID
	Limit            int32      `query:"limit"`
	Since            utils.Time `query:"since"`
}

type FeedRepository interface {
	GetPosts(opts GetPostsOpts) ([]Post, error)
	GetLikedPosts(opts GetPostsOpts) ([]Post, error)
	GetSubscriptionPosts(opts GetPostsOpts) ([]Post, error)
	GetPostsByUser(userId utils.UUID, opts GetPostsOpts) ([]Post, error)

	Subscribe(subscriberId, userId utils.UUID) error
	Unsubscribe(subscriberId, userId utils.UUID) error

	GetPost(postId utils.UUID) (*Post, error)
	GetComments(opts GetCommentsOpts) ([]Comment, error)
	GetCommentsTree(opts GetCommentsOpts) ([]Comment, error)

	RatePost(userId, postId utils.UUID, rating int) error
	RateComment(userId, commentId utils.UUID, rating int) error

	Comment(postId utils.UUID, comment CommentModel) (utils.UUID, error)
	DeleteComment(userId, commentId utils.UUID) error
	UpdateComment(commentId utils.UUID, data CommentModel) error
}

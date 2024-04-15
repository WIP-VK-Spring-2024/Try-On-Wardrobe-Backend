package feed

import (
	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/feed"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mailru/easyjson"
)

type FeedHandler struct {
	feed domain.FeedRepository
}

func New(db *pgxpool.Pool) *FeedHandler {
	return &FeedHandler{
		feed: feed.New(db),
	}
}

func (h *FeedHandler) GetPosts(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var opts domain.GetPostsOpts
	if err := ctx.QueryParser(&opts); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	opts.RequestingUserID = session.UserID

	posts, err := h.feed.GetPosts(opts)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(posts)
}

func (h *FeedHandler) GetComments(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var opts domain.GetCommentsOpts
	if err := ctx.QueryParser(&opts); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	opts.RequestingUserID = session.UserID

	comments, err := h.feed.GetComments(opts)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(comments)
}

func (h *FeedHandler) CreateComment(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	postId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrPostIdInvalid
	}

	var comment domain.CommentModel
	if err = easyjson.Unmarshal(ctx.Body(), &comment); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	comment.UserID = session.UserID

	err = h.feed.Comment(postId, comment)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

//easyjson:json
type rateRequest struct {
	Rating int
}

func (h *FeedHandler) RatePost(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	postId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrPostIdInvalid
	}

	var req rateRequest
	if err = easyjson.Unmarshal(ctx.Body(), &req); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	if req.Rating > 1 {
		req.Rating = 1
	}
	if req.Rating < -1 {
		req.Rating = -1
	}

	err = h.feed.RatePost(session.UserID, postId, req.Rating)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *FeedHandler) RateComment(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	commentId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrPostIdInvalid
	}

	var req rateRequest
	if err = easyjson.Unmarshal(ctx.Body(), &req); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	if req.Rating > 1 {
		req.Rating = 1
	}
	if req.Rating < -1 {
		req.Rating = -1
	}

	err = h.feed.RateComment(session.UserID, commentId, req.Rating)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

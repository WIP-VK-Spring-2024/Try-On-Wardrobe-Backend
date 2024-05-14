package feed

import (
	"strconv"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/mailru/easyjson"
)

type FeedHandler struct {
	feed   domain.FeedRepository
	recsys domain.Recsys

	getPosts      fiber.Handler
	getLikedPosts fiber.Handler
	getSubPosts   fiber.Handler
}

func New(feed domain.FeedRepository, recsys domain.Recsys) *FeedHandler {
	return &FeedHandler{
		feed:          feed,
		recsys:        recsys,
		getPosts:      getPostsTemplate(feed.GetPosts),
		getLikedPosts: getPostsTemplate(feed.GetLikedPosts),
		getSubPosts:   getPostsTemplate(feed.GetSubscriptionPosts),
	}
}

func (h *FeedHandler) GetPosts(ctx *fiber.Ctx) error {
	return h.getPosts(ctx)
}

func (h *FeedHandler) GetLikedPosts(ctx *fiber.Ctx) error {
	return h.getLikedPosts(ctx)
}

func (h *FeedHandler) GetSubscriptionPosts(ctx *fiber.Ctx) error {
	return h.getSubPosts(ctx)
}

func getPostsTemplate(getter func(domain.GetPostsOpts) ([]domain.Post, error)) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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

		posts, err := getter(opts)
		if err != nil {
			return app_errors.New(err)
		}

		return ctx.JSON(posts)
	}
}

func (h FeedHandler) GetPostsByUser(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	userId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserIdInvalid
	}

	var opts domain.GetPostsOpts
	if err := ctx.QueryParser(&opts); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	opts.RequestingUserID = session.UserID

	posts, err := h.feed.GetPostsByUser(userId, opts)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(posts)
}

func (h FeedHandler) GetRecommendedPosts(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	limit, _ := strconv.Atoi(ctx.Query("limit"))
	if limit == 0 {
		limit = 9
	}

	samplesAmount, _ := strconv.Atoi(ctx.Query("sample_amount"))
	if samplesAmount == 0 {
		samplesAmount = 100
	}

	posts, err := h.recsys.GetRecommendations(ctx.UserContext(), limit, domain.RecsysRequest{
		UserID:        session.UserID,
		SamplesAmount: samplesAmount,
	})
	if err == app_errors.ErrModelUnavailable {
		return err
	}
	if err != nil {
		return app_errors.New(err)
	}

	if posts == nil {
		posts = []domain.Post{}
	}

	return ctx.JSON(posts)
}

func (h *FeedHandler) Subscribe(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	userId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserIdInvalid
	}

	if session.UserID == userId {
		return app_errors.ErrSubscribeTarget
	}

	err = h.feed.Subscribe(session.UserID, userId)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *FeedHandler) Unsubscribe(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	userId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserIdInvalid
	}

	if session.UserID == userId {
		return app_errors.ErrUnsubscribeTarget
	}

	err = h.feed.Unsubscribe(session.UserID, userId)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

type getCommentsRequest struct {
	domain.GetCommentsOpts
	Tree bool `query:"tree"`
}

func (h *FeedHandler) GetComments(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var opts getCommentsRequest
	if err := ctx.QueryParser(&opts); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	postId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrPostIdInvalid
	}

	opts.RequestingUserID = session.UserID
	opts.PostID = postId

	var comments []domain.Comment

	if opts.Tree {
		comments, err = h.feed.GetCommentsTree(opts.GetCommentsOpts)
	} else {
		comments, err = h.feed.GetComments(opts.GetCommentsOpts)
	}
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(comments)
}

//easyjson:json
type createCommentResponse struct {
	Uuid utils.UUID
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

	id, err := h.feed.Comment(postId, comment)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(&createCommentResponse{
		Uuid: id,
	})
}

func (h *FeedHandler) DeleteComment(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	commentId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrCommentIdInvalid
	}

	err = h.feed.DeleteComment(session.UserID, commentId)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h *FeedHandler) UpdateComment(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	commentId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrCommentIdInvalid
	}

	var comment domain.CommentModel
	if err = easyjson.Unmarshal(ctx.Body(), &comment); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}
	comment.UserID = session.UserID

	err = h.feed.UpdateComment(commentId, comment)
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

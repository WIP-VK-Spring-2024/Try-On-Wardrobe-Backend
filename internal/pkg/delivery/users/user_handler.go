package users

import (
	"net/http"
	"strings"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	userRepo "try-on/internal/pkg/repository/sqlc/users"
	"try-on/internal/pkg/usecase/session"
	"try-on/internal/pkg/usecase/users"
	"try-on/internal/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	easyjson "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	users    domain.UserUsecase
	sessions domain.SessionUsecase
	file     domain.FileManager
	cfg      *config.Static
}

func New(
	db *pgxpool.Pool,
	fileManager domain.FileManager,
	sessionCfg *config.Session,
	cfg *config.Static,
) *UserHandler {
	userRepo := userRepo.New(db)

	return &UserHandler{
		users:    users.New(userRepo),
		sessions: session.New(userRepo, sessionCfg),
		file:     fileManager,
		cfg:      cfg,
	}
}

//easyjson:json
type tokenResponse struct {
	Token  string
	UserID utils.UUID
}

func (h *UserHandler) Create(ctx *fiber.Ctx) error {
	var user domain.User
	if err := easyjson.Unmarshal(ctx.Body(), &user); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	err := h.users.Create(&user)
	if err != nil {
		return app_errors.New(err)
	}

	token, err := h.sessions.IssueToken(user.ID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(tokenResponse{
		Token:  token,
		UserID: user.ID,
	})
}

func (h *UserHandler) Update(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	userId, err := utils.ParseUUID(ctx.Params("id"))
	if err != nil {
		return app_errors.ErrUserIdInvalid
	}

	if userId != session.UserID {
		return app_errors.ErrNotOwner
	}

	var fileName string

	fileHeader, err := ctx.FormFile("img")
	switch {
	case err == nil && fileHeader != nil:
		break
	case err != fasthttp.ErrMissingFile:
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	default:
		fileName = userId.String() + "_" + fileHeader.Filename
	}

	var user domain.User
	if err := ctx.BodyParser(&user); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}
	user.ID = userId
	user.Avatar = fileName

	err = h.users.Update(user)
	if err != nil {
		return app_errors.New(err)
	}

	if fileName == "" {
		return ctx.SendString(common.EmptyJson)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return app_errors.New(err)
	}
	defer file.Close()

	err = h.file.Save(ctx.UserContext(), h.cfg.Avatars, fileName, file)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.SendString(common.EmptyJson)
}

func (h UserHandler) SearchUsers(ctx *fiber.Ctx) error {
	name := strings.TrimSpace(ctx.Query("name"))
	if name == "" {
		return app_errors.ResponseError{
			Code: http.StatusBadRequest,
			Msg:  "query param 'name' should be non-empty",
		}
	}

	users, err := h.users.SearchUsers(name)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(users)
}

func (h UserHandler) GetSubscriptions(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	users, err := h.users.GetSubscriptions(session.UserID)
	if err != nil {
		return app_errors.New(err)
	}

	return ctx.JSON(users)
}

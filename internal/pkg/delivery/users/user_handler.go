package users

import (
	"context"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/common"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	userRepo "try-on/internal/pkg/repository/sqlc/users"
	"try-on/internal/pkg/usecase/session"
	userImagesUsecase "try-on/internal/pkg/usecase/user_images"
	"try-on/internal/pkg/usecase/users"
	"try-on/internal/pkg/utils"
	"try-on/internal/pkg/utils/validate"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	easyjson "github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	users      domain.UserUsecase
	userImages domain.UserImageUsecase
	sessions   domain.SessionUsecase
	file       domain.FileManager
	cfg        *config.Static
}

func New(
	db *pgxpool.Pool,
	fileManager domain.FileManager,
	sessionCfg *config.Session,
	cfg *config.Static,
) *UserHandler {
	userRepo := userRepo.New(db)

	return &UserHandler{
		users:      users.New(userRepo),
		userImages: userImagesUsecase.New(db),
		sessions:   session.New(userRepo, sessionCfg),
		file:       fileManager,
		cfg:        cfg,
	}
}

//easyjson:json
type registerResponse struct {
	Token    string
	UserName string
	UserID   utils.UUID
	Email    string
	Gender   domain.Gender
	Privacy  domain.Privacy
}

func (h *UserHandler) Create(ctx *fiber.Ctx) error {
	var user domain.User
	if err := easyjson.Unmarshal(ctx.Body(), &user); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}

	err := validate.Struct(&user)
	if err != nil {
		return app_errors.ValidationError(err)
	}

	err = h.users.Create(&user)
	if err != nil {
		return app_errors.New(err)
	}

	token, err := h.sessions.IssueToken(user.ID)
	if err != nil {
		return app_errors.New(err)
	}

	err = h.createDefaultPhoto(&user, ctx.UserContext())
	if err != nil {
		middleware.LogWarning(ctx, err)
	}
	return ctx.JSON(registerResponse{
		Token:    token,
		UserID:   user.ID,
		UserName: user.Name,
		Gender:   user.Gender,
		Privacy:  user.Privacy,
		Email:    user.Email,
	})
}

func (h *UserHandler) createDefaultPhoto(user *domain.User, ctx context.Context) error {
	defaultImg, err := h.file.Get(ctx, h.cfg.FullBody, h.cfg.DefaultImgPaths[user.Gender])
	if err != nil {
		return err
	}
	defer defaultImg.Close()

	img := &domain.UserImage{
		UserID: user.ID,
		Image:  h.cfg.FullBody,
	}

	err = h.userImages.Create(img)
	if err != nil {
		return err
	}

	err = h.file.Save(ctx, h.cfg.FullBody, img.ID.String(), defaultImg)
	if err != nil {
		return err
	}

	return nil
}

//easyjson:json
type updateResponse struct {
	Avatar string
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
	case err != nil && err != fasthttp.ErrMissingFile:
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	case fileHeader == nil:
		break
	default:
		fileName = userId.String()
	}

	var user domain.User
	if err := ctx.BodyParser(&user); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}
	user.ID = userId
	user.Avatar = h.cfg.Avatars + "/" + fileName

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

	return ctx.JSON(updateResponse{
		Avatar: h.cfg.Avatars,
	})
}

func (h UserHandler) SearchUsers(ctx *fiber.Ctx) error {
	session := middleware.Session(ctx)
	if session == nil {
		return app_errors.ErrUnauthorized
	}

	var opts domain.SearchUserOpts
	if err := ctx.QueryParser(&opts); err != nil {
		middleware.LogWarning(ctx, err)
		return app_errors.ErrBadRequest
	}
	opts.UserID = session.UserID

	users, err := h.users.SearchUsers(opts)
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

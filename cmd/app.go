package main

import (
	"errors"
	"log"
	"net/http"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	clothes "try-on/internal/pkg/clothes/delivery"
	clothesRepo "try-on/internal/pkg/clothes/repository"
	clothesUsecase "try-on/internal/pkg/clothes/usecase"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/ml"
	session "try-on/internal/pkg/session/delivery"
	tryOn "try-on/internal/pkg/try-on/delivery"
	"try-on/internal/pkg/utils"

	"github.com/wagslane/go-rabbitmq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	api    *fiber.App
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func (app *App) Run() error {
	db, err := app.getDB()
	if err != nil {
		return err
	}

	err = applyMigrations(app.cfg.Sql, db)
	if err != nil {
		return err
	}

	log.Println("Connecting to rabbit", app.cfg.Rabbit.DSN())
	rabbitConn, err := rabbitmq.NewConn(app.cfg.Rabbit.DSN())
	if err != nil {
		return err
	}

	clothesProcessor, err := ml.New(
		app.cfg.Rabbit.RequestQueue,
		app.cfg.Rabbit.ResponseQueue,
		rabbitConn,
	)
	if err != nil {
		return err
	}

	recover := recover.New(recover.Config{
		EnableStackTrace: true,
	})

	logger := logger.New(logger.Config{
		Format:     config.JsonLogFormat,
		TimeFormat: config.TimeFormat,
	})

	cors := cors.New(cors.Config{
		AllowOrigins:     app.cfg.Cors.Domain,
		AllowCredentials: app.cfg.Cors.AllowCredentials,
		MaxAge:           app.cfg.Cors.MaxAge,
	})

	sessionHandler := session.New(db, &app.cfg.Session)

	checkSession := middleware.CheckSession(middleware.SessionConfig{
		TokenName:    app.cfg.Session.TokenName,
		Sessions:     sessionHandler.Sessions,
		NoAuthRoutes: []string{"/register", "/login"},
		// SecureRoutes: []string{"/renew", "/clothes"},
	})

	clothesUsecase := clothesUsecase.New(clothesRepo.New(db))

	clothesHandler := clothes.New(clothesUsecase, clothesProcessor, &app.cfg.Static)

	tryOnHandler := tryOn.New(db, clothesProcessor, clothesUsecase, app.logger, &app.cfg.Static)

	app.api.Use(recover, logger, cors, middleware.AddLogger(app.logger), checkSession)

	app.api.Post("/register", sessionHandler.Register)
	app.api.Post("/login", sessionHandler.Login)
	app.api.Post("/renew", sessionHandler.Renew)

	app.api.Post("/clothes", clothesHandler.Upload)
	app.api.Get("/clothes", clothesHandler.GetOwn)
	app.api.Get("/clothes/:id", clothesHandler.GetByID)
	app.api.Delete("/clothes/:id", clothesHandler.Delete)
	app.api.Put("/clothes/:id", clothesHandler.Update)
	app.api.Get("/user/:id/clothes", clothesHandler.GetByUser)

	app.api.Post("/user/try-on/:clothing_id", tryOnHandler.TryOn)
	app.api.Get("/user/try-on/:clothing_id", tryOnHandler.GetTryOnResult)

	app.api.Static("/static", app.cfg.Static.Dir)

	tryOnHandler.ListenTryOnResults()

	return app.api.Listen(app.cfg.Addr)
}

func NewApp(cfg *config.Config, logger *zap.SugaredLogger) *App {
	return &App{
		api: fiber.New(
			fiber.Config{
				ErrorHandler: errorHandler,
				JSONEncoder:  utils.EasyJsonMarshal,
				JSONDecoder:  utils.EasyJsonUnmarshal,
			},
		),
		cfg:    cfg,
		logger: logger,
	}
}

func (app *App) getDB() (*gorm.DB, error) {
	pg, err := initPostgres(&app.cfg.Postgres)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: pg,
	}), &gorm.Config{
		// Logger: gormLogger.Discard,
		// FullSaveAssociations: true,
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	var e *fiber.Error
	if errors.As(err, &e) {
		return ctx.Status(e.Code).JSON(
			&app_errors.ResponseError{
				Msg: e.Message,
			},
		)
	}

	msg := "Internal Server Error"

	var errorMsg *app_errors.ResponseError
	if errors.As(err, &errorMsg) {
		return ctx.Status(errorMsg.Code).JSON(errorMsg)
	}

	middleware.LogError(ctx, err)

	return ctx.Status(http.StatusInternalServerError).JSON(
		&app_errors.ResponseError{
			Msg: msg,
		},
	)
}

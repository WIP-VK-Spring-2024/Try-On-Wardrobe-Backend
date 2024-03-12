package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/config"
	session "try-on/internal/pkg/session/delivery"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mailru/easyjson"
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

	err = app.registerRoutes(db)
	if err != nil {
		return err
	}
	return app.api.Listen(app.cfg.Addr)
}

func NewApp(cfg *config.Config, logger *zap.SugaredLogger) *App {
	return &App{
		api: fiber.New(
			fiber.Config{
				ErrorHandler: errorHandler,
				JSONEncoder:  easyjsonMarshal,
				JSONDecoder:  easyjsonUnmarshal,
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
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *App) registerRoutes(db *gorm.DB) error {
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

	sessionHandler := session.NewSessionHandler(db, &app.cfg.Session)

	addSession := middleware.CheckSession(middleware.SessionConfig{
		TokenName:    app.cfg.Session.TokenName,
		Sessions:     sessionHandler.Sessions,
		NoAuthRoutes: []string{"/register", "/login"},
		SecureRoutes: []string{"/auth-only"},
	})

	app.api.Use(recover, logger, cors, middleware.AddLogger(app.logger), addSession)

	app.api.Post("/register", sessionHandler.Register)
	app.api.Post("/login", sessionHandler.Login)

	app.api.Get("/check", func(ctx *fiber.Ctx) error {
		session := middleware.Session(ctx)
		if session == nil {
			return ctx.SendString("Not logged in\n")
		}
		return ctx.SendString("Logged in\n")
	})

	app.api.Get("/auth-only", func(ctx *fiber.Ctx) error {
		session := middleware.Session(ctx)
		return ctx.SendString(fmt.Sprintln("Hello,", session.UserID, '!'))
	})

	app.api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, world!\n")
	})

	return nil
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		code = fiberError.Code
		msg = err.Error()
	}

	if code == fiber.StatusInternalServerError {
		middleware.LogError(ctx, err)
	}

	return ctx.Status(code).JSON(
		&app_errors.ErrorMsg{
			Msg: msg,
		},
	)
}

func easyjsonMarshal(value interface{}) ([]byte, error) {
	marshaler, ok := value.(easyjson.Marshaler)
	if ok {
		return easyjson.Marshal(marshaler)
	}
	return json.Marshal(value)
}

func easyjsonUnmarshal(data []byte, value interface{}) error {
	unmarshaler, ok := value.(easyjson.Unmarshaler)
	if ok {
		return easyjson.Unmarshal(data, unmarshaler)
	}
	return json.Unmarshal(data, value)
}

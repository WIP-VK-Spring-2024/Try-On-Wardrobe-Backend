package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"
	session "try-on/internal/pkg/session/delivery"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
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

	err = applyMigrations(app.cfg.SqlDir, db)
	if err != nil {
		return err
	}

	app.registerMiddleware()

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
	dsn := app.cfg.Postgres.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Discard,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *App) registerRoutes(db *gorm.DB) error {
	sessionHandler := session.NewSessionHandler(db, session.Config{
		Session: app.cfg.Session,
		Redis:   app.cfg.Redis,
	})

	authRoutes := fiber.New()
	// authRoutes.Use(middleware.AddSession(middleware.SessionConfig{
	// 	CookieName: app.cfg.Session.CookieName,
	// 	Sessions: ,
	// }), middleware.CheckSession)

	app.api.Post("/register", sessionHandler.Register)
	app.api.Post("/login", sessionHandler.Login)
	authRoutes.Post("/logout", sessionHandler.Logout)

	app.api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, world!\n")
	})

	app.api.Mount("/", authRoutes)

	return nil
}

func (app *App) registerMiddleware() {
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

	app.api.Use(recover, logger, cors, middleware.AddLogger(app.logger))
}

func applyMigrations(scriptsDir string, db *gorm.DB) error {
	files, err := os.ReadDir(scriptsDir)
	if err != nil {
		return err
	}

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(scriptsDir + "/" + fileInfo.Name())
		if err != nil {
			return err
		}

		bytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		err = db.Exec(string(bytes)).Error
		if err != nil {
			return errors.Join(fmt.Errorf("failed applying migration '%s'", fileInfo.Name()), err)
		}
	}

	return db.AutoMigrate(
		&domain.User{},
		&domain.Clothes{},
		&domain.Tag{},
		&domain.Style{},
		&domain.Type{},
		&domain.Subtype{},
	)
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

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	clothes "try-on/internal/pkg/clothes/delivery"
	"try-on/internal/pkg/config"
	"try-on/internal/pkg/file_manager/filesystem"
	"try-on/internal/pkg/ml"
	session "try-on/internal/pkg/session/delivery"

	amqp "github.com/rabbitmq/amqp091-go"

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

	log.Println("Connecting to rabbit", app.cfg.Rabbit.Addr())
	rabbitConn, err := amqp.Dial(app.cfg.Rabbit.Addr())
	if err != nil {
		return err
	}
	defer rabbitConn.Close()

	rabbitChan, err := rabbitConn.Channel()
	if err != nil {
		return err
	}
	defer rabbitChan.Close()

	err = app.registerRoutes(db, rabbitChan)
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
		FullSaveAssociations: true,
		TranslateError:       true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *App) registerRoutes(db *gorm.DB, rabbitChan *amqp.Channel) error {
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
		SecureRoutes: []string{"/renew", "/clothes", "/user"},
	})

	clothesProcessor, err := ml.New(app.cfg.Rabbit.RequestQueue, rabbitChan)
	if err != nil {
		return err
	}

	clothesHandler := clothes.New(db, filesystem.New(app.cfg.ImageDir), clothesProcessor)

	app.api.Use(recover, logger, cors, middleware.AddLogger(app.logger), checkSession)

	app.api.Post("/register", sessionHandler.Register)
	app.api.Post("/login", sessionHandler.Login)
	app.api.Post("/renew", sessionHandler.Renew)

	app.api.Post("/clothes", clothesHandler.Upload)
	app.api.Get("/clothes/:id", clothesHandler.GetByID)
	app.api.Get("/user/:id/clothes", clothesHandler.GetByUser)

	return nil
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	msg := "Internal Server Error"

	var errorMsg app_errors.ErrorMsg
	if errors.As(err, &errorMsg) {
		return ctx.Status(errorMsg.Code).JSON(errorMsg)
	}

	middleware.LogError(ctx, err)

	return ctx.Status(http.StatusInternalServerError).JSON(
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

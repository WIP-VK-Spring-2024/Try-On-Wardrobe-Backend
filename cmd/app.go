package main

import (
	"try-on/internal/pkg/config"
	session "try-on/internal/pkg/session/delivery"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	api *fiber.App
	cfg *config.Config
}

func (app *App) Run() error {
	app.registerMiddleware()
	app.registerRoutes()
	return app.api.Listen(app.cfg.Port)
}

func NewApp(cfg *config.Config) *App {
	return &App{
		api: fiber.New(),
		cfg: cfg,
	}
}

func (app *App) getDB() (*gorm.DB, error) {
	dsn := app.cfg.Postgres.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *App) registerRoutes() error {
	db, err := app.getDB()
	if err != nil {
		return err
	}

	sessionHandler := session.NewSessionHandler(db, session.Config{
		Session: app.cfg.Session,
		Redis:   app.cfg.Redis,
	})

	authRoutes := fiber.New()
	// authRoutes.Use(middleware.AddSession(middleware.SessionConfig{
	// 	CookieName: app.cfg.Session.CookieName,
	// 	Sessions: ,
	// }), middleware.CheckSession)

	app.api.Post("/login", sessionHandler.Login)
	authRoutes.Post("/logout", sessionHandler.Logout)
	app.api.Post("/register", sessionHandler.Register)

	app.api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, world!\n")
	})

	app.api.Mount("/", authRoutes)

	return nil
}

func (app *App) registerMiddleware() {
	api := fiber.New()

	recover := recover.New(recover.Config{
		EnableStackTrace: true,
	})

	logger := logger.New(logger.Config{
		Format: config.JsonLogFormat,
	})

	cors := cors.New(cors.Config{
		AllowOrigins:     app.cfg.Cors.Domain,
		AllowCredentials: app.cfg.Cors.AllowCredentials,
		MaxAge:           app.cfg.Cors.MaxAge,
	})

	api.Use(recover, logger, cors)
}

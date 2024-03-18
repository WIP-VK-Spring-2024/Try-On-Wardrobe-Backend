package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"try-on/internal/pkg/config"

	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

func applyMigrations(cfg config.Sql, db *gorm.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: cfg.Dir,
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	n, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		errors.Join(errors.New("sql-migrate migrations failed"), err)
	}
	fmt.Printf("Applied %d migrations\n", n)

	return nil
}

func initPostgres(config *config.Postgres) (*sql.DB, error) {
	till := time.Now().Add(time.Second * config.InitTimeout)
	log.Println("Connecting to postgres:", config.DSN())

	db, err := sql.Open("pgx", config.DSN())
	if err != nil {
		return nil, err
	}

	for time.Now().Before(till) {
		log.Println("Trying to open pg connection")

		err = db.Ping()
		if err == nil {
			log.Println("pg connection successfully opened")
			break
		}

		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, errors.New("connection to postgres timed out")
	}

	db.SetMaxOpenConns(config.MaxConn)
	return db, nil
}

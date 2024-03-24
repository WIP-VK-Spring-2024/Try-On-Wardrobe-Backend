package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"try-on/internal/pkg/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

func applyMigrations(cfg config.Sql, pg *pgxpool.Pool) error {
	migrations := &migrate.FileMigrationSource{
		Dir: cfg.Dir,
	}

	sqlDB := stdlib.OpenDBFromPool(pg)

	n, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Join(errors.New("sql-migrate migrations failed"), err)
	}
	fmt.Printf("Applied %d migrations\n", n)

	return nil
}

func initPostgres(config *config.Postgres) (*pgxpool.Pool, error) {
	till := time.Now().Add(time.Second * config.InitTimeout)
	log.Println("Connecting to postgres:", config.DSN())

	cfg, err := pgxpool.ParseConfig(config.DSN())
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	for time.Now().Before(till) {
		log.Println("Trying to open pg connection")

		err = db.Ping(context.Background())
		if err == nil {
			log.Println("pg connection successfully opened")
			break
		}

		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, errors.New("connection to postgres timed out")
	}

	conn, err := db.Acquire(context.Background())
	if err != nil {
		return nil, errors.Join(errors.New("failed acquiring connection"), err)
	}

	t, err := conn.Conn().LoadType(context.Background(), "season")
	if err != nil {
		return nil, errors.Join(errors.New("failed registering season type"), err)
	}
	conn.Conn().TypeMap().RegisterType(t)

	t, err = conn.Conn().LoadType(context.Background(), "_seasons")
	if err != nil {
		return nil, errors.Join(errors.New("failed registering season[] type"), err)
	}
	conn.Conn().TypeMap().RegisterType(t)

	return db, nil
}

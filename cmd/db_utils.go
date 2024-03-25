package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
)

func applyMigrations(cfg config.Sql, pgCfg *config.Postgres) error {
	migrations := &migrate.FileMigrationSource{
		Dir: cfg.Dir,
	}

	db, err := sql.Open("pgx", pgCfg.DSN())
	if err != nil {
		return errors.Join(errors.New("failed opening connections for migrations"), err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Join(errors.New("sql-migrate migrations failed"), err)
	}
	fmt.Printf("Applied %d migrations\n", n)

	return nil
}

var customTypes = []string{"season", "season[]"}

func initPostgres(config *config.Postgres) (*pgxpool.Pool, error) {
	till := time.Now().Add(time.Second * config.InitTimeout)
	log.Println("Connecting to postgres:", config.DSN())

	cfg, err := pgxpool.ParseConfig(config.PoolDSN())
	if err != nil {
		return nil, err
	}

	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterDefaultPgType(domain.Spring, "season")

		for _, customType := range customTypes {
			t, err := conn.LoadType(context.Background(), customType)
			if err != nil {
				return errors.Join(fmt.Errorf("failed registering type %s", customType), err)
			}
			conn.TypeMap().RegisterType(t)
		}
		return nil
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

	return db, nil
}

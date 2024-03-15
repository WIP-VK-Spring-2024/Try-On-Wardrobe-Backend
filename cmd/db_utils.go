package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"try-on/internal/pkg/config"
	"try-on/internal/pkg/domain"

	"gorm.io/gorm"
)

func applyMigrations(cfg config.Sql, db *gorm.DB) error {
	err := execMultipleScripts(db, cfg.Dir+"/", cfg.BeforeGorm)
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.ClothesModel{},
		&domain.Tag{},
		&domain.Style{},
		&domain.Type{},
		&domain.Subtype{},
	)
	if err != nil {
		return errors.Join(errors.New("gorm migrations failed"), err)
	}

	return execMultipleScripts(db, cfg.Dir+"/", cfg.AfterGorm)
}

func execMultipleScripts(db *gorm.DB, prefix string, paths []string) error {
	for _, fileName := range paths {
		err := execSqlScript(db, prefix+fileName)
		if err != nil {
			return errors.Join(fmt.Errorf("failed applying migration '%s'", fileName), err)
		}
	}
	return nil
}

func execSqlScript(db *gorm.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return db.Exec(string(bytes)).Error
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
			log.Println("Ping sucessful")
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

package db

import (
	"database/sql"

	"github.com/pressly/goose/v3"

	"github.com/romanyakovlev/gophermart/internal/logger"
)

const migrationsDir = "./migrations"

func InitDB(DatabaseURI string, sugar *logger.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", DatabaseURI)
	if err != nil {
		sugar.Errorf("Server error: %v", err)
		return nil, err
	}

	if DatabaseURI != "" {
		if err := goose.Up(db, migrationsDir); err != nil {
			sugar.Fatalf("goose Up failed: %v\n", err)
		}
	}

	return db, nil
}

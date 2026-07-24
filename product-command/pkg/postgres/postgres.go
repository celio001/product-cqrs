package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/celio001/product-command/config"
)

func ConectPostgres() (*sql.DB, error) {
	dsn := config.GetString("POSTGRES_DB_DSN")
	if dsn == "" {
		return nil, errors.New("Empty Postgres DSN")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

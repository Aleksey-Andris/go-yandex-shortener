package postgresstorage

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	linkTable   = "ys_link"
	shortURL    = "short_url"
	originalURL = "original_url"
)

func NewPostgresDB(cfg string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := initTable(db); err != nil {
		return nil, err
	}
	return db, nil
}

func initTable(db *sqlx.DB) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, %s VARCHAR(255) NOT NULL UNIQUE, %s VARCHAR(255) NOT NULL);", linkTable, shortURL, originalURL)
	_, err := db.Exec(query)
	return err
}

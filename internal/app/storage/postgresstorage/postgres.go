package postgresstorage

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB(configDB string) (*sql.DB, error) {
	db, err := sql.Open("pgx", configDB)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

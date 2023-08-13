package postgresstorage

import (
	"fmt"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/jmoiron/sqlx"
)

type linkStorage struct {
	db *sqlx.DB
}

func NewLinkStorage(db *sqlx.DB) (*linkStorage, error) {
	s := &linkStorage{db: db}
	return s, nil
}

func (s *linkStorage) GetOneByIdent(ident string) (domain.Link, error) {
	var link domain.Link
	query := fmt.Sprintf("SELECT id, %s, %s FROM %s WHERE %s = $1;", shortURL, originalURL, linkTable, shortURL)
	err := s.db.Get(&link, query, ident)
	return link, err
}

func (s *linkStorage) Create(ident, fulLink string) (domain.Link, error) {
	var link domain.Link
	query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES($1, $2) RETURNING id, %s, %s;",
		linkTable, shortURL, originalURL, shortURL, originalURL)
	err := s.db.Get(&link, query, ident, fulLink)
	return link, err
}

func (s *linkStorage) Close() error {
	return s.db.Close()
}

func (s *linkStorage) GetMaxID() (int32, error) {
	var maxID int32
	query := fmt.Sprintf("SELECT COALESCE(MAX(id), 1) FROM %s;", linkTable)
	row := s.db.QueryRow(query)
	if err := row.Scan(&maxID); err != nil {
		return 0, err
	}
	return maxID, nil
}

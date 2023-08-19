package postgresstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type linkStorage struct {
	db *sqlx.DB
}

func NewLinkStorage(db *sqlx.DB) (*linkStorage, error) {
	s := &linkStorage{db: db}
	return s, nil
}

func (s *linkStorage) GetOneByIdent(ctx context.Context, ident string) (domain.Link, error) {
	var link domain.Link
	query := fmt.Sprintf("SELECT id, %s, %s FROM %s WHERE %s = $1;", shortURL, originalURL, linkTable, shortURL)
	err := s.db.GetContext(ctx, &link, query, ident)
	return link, err
}

func (s *linkStorage) Create(ctx context.Context, ident, fulLink string) (domain.Link, error) {
	var link domain.Link
	query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES($1, $2) RETURNING id, %s, %s;",
		linkTable, shortURL, originalURL, shortURL, originalURL)
	err := s.db.GetContext(ctx, &link, query, ident, fulLink)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
	}

	query = fmt.Sprintf("SELECT id, %s, %s FROM %s WHERE %s = $1;", shortURL, originalURL, linkTable, originalURL)
	if err := s.db.GetContext(ctx, &link, query, fulLink); err != nil {
		return link, err
	}
	return link, err
}

func (s *linkStorage) CreateLinks(ctx context.Context, links []domain.Link) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES($1, $2);",
		linkTable, shortURL, originalURL)
	stm, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	for _, v := range links {
		_, err := stm.ExecContext(ctx, v.Ident, v.FulLink)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *linkStorage) Close() error {
	return s.db.Close()
}

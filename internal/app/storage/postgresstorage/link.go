package postgresstorage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
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
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1;", linkTable, shortURL)
	err := s.db.GetContext(ctx, &link, query, ident)
	return link, err
}

func (s *linkStorage) Create(ctx context.Context, ident, fulLink string, userID int32) (domain.Link, error) {
	var link domain.Link

	query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES($1, $2, $3) RETURNING id, %s, %s, %s;",
		linkTable, shortURL, originalURL, userIDStor, shortURL, originalURL, userIDStor)
	err := s.db.GetContext(ctx, &link, query, ident, fulLink, userID)

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

func (s *linkStorage) CreateLinks(ctx context.Context, links []domain.Link, userID int32) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES($1, $2, $3);",
		linkTable, shortURL, originalURL, userIDStor)
	stm, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	for _, v := range links {
		_, err := stm.ExecContext(ctx, v.Ident, v.FulLink, userID)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *linkStorage) GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error) {
	var linkListByUserIDRes []dto.LinkListByUserIDRes
	query := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s = $1 AND %s = $2;", shortURL, originalURL, linkTable, userIDStor, isDeleted)
	err := s.db.SelectContext(ctx, &linkListByUserIDRes, query, userID, false)
	return linkListByUserIDRes, err
}

func (s *linkStorage) DeleteByIdents(ctx context.Context, idents ...string) error {
	var values []string
	var args []any
	for i, v := range idents {
		params := fmt.Sprintf("$%d", i+1)
		values = append(values, params)
		args = append(args, v)
	}
	query := fmt.Sprintf("UPDATE %s SET %s = true WHERE %s IN (", linkTable, isDeleted, shortURL) + strings.Join(values, ",") + ");"
	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func (s *linkStorage) GetByIdents(ctx context.Context, idents ...string) ([]domain.Link, error) {
	var values []string
	var args []any
	for i, v := range idents {
		params := fmt.Sprintf("$%d", i+1)
		values = append(values, params)
		args = append(args, v)
	}
	var links []domain.Link
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s IN (", linkTable, shortURL) + strings.Join(values, ",") + ");"
	err := s.db.SelectContext(ctx, &links, query, args...)
	return links, err
}

func (s *linkStorage) Close() error {
	return s.db.Close()
}

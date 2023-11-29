package postgresstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
)

type userStorage struct {
	db *sqlx.DB
}

// NewUserStorage - constructor for userStorage.
func NewUserStorage(db *sqlx.DB) (*userStorage, error) {
	s := &userStorage{db: db}
	return s, nil
}

// CreateUser - user creating method.
func (s *userStorage) CreateUser(ctx context.Context) (int32, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()
	var user domain.User
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1) RETURNING id;", userTable, createDate)
	err = s.db.GetContext(ctx, &user, query, time.Now())
	if err != nil {
		return -1, err
	}
	if err := tx.Commit(); err != nil {
		return -1, err
	}
	return user.ID, err
}

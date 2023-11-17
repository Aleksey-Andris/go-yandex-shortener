package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	signingKey = "ushjdhui38487"
	tokenExp   = time.Hour * 3
)

type UserStorage interface {
	CreateUser(ctx context.Context) (int32, error)
}

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID int32
}

type authService struct {
	storage UserStorage
}

func NewAauthService(userStorage UserStorage) *authService {
	return &authService{
		storage: userStorage,
	}
}

func (s *authService) ParseToken(tokenString string) (int32, bool, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(signingKey), nil
		})

	if err != nil {
		return 0, false, err
	}

	return claims.UserID, token.Valid, err
}

func (s *authService) BuildJWTString(userID int32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) CreateUser(ctx context.Context) (int32, error) {
	return s.storage.CreateUser(ctx)
}

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

// UserStorage - interface describing the contract for working with users entities in storage.
type UserStorage interface {
	// CreateUser - user creating method.
	CreateUser(ctx context.Context) (int32, error)
}

// TokenClaims - structure containing information to create a user token.
type TokenClaims struct {
	jwt.RegisteredClaims
	// UserID - the unique ident of an entity is a surrogate key.
	UserID int32
}

// authService - structure representing a usecase for user.
type authService struct {
	storage UserStorage
}

// NewAauthService - constructor for authService.
func NewAauthService(userStorage UserStorage) *authService {
	return &authService{
		storage: userStorage,
	}
}

// ParseToken - GWT token parsing method, returns user ID, token validity and error.
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

// BuildJWTString - GWT token creating method.
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

// CreateUser - user creating method.
func (s *authService) CreateUser(ctx context.Context) (int32, error) {
	return s.storage.CreateUser(ctx)
}

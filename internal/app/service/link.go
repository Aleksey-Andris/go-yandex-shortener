// The service package contains settings forservice layer(usecases).
package service

import (
	"context"
	"crypto/md5"
	"math/rand"
	"time"

	"github.com/speps/go-hashids"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
)

const (
	salt = "Qw6"
)

// LinkStorage - interface describing the contract for working with link entities in storage.
type LinkStorage interface {
	// GetOneByIdent - returns the link entity short ident.
	GetOneByIdent(ctx context.Context, ident string) (domain.Link, error)
	// Create - create the link entity.
	Create(ctx context.Context, idemt, fulLink string, userID int32) (domain.Link, error)
	// CreateLinks - create the link entityes.
	CreateLinks(ctx context.Context, links []domain.Link, userID int32) error
	// GetLinksByUserID - returns all user's links.
	GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error)
	// DeleteLinksByIdent - removes links by their idents.
	DeleteByIdents(ctx context.Context, idents ...string) error
	// GetByIdents - returns links by their idents.
	GetByIdents(ctx context.Context, idents ...string) ([]domain.Link, error)
	// Close - safely stops the database.
	Close() error
}

// linkService -  structure representing a usecase for links.
type linkService struct {
	storage LinkStorage
}

// NewLinkService - // NewServices - constructor for linkService.
func NewLinkService(storage LinkStorage) *linkService {
	return &linkService{
		storage: storage,
	}
}

// GetIdent - returns the short link by the full URL.
func (s *linkService) GetIdent(ctx context.Context, fulLink string, userID int32) (string, error) {
	ident := s.GenerateIdent(fulLink)
	link, err := s.storage.Create(ctx, ident, fulLink, userID)
	return link.Ident, err
}

	// GetIdents - returns the short links by the full URLs.
func (s *linkService) GetIdents(ctx context.Context, linkReq []dto.LinkListReq, userID int32) ([]dto.LinkListRes, error) {
	result := make([]dto.LinkListRes, 0)
	links := make([]domain.Link, 0)
	for _, v := range linkReq {
		ident := s.GenerateIdent(v.OriginalURL)
		result = append(result, dto.LinkListRes{CorrelationID: v.CorrelationID, ShortURL: ident})
		links = append(links, domain.Link{Ident: ident, FulLink: v.OriginalURL})
	}
	err := s.storage.CreateLinks(ctx, links, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetLinksByUserID - returns all user's links.
func (s *linkService) GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error) {
	return s.storage.GetLinksByUserID(ctx, userID)
}

// GetFulLink - returns the full link by the shortened URL ident.
func (s *linkService) GetFulLink(ctx context.Context, ident string) (domain.Link, error) {
	link, err := s.storage.GetOneByIdent(ctx, ident)
	if err != nil {
		return link, err
	}
	return link, nil
}

// DeleteLinksByIdent - removes links by their idents.
func (s *linkService) DeleteLinksByIdent(ctx context.Context, idents ...string) error {
	return s.storage.DeleteByIdents(ctx, idents...)
}

// CanDelete - checks whether the user has the right to delete these links.
func (s *linkService) CanDelete(ctx context.Context, userID int32, idents ...string) (bool, error) {
	links, err := s.storage.GetByIdents(ctx, idents...)
	if err != nil {
		return false, err
	}
	for _, link := range links {
		if link.UserID != userID {
			return false, err
		}
	}
	return true, nil
}

// GenerateIdent - generates a unique shorts ident.
func (s *linkService) GenerateIdent(url string) string {
	rand.NewSource(time.Now().UnixNano())
	hash := md5.New()
	hash.Write([]byte(url))
	hd := hashids.NewData()
	hd.Salt = string(hash.Sum([]byte(salt)))
	h, _ := hashids.NewWithData(hd)
	ident, _ := h.Encode([]int{rand.Int()})
	return ident
}

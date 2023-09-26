package service

import (
	"context"
	"crypto/md5"
	"math/rand"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/speps/go-hashids"
)

const (
	salt = "Qw6"
)

type LinkStorage interface {
	GetOneByIdent(ctx context.Context, ident string) (domain.Link, error)
	Create(ctx context.Context, idemt, fulLink string, userID int32) (domain.Link, error)
	CreateLinks(ctx context.Context, links []domain.Link, userID int32) error
	GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error)
	Close() error
}

type linkService struct {
	storage LinkStorage
}

func NewLinkService(storage LinkStorage) *linkService {
	return &linkService{
		storage: storage,
	}
}

func (s *linkService) GetIdent(ctx context.Context, fulLink string, userID int32) (string, error) {
	ident := s.GenerateIdent(fulLink)
	link, err := s.storage.Create(ctx, ident, fulLink, userID)
	return link.Ident, err
}

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

func (s *linkService) GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error) {
	return s.storage.GetLinksByUserID(ctx, userID)
}

func (s *linkService) GetFulLink(ctx context.Context, ident string) (string, error) {
	link, err := s.storage.GetOneByIdent(ctx, ident)
	if err != nil {
		return "", err
	}
	return link.FulLink, nil
}

func (s *linkService) GenerateIdent(url string) string {
	hash := md5.New()
	hash.Write([]byte(url))
	hd := hashids.NewData()
	hd.Salt = string(hash.Sum([]byte(salt)))
	h, _ := hashids.NewWithData(hd)
	ident, _ := h.Encode([]int{rand.Int()})
	return ident
}

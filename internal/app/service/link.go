package service

import (
	"context"
	"math/rand"
	"sync/atomic"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/speps/go-hashids"
)

var (
	salt = "Qw6"
)

type LinkStorage interface {
	GetOneByIdent(ctx context.Context, ident string) (domain.Link, error)
	Create(ctx context.Context, idemt, fulLink string) (domain.Link, error)
	CreateLinks(ctx context.Context, links []domain.Link) error
	Close() error
}

type linkService struct {
	storage   LinkStorage
	count int32
}

func NewLinkService(storage LinkStorage) *linkService {
	return &linkService{
		storage:   storage,
	}
}

func (s *linkService) GetIdent(ctx context.Context, fulLink string) (string, error) {
	ident := s.GenerateIdent()
	link, err := s.storage.Create(ctx, ident, fulLink)
	return link.Ident, err
}

func (s *linkService) GetIdents(ctx context.Context, linkReq []dto.LinkListReq) ([]dto.LinkListRes, error) {
	result := make([]dto.LinkListRes, 0)
	links := make([]domain.Link, 0)
	for _, v := range linkReq {
		ident := s.GenerateIdent()
		result = append(result, dto.LinkListRes{CorrelationID: v.CorrelationID, ShortURL: ident})
		links = append(links, domain.Link{Ident: ident, FulLink: v.OriginalURL})
	}
	if err := s.storage.CreateLinks(ctx, links); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *linkService) GetFulLink(ctx context.Context, ident string) (string, error) {
	link, err := s.storage.GetOneByIdent(ctx, ident)
	if err != nil {
		return "", err
	}
	return link.FulLink, nil
}

func (s *linkService) GenerateIdent() string {
	hd := hashids.NewData()
	hd.Salt = salt
	h, _ := hashids.NewWithData(hd)
	ident, _ := h.Encode([]int{rand.Int() + int(atomic.LoadInt32(&s.count))})
	atomic.AddInt32(&s.count, 1)
	return ident
}

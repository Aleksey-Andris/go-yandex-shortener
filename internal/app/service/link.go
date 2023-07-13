package service

import (
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/speps/go-hashids"
	"sync/atomic"
)

var (
	count = int32(1)
	salt  = "Qw6"
)

type LinkStorage interface {
	GetOneByIdent(ident string) (domain.Link, error)
	Create(shortLink, fulLink string) (domain.Link, error)
}

type linkService struct {
	storage LinkStorage
}

func NewLinkService(storage LinkStorage) *linkService {
	return &linkService{storage: storage}
}

func (s *linkService) GetIdent(fulLink string) (string, error) {
	ident := s.GenerateIdent(fulLink)

	link, err := s.storage.Create(ident, fulLink)
	if err != nil {
		return "", err
	}
	return link.Ident(), nil
}

func (s *linkService) GetFulLink(ident string) (string, error) {
	link, err := s.storage.GetOneByIdent(ident)
	if err != nil {
		return "", err
	}
	return link.FulLink(), nil
}

func (s *linkService) GenerateIdent(fulLink string) string {
	hd := hashids.NewData()
	hd.Salt = salt
	h, _ := hashids.NewWithData(hd)
	ident, _ := h.Encode([]int{int(atomic.LoadInt32(&count))})
	atomic.AddInt32(&count, 1)
	return ident
}

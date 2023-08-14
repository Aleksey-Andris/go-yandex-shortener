package service

import (
	"sync/atomic"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/speps/go-hashids"
)

var (
	salt = "Qw6"
)

type LinkStorage interface {
	GetOneByIdent(ident string) (domain.Link, error)
	Create(idemt, fulLink string) (domain.Link, error)
	CreateLinks(links []domain.Link) error
	Close() error
	GetMaxID() (int32, error)
}

type linkService struct {
	storage   LinkStorage
	nextCount int32
}

func NewLinkService(storage LinkStorage, nextCount int32) *linkService {
	return &linkService{
		storage:   storage,
		nextCount: nextCount,
	}
}

func (s *linkService) GetIdent(fulLink string) (string, error) {
	ident := s.GenerateIdent(fulLink)

	link, err := s.storage.Create(ident, fulLink)
	if err != nil {
		return "", err
	}
	return link.Ident, nil
}

func (s *linkService) GetIdents(linkReq []dto.LinkListReq) ([]dto.LinkListRes, error) {
	result := make([]dto.LinkListRes, 0)
	links := make([]domain.Link, 0)
	for _, v := range linkReq {
		ident := s.GenerateIdent(v.OriginalURL)
		result = append(result, dto.LinkListRes{CorrelationID: v.CorrelationID, ShortURL: ident})
		links = append(links, domain.Link{Ident: ident, FulLink: v.OriginalURL})
	}
	if err := s.storage.CreateLinks(links); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *linkService) GetFulLink(ident string) (string, error) {
	link, err := s.storage.GetOneByIdent(ident)
	if err != nil {
		return "", err
	}
	return link.FulLink, nil
}

func (s *linkService) GenerateIdent(fulLink string) string {
	hd := hashids.NewData()
	hd.Salt = salt
	h, _ := hashids.NewWithData(hd)
	ident, _ := h.Encode([]int{int(atomic.LoadInt32(&s.nextCount))})
	atomic.AddInt32(&s.nextCount, 1)
	return ident
}

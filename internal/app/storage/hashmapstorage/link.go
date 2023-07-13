package hashmapstorage

import (
	"errors"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
)

type linkStorage struct {
	linkMap map[string]domain.Link
}

func NewLinkStorage(linkMap map[string]domain.Link) *linkStorage {
	newLinkMap := linkMap
	return &linkStorage{linkMap: newLinkMap}
}

func (s *linkStorage) GetOneByIdent(ident string) (domain.Link, error) {
	link, ok := s.linkMap[ident]
	if !ok {
		link = domain.Link{}
		return link, errors.New("not found")
	} else {
		return link, nil
	}
}

func (s *linkStorage) Create(ident, fulLink string) (domain.Link, error) {
	link := domain.Link{}
	link.SetIdent(ident)
	link.SetFulLink(fulLink)

	s.linkMap[ident] = link

	return link, nil
}

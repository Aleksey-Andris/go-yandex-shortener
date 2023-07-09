package hash_map_storage

import (
	"errors"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
)

type linkStorage struct {
	linkMap map[string]domain.Link
}

func NewLinkStorage() *linkStorage {
	newLinkMap := make(map[string]domain.Link)
	return &linkStorage{linkMap: newLinkMap}
}

func (s *linkStorage) GetOneByShortLink(shortLink string) (domain.Link, error) {
	link, ok := s.linkMap[shortLink]
	if !ok {
		link = domain.Link{}
		return link, errors.New("not found")
	} else {
		return link, nil
	}
}

func (s *linkStorage) Create(shortLink, fulLink string) (domain.Link, error) {
	link := domain.Link{}
	link.SetShortLink(shortLink)
	link.SetFulLink(fulLink)

	s.linkMap[shortLink] = link

	return link, nil
}

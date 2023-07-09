package service

import (
	"crypto/md5"
	"fmt"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
)

type LinkStorage interface {
	GetOneByShortLink(shortLink string) (domain.Link, error)
	Create(shortLink, fulLink string) (domain.Link, error)
}

type linkService struct {
	storage LinkStorage
}

func NewLinkService(storage LinkStorage) *linkService {
	return &linkService{storage: storage}
}

func (s *linkService) GetShortLink(fulLink string) (string, error) {
	shortLink := s.GenerateShortLink(fulLink)

	link, err := s.storage.GetOneByShortLink(shortLink)
	if err != nil {
		link, err = s.storage.Create(shortLink, fulLink)
	}

	return link.ShortLink(), nil
}

func (s *linkService) GetFulLink(shortLink string) (string, error) {
	link, err := s.storage.GetOneByShortLink(shortLink)
	if err != nil {
		return "", err
	}
	return link.FulLink(), nil
}

func (s *linkService) GenerateShortLink(fulLink string) string {
	hash := md5.New()
	hash.Write([]byte(fulLink))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

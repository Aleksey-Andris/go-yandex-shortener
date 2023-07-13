package service

import (
	"crypto/md5"
	"fmt"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
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

	link, err := s.storage.GetOneByIdent(ident)
	if err != nil {
		link, err = s.storage.Create(ident, fulLink)
		if err != nil {
			return "", err
		}
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
	hash := md5.New()
	hash.Write([]byte(fulLink))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

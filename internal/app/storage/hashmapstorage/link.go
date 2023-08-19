package hashmapstorage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
)

type linkStorage struct {
	sync.Mutex
	linkMap    map[string]domain.Link
	record     bool
	filePath   string
	file       *os.File
	decoder    *json.Decoder
	encoder    *json.Encoder
}

func NewLinkStorage(linkMap map[string]domain.Link, filePath string) (*linkStorage, error) {

	storage := &linkStorage{
		linkMap:    linkMap,
		record:     filePath != "",
		filePath:   filePath,
	}
	if filePath != "" {
		if err := storage.loadFromFile(); err != nil {
			return &linkStorage{}, err
		}
	}
	return storage, nil
}

func (s *linkStorage) GetOneByIdent(ctx context.Context, ident string) (domain.Link, error) {
	s.Lock()
	defer s.Unlock()
	link, ok := s.linkMap[ident]
	if !ok {
		link = domain.Link{}
		return link, errors.New("not found")
	}
	return link, nil
}

func (s *linkStorage) Create(ctx context.Context, ident, fulLink string) (domain.Link, error) {
	s.Lock()
	defer s.Unlock()
	link := domain.Link{
		Ident:   ident,
		FulLink: fulLink,
	}
	if s.record {
		if err := s.encoder.Encode(&link); err != nil {
			return domain.Link{}, err
		}
	}
	s.linkMap[ident] = link
	return link, nil
}

func (s *linkStorage) CreateLinks(ctx context.Context, links []domain.Link) error {
	s.Lock()
	defer s.Unlock()
	if s.record {
		for _, v := range links {
			if err := s.encoder.Encode(&v); err != nil {
				return err
			}
		}
	}
	for _, v := range links {
		s.linkMap[v.Ident] = v
	}
	return nil
}

func (s *linkStorage) loadFromFile() error {
	var err error
	s.file, err = os.OpenFile(s.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	s.decoder = json.NewDecoder(s.file)
	s.encoder = json.NewEncoder(s.file)

	var link domain.Link
	for {
		err = s.decoder.Decode(&link)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		s.linkMap[link.Ident] = link
	}
	return nil
}

func (s *linkStorage) Close() error {
	if s.file == nil {
		return nil
	}
	return s.file.Close()
}

package hashmapstorage

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
)

type linkStorage struct {
	sync.Mutex
	linkMap   map[string]domain.Link
	record    bool
	filePath  string
	file      *os.File
	decoder   *json.Decoder
	encoder   *json.Encoder
	seqUserID int32
}

func NewLinkStorage(linkMap map[string]domain.Link, filePath string) (*linkStorage, error) {

	storage := &linkStorage{
		linkMap:   linkMap,
		record:    filePath != "",
		filePath:  filePath,
		seqUserID: 1,
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

func (s *linkStorage) Create(ctx context.Context, ident, fulLink string, userID int32) (domain.Link, error) {
	s.Lock()
	defer s.Unlock()
	if s.seqUserID < userID {
		s.seqUserID = userID
	}
	link := domain.Link{
		Ident:   ident,
		FulLink: fulLink,
		UserID:  userID,
	}
	if s.record {
		if err := s.encoder.Encode(&link); err != nil {
			return domain.Link{}, err
		}
	}
	s.linkMap[ident] = link
	return link, nil
}

func (s *linkStorage) CreateLinks(ctx context.Context, links []domain.Link, userID int32) error {
	s.Lock()
	defer s.Unlock()
	if s.seqUserID < userID {
		s.seqUserID = userID
	}
	if s.record {
		for _, v := range links {
			if err := s.encoder.Encode(&v); err != nil {
				return err
			}
		}
	}
	for _, v := range links {
		v.UserID = userID
		s.linkMap[v.Ident] = v
	}
	return nil
}

func (s *linkStorage) GetLinksByUserId(ctx context.Context, userID int32) ([]dto.LinkListByUserIdRes, error) {
	s.Lock()
	defer s.Unlock()
	var linkListByUserIdRes []dto.LinkListByUserIdRes
	for _, v := range s.linkMap {
		if v.UserID == userID {
			linkListByUserIdRes = append(linkListByUserIdRes, dto.LinkListByUserIdRes{
				OriginalURL: v.FulLink,
				ShortURL:    v.Ident,
			})
		}
	}
	return linkListByUserIdRes, nil
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
		if s.seqUserID < link.UserID {
			s.seqUserID = link.UserID
		}
		s.linkMap[link.Ident] = link
	}
	return nil
}

func (s *linkStorage) CreateUser(ctx context.Context) (int32, error) {
	s.Lock()
	defer s.Unlock()
	s.seqUserID++
	return s.seqUserID, nil
}

func (s *linkStorage) Close() error {
	if s.file == nil {
		return nil
	}
	return s.file.Close()
}

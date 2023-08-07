package hashmapstorage

import (
	"encoding/json"
	"errors"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"os"
	"sync"
)

type linkStorage struct {
	sync.Mutex
	linkMap    map[string]domain.Link
	sequenceID int32
	record     bool
	filePath   string
	file       *os.File
	decoder    *json.Decoder
	encoder    *json.Encoder
}

func NewLinkStorage(linkMap map[string]domain.Link, filePath string) (*linkStorage, error) {

	storage := &linkStorage{
		linkMap:    linkMap,
		sequenceID: 1,
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

func (s *linkStorage) GetOneByIdent(ident string) (domain.Link, error) {
	s.Lock()
	defer s.Unlock()
	link, ok := s.linkMap[ident]
	if !ok {
		link = domain.Link{}
		return link, errors.New("not found")
	}
	return link, nil
}

func (s *linkStorage) Create(ident, fulLink string) (domain.Link, error) {
	s.Lock()
	defer s.Unlock()
	link := domain.Link{
		ID:      s.sequenceID,
		Ident:   ident,
		FulLink: fulLink,
	}
	if s.record {
		if err := s.encoder.Encode(&link); err != nil {
			return domain.Link{}, err
		}
	}
	s.linkMap[ident] = link
	s.sequenceID++
	return link, nil
}

func (s *linkStorage) loadFromFile() error {
	s.Lock()
	defer s.Unlock()
	data, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	s.decoder = json.NewDecoder(data)
	s.encoder = json.NewEncoder(data)

	var link domain.Link
	for {
		err = s.decoder.Decode(&link)
		if err != nil {
			break
		}
		s.linkMap[link.Ident] = link
		if link.ID > s.sequenceID {
			s.sequenceID = link.ID
		}
	}
	return nil
}

func (s *linkStorage) Close() error {
	return s.file.Close()
}

func (s *linkStorage) GetSequense() int32 {
	s.Lock()
	defer s.Unlock()
	return s.sequenceID
}

package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/mock/mockservice"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_linkHandler_GetShortLink(t *testing.T) {
	tests := []struct {
		name               string
		requestURL         string
		requestContentType string
		requestBody        []byte
		expectedStatusCode int
	}{
		{
			name:               "simple case",
			requestURL:         "/",
			requestContentType: "text/plain",
			requestBody:        []byte("https://practicum.test1.ru/"),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "two Content Type",
			requestURL:         "/",
			requestContentType: "text/plain; utf8",
			requestBody:        []byte("https://practicum.test2.ru/"),
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "incorrect Content Type",
			requestURL:         "/",
			requestContentType: "text/pain",
			requestBody:        []byte("https://practicum.test3.ru/"),
			expectedStatusCode: http.StatusBadRequest,
		},

		{
			name:               "empty Content Type",
			requestURL:         "/",
			requestContentType: "",
			requestBody:        []byte("https://practicum.test4.ru/"),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	linkStorage, _ := hashmapstorage.NewLinkStorage(make(map[string]domain.Link), "")
	linkService := service.NewLinkService(linkStorage, 1)
	linkHandler := NewLinkHandler(linkService, "http://localhost:8080")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.requestURL, bytes.NewReader(tt.requestBody))
			request.Header.Set("Content-Type", tt.requestContentType)

			rec := httptest.NewRecorder()
			linkHandler.GetShortLink(rec, request)

			res := rec.Result()

			assert.Equal(t, res.StatusCode, tt.expectedStatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

func Test_linkHandler_GetFulLink(t *testing.T) {
	tests := []struct {
		name               string
		requestURL         string
		paramURL           string
		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name:               "simple case",
			requestURL:         "/",
			paramURL:           "123456",
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedLocation:   "https://practicum.test5.ru/",
		},

		{
			name:               "not found case",
			requestURL:         "/",
			paramURL:           "123457",
			expectedStatusCode: http.StatusBadRequest,
			expectedLocation:   "",
		},
	}
	linkMap := make(map[string]domain.Link)
	link := domain.Link{
		ID: 1,
		Ident: "123456",
		FulLink: "https://practicum.test5.ru/",
	}
	linkMap[link.Ident] = link
	linkStorage, _ := hashmapstorage.NewLinkStorage(linkMap, "")
	linkService := service.NewLinkService(linkStorage, 1)
	linkHandler := NewLinkHandler(linkService, "http://localhost:8080")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ident", tt.paramURL)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			linkHandler.GetFulLink(rec, request)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedStatusCode)
			assert.Equal(t, res.Header.Get("Location"), tt.expectedLocation)
		})
	}
}

func Test_linkHandler_GetShortLinkByJson(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	linkStorage := mockservice.NewMockLinkStorage(c)
	linkService := service.NewLinkService(linkStorage, 1)
	linkHandler := NewLinkHandler(linkService, "http://localhost:8080")
	testServ := httptest.NewServer(linkHandler.InitRouter())
	defer testServ.Close()

	type mocBehavior func(s *mockservice.MockLinkStorage)
	tests := []struct {
		name               string
		requestURL         string
		requestBody        string
		requestContentType string
		expectedStatusCode int
		expectedErr        bool
		mocBehavior        mocBehavior
	}{
		{
			name:               "simple case (json)",
			requestURL:         "/api/shorten",
			requestBody:        `{"url": "https://practicum.test.ru/"}`,
			requestContentType: "application/json",
			expectedStatusCode: http.StatusCreated,
			expectedErr:        false,
			mocBehavior: func(s *mockservice.MockLinkStorage) {
				link := domain.Link{
					ID: 1,
					Ident: "some_ident",
					FulLink: "some_link",
				}
				s.EXPECT().Create(gomock.Any(), gomock.Any()).Return(link, nil)
			},
		},

		{
			name:               "incorrect Content Type (json)",
			requestURL:         "/api/shorten",
			requestContentType: "application/gson",
			requestBody:        `{"url": "https://practicum.test.ru/"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			mocBehavior:        func(s *mockservice.MockLinkStorage) {},
		},

		{
			name:               "empty Content Type (json)",
			requestURL:         "/api/shorten",
			requestContentType: "",
			requestBody:        `{"url": "https://practicum.test.ru/"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			mocBehavior:        func(s *mockservice.MockLinkStorage) {},
		},

		{
			name:               "empty body",
			requestURL:         "/api/shorten",
			requestContentType: "application/json",
			requestBody:        "",
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			mocBehavior:        func(s *mockservice.MockLinkStorage) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocBehavior(linkStorage)

			req, err := http.NewRequest(http.MethodPost, testServ.URL+tt.requestURL, bytes.NewBufferString(tt.requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", tt.requestContentType)

			res, err := testServ.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			if !tt.expectedErr {
				err = json.NewDecoder(res.Body).Decode(&dto.LinkRes{})
				require.NoError(t, err)
			}
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s with gzip", tt.name), func(t *testing.T) {
			tt.mocBehavior(linkStorage)

			buf := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buf)
			_, err := zb.Write([]byte(tt.requestBody))
			require.NoError(t, err)
			err = zb.Close()
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, testServ.URL+tt.requestURL, buf)
			require.NoError(t, err)
			req.Header.Set("Content-Type", tt.requestContentType)
			req.Header.Set("Content-Encoding", "gzip")

			res, err := testServ.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			if !tt.expectedErr {
				err = json.NewDecoder(res.Body).Decode(&dto.LinkRes{})
				require.NoError(t, err)
			}
		})
	}
}

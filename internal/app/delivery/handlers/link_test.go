package handlers

import (
	"bytes"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

	linkStorage := hashmapstorage.NewLinkStorage(make(map[string]domain.Link))
	linkService := service.NewLinkService(linkStorage)
	linkHandler := NewLinkHandler(linkService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.requestURL, bytes.NewReader(tt.requestBody))
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
		expectedStatusCode int
		expectedLocation   string
	}{
		{
			name:               "simple case",
			requestURL:         "/123456",
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedLocation:   "https://practicum.test5.ru/",
		},

		{
			name:               "not found case",
			requestURL:         "/1234567",
			expectedStatusCode: http.StatusBadRequest,
			expectedLocation:   "",
		},
	}
	linkMap := make(map[string]domain.Link)
	link := domain.Link{}
	link.SetIdent("123456")
	link.SetFulLink("https://practicum.test5.ru/")
	linkMap["123456"] = link
	linkStorage := hashmapstorage.NewLinkStorage(linkMap)
	linkService := service.NewLinkService(linkStorage)
	linkHandler := NewLinkHandler(linkService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)

			rec := httptest.NewRecorder()
			linkHandler.GetFulLink(rec, request)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.expectedStatusCode)
			assert.Equal(t, res.Header.Get("Location"), tt.expectedLocation)
		})
	}
}

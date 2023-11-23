package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/mock/mockservice"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
)

// Test_Handler_GetShortLink - tests for GetShortLink handler
func Test_Handler_GetShortLink(t *testing.T) {
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

	linkStorage, _ := hashmapstorage.NewLinkStorage(make(map[string]*domain.Link), make(map[int32][]*domain.Link), "")
	userStorage, _ := hashmapstorage.NewLinkStorage(make(map[string]*domain.Link), make(map[int32][]*domain.Link), "")
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.requestURL, bytes.NewReader(tt.requestBody))
			request.Header.Set("Content-Type", tt.requestContentType)

			rec := httptest.NewRecorder()
			handler.GetShortLink(rec, request)

			res := rec.Result()

			assert.Equal(t, res.StatusCode, tt.expectedStatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
		})
	}
}

// Test_Handler_GetFulLink - tests for GetFulLink handler
func Test_Handler_GetFulLink(t *testing.T) {
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

		{
			name:               "is deleted",
			requestURL:         "/",
			paramURL:           "123458",
			expectedStatusCode: http.StatusGone,
			expectedLocation:   "",
		},
	}
	linkMap := make(map[string]*domain.Link)
	link := domain.Link{
		ID:      1,
		Ident:   "123456",
		FulLink: "https://practicum.test5.ru/",
	}
	linkDeleted := domain.Link{
		ID:          2,
		Ident:       "123458",
		FulLink:     "https://practicum.test8.ru/",
		DeletedFlag: true,
	}
	linkMap[link.Ident] = &link
	linkMap[linkDeleted.Ident] = &linkDeleted
	linkStorage, _ := hashmapstorage.NewLinkStorage(linkMap, make(map[int32][]*domain.Link), "")
	userStorage, _ := hashmapstorage.NewLinkStorage(linkMap, make(map[int32][]*domain.Link), "")
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.requestURL, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("ident", tt.paramURL)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			handler.GetFulLink(rec, request)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Equal(t, tt.expectedLocation, res.Header.Get("Location"))
		})
	}
}

// Test_Handler_GetShortLinkByJson - tests for GetShortLinkByJson handler
func Test_Handler_GetShortLinkByJson(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	userStorage := mockservice.NewMockUserStorage(c)
	linkStorage := mockservice.NewMockLinkStorage(c)
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")
	testServ := httptest.NewServer(handler.InitRouter())
	defer testServ.Close()

	type mocBehavior func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage)
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
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				link := domain.Link{
					ID:      1,
					Ident:   "some_ident",
					FulLink: "some_link",
				}
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
				sl.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(link, nil)
			},
		},

		{
			name:               "incorrect Content Type (json)",
			requestURL:         "/api/shorten",
			requestContentType: "application/gson",
			requestBody:        `{"url": "https://practicum.test.ru/"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
			},
		},

		{
			name:               "empty Content Type (json)",
			requestURL:         "/api/shorten",
			requestContentType: "",
			requestBody:        `{"url": "https://practicum.test.ru/"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
			},
		},

		{
			name:               "empty body",
			requestURL:         "/api/shorten",
			requestContentType: "application/json",
			requestBody:        "",
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocBehavior(userStorage, linkStorage)

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
			tt.mocBehavior(userStorage, linkStorage)

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

// Test_Handler_GetShortLinkByListJSON - tests for GetShortLinkByListJSON handler
func Test_Handler_GetShortLinkByListJSON(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	userStorage := mockservice.NewMockUserStorage(c)
	linkStorage := mockservice.NewMockLinkStorage(c)
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")
	testServ := httptest.NewServer(handler.InitRouter())
	defer testServ.Close()

	type mocBehavior func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage)
	tests := []struct {
		name               string
		requestURL         string
		requestBody        string
		requestContentType string
		expectedStatusCode int
		expectedErr        bool
		expectedListSize   int
		mocBehavior        mocBehavior
	}{
		{
			name:               "simple case (json)",
			requestURL:         "/api/shorten/batch",
			requestBody:        `[{"correlation_id": "string_ident1","original_url":"https://practicum.test1.ru/"},{"correlation_id":"string_ident2","original_url":"https://practicum.test2.ru/"}]`,
			requestContentType: "application/json",
			expectedStatusCode: http.StatusCreated,
			expectedErr:        false,
			expectedListSize:   2,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
				sl.EXPECT().CreateLinks(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},

		{
			name:               "incorrect Content Type (json)",
			requestURL:         "/api/shorten/batch",
			requestContentType: "application/gson",
			requestBody:        `[{"correlation_id": "string_ident1","original_url":"https://practicum.test1.ru/"},{"correlation_id":"string_ident2","original_url":"https://practicum.test2.ru/"}]`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
			},
		},

		{
			name:               "empty body",
			requestURL:         "/api/shorten/batch",
			requestContentType: "application/json",
			requestBody:        "",
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocBehavior(userStorage, linkStorage)

			req, err := http.NewRequest(http.MethodPost, testServ.URL+tt.requestURL, bytes.NewBufferString(tt.requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", tt.requestContentType)

			res, err := testServ.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			if !tt.expectedErr {
				var buf bytes.Buffer
				var linkRes []dto.LinkListRes

				_, err := buf.ReadFrom(res.Body)
				require.NoError(t, err)

				err = json.Unmarshal(buf.Bytes(), &linkRes)
				require.NoError(t, err)
			}
		})
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s with gzip", tt.name), func(t *testing.T) {
			tt.mocBehavior(userStorage, linkStorage)

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
				var buf bytes.Buffer
				var linkRes []dto.LinkListRes

				_, err := buf.ReadFrom(res.Body)
				require.NoError(t, err)

				err = json.Unmarshal(buf.Bytes(), &linkRes)
				require.NoError(t, err)

				assert.True(t, tt.expectedListSize == len(linkRes))
			}
		})
	}
}

// Test_Handler_GetLinksByUser - tests for GetLinksByUser handler
func Test_Handler_GetLinksByUser(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	userStorage := mockservice.NewMockUserStorage(c)
	linkStorage := mockservice.NewMockLinkStorage(c)
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")
	testServ := httptest.NewServer(handler.InitRouter())
	defer testServ.Close()

	type mocBehavior func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage)
	tests := []struct {
		name               string
		requestURL         string
		expectedStatusCode int
		IsNoContent        bool
		expectedListSize   int
		mocBehavior        mocBehavior
	}{
		{
			name:               "geting link by user - simple case",
			requestURL:         "/api/user/urls",
			expectedStatusCode: http.StatusOK,
			IsNoContent:        false,
			expectedListSize:   2,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				link1 := dto.LinkListByUserIDRes{
					OriginalURL: "some_oriq_url",
					ShortURL:    "some_orig_url",
				}
				link2 := dto.LinkListByUserIDRes{
					OriginalURL: "some_oriq_url_2",
					ShortURL:    "some_orig_url_3",
				}
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
				sl.EXPECT().GetLinksByUserID(gomock.Any(), gomock.Any()).Return([]dto.LinkListByUserIDRes{link1, link2}, nil)
			},
		},
		{
			name:               "geting link by user - no content",
			requestURL:         "/api/user/urls",
			expectedStatusCode: http.StatusNoContent,
			IsNoContent:        true,
			expectedListSize:   2,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(1), nil)
				sl.EXPECT().GetLinksByUserID(gomock.Any(), gomock.Any()).Return([]dto.LinkListByUserIDRes{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocBehavior(userStorage, linkStorage)

			req, err := http.NewRequest(http.MethodGet, testServ.URL+tt.requestURL, nil)
			require.NoError(t, err)

			res, err := testServ.Client().Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			if !tt.IsNoContent {
				var buf bytes.Buffer
				var linkRes []dto.LinkListByUserIDRes

				_, err := buf.ReadFrom(res.Body)
				require.NoError(t, err)

				err = json.Unmarshal(buf.Bytes(), &linkRes)
				require.NoError(t, err)
			}
		})
	}
}

// Test_Handler_DeleteLinksByIdents - tests for DeleteLinksByIdents handler
func Test_Handler_DeleteLinksByIdents(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	userStorage := mockservice.NewMockUserStorage(c)
	linkStorage := mockservice.NewMockLinkStorage(c)
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")
	testServ := httptest.NewServer(handler.InitRouter())
	defer testServ.Close()
	type mocBehavior func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage)
	tests := []struct {
		name               string
		requestURL         string
		requestBody        string
		requestContentType string
		expectedStatusCode int
		mocBehavior        mocBehavior
	}{
		{
			name:               "delete link by idents - forbidden",
			requestURL:         "/api/user/urls",
			requestBody:        `["some_ident1", "some_ident2", "some_ident3" ]`,
			requestContentType: "application/json",
			expectedStatusCode: http.StatusForbidden,
			mocBehavior: func(sa *mockservice.MockUserStorage, sl *mockservice.MockLinkStorage) {
				link1 := domain.Link{
					ID:      1,
					Ident:   "some_ident1",
					FulLink: "some_link1",
					UserID:  1,
				}
				link2 := domain.Link{
					ID:      2,
					Ident:   "some_ident2",
					FulLink: "some_link2",
					UserID:  1,
				}
				link3 := domain.Link{
					ID:      3,
					Ident:   "some_ident3",
					FulLink: "some_link3",
					UserID:  1,
				}
				links := []domain.Link{link1, link2, link3}
				sa.EXPECT().CreateUser(gomock.Any()).Return(int32(2), nil)
				sl.EXPECT().GetByIdents(gomock.Any(), "some_ident1", "some_ident2", "some_ident3").Return(links, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocBehavior(userStorage, linkStorage)

			req, err := http.NewRequest(http.MethodDelete, testServ.URL+tt.requestURL, bytes.NewBufferString(tt.requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", tt.requestContentType)

			res, err := testServ.Client().Do(req)

			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
		})
	}
}

/*
// BenchmarkGetShortLink - benchmark for GetShortLink handler
func BenchmarkGetShortLink(b *testing.B) {
	random := rand.NewSource(time.Now().UnixNano())
	requestURL := "/"
	requestContentType := "text/plain"
	requestBody := "https://practicum.test1.ru/"
	linkStorage, _ := hashmapstorage.NewLinkStorage(make(map[string]*domain.Link), make(map[int32][]*domain.Link), "")
	userStorage, _ := hashmapstorage.NewLinkStorage(make(map[string]*domain.Link), make(map[int32][]*domain.Link), "")
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		newRequestBody := requestBody + strconv.Itoa(int(random.Int63()))
		request := httptest.NewRequest(http.MethodPost, requestURL, bytes.NewBufferString(newRequestBody))
		request.Header.Set("Content-Type", requestContentType)
		rec := httptest.NewRecorder()
		b.StartTimer()
		handler.GetShortLink(rec, request)
	}
}

// BenchmarkGetFulLink - benchmark for GetFulLink handler
func BenchmarkGetFulLink(b *testing.B) {
	requestURL := "/"
	paramURL := "123456"
	linkMap := make(map[string]*domain.Link)
	link := domain.Link{
		ID:      1,
		Ident:   "123456",
		FulLink: "https://practicum.test5.ru/",
	}
	linkMap[link.Ident] = &link
	linkStorage, _ := hashmapstorage.NewLinkStorage(linkMap, make(map[int32][]*domain.Link), "")
	userStorage, _ := hashmapstorage.NewLinkStorage(linkMap, make(map[int32][]*domain.Link), "")
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		request := httptest.NewRequest(http.MethodGet, requestURL, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("ident", paramURL)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))
		rec := httptest.NewRecorder()
		b.StartTimer()
		handler.GetFulLink(rec, request)
	}
} */

// BenchmarkGetLinksByUser - benchmark for GetLinksByUser  handler
func BenchmarkGetLinksByUser(b *testing.B) {
	requestURL := "/api/user/urls"
	linkMap := make(map[string]*domain.Link)
	linkByUserIDMap := make(map[int32][]*domain.Link)
	linkStorage, _ := hashmapstorage.NewLinkStorage(linkMap, linkByUserIDMap, "")
	userStorage, _ := hashmapstorage.NewLinkStorage(linkMap, linkByUserIDMap, "")
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")
	testServ := httptest.NewServer(handler.InitRouter())
	defer testServ.Close()

	userID := int32(1)
	ident := "123456"
	fulLink := "https://practicum.test5.ru/"
	linkSlice := make([]*domain.Link, 0)
	for i := 0; i < 20; i++ {
		ident := ident + strconv.Itoa(i)
		fulLink := fulLink + strconv.Itoa(i)
		link := domain.Link{
			UserID:  userID,
			Ident:   ident,
			FulLink: fulLink,
		}
		linkSlice = append(linkSlice, &link)
		linkMap[link.Ident] = &link
	}
	linkByUserIDMap[userID] = linkSlice
	for i := 20; i < 10000000; i++ {
		ident := ident + strconv.Itoa(i)
		fulLink := fulLink + strconv.Itoa(i)
		link := domain.Link{
			UserID:  int32(i),
			Ident:   ident,
			FulLink: fulLink,
		}
		linkByUserIDMap[int32(i)] = []*domain.Link{&link}
		linkMap[link.Ident] = &link
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		request, _ := http.NewRequest(http.MethodGet, testServ.URL+requestURL, nil)
		token, _ := handler.services.AuthService.BuildJWTString(int32(userID))
		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		})
		b.StartTimer()
		res, _ := testServ.Client().Do(request)
		b.StopTimer()
		res.Body.Close()
	}
}

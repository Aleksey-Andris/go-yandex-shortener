package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
)

func Example() {
	//Preparing the server infrastructure.
	linkMap := make(map[string]domain.Link)
	linkStorage, _ := hashmapstorage.NewLinkStorage(linkMap, "")
	userStorage, _ := hashmapstorage.NewLinkStorage(linkMap, "")
	servises := NewServices(linkStorage, userStorage)
	handler := NewHandler(servises, "http://localhost:8080")
	testServ := httptest.NewServer(handler.InitRouter())
	errRedirectBlocked := errors.New("HTTP redirect blocked")
	testServ.Client().CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	}

	defer testServ.Close()

	//Preparing data.
	userID := 1
	ident := "123456"
	fulLink := "https://practicum.test5.ru/"
	link := domain.Link{
		UserID:  int32(userID),
		Ident:   ident,
		FulLink: fulLink,
	}
	linkMap[link.Ident] = link
	token, _ := handler.services.AuthService.BuildJWTString(int32(userID))

	//Example request for "POST: .../"
	request, _ := http.NewRequest(http.MethodPost, testServ.URL+"/", bytes.NewBufferString("https://practicum.test1.ru/"))
	request.Header.Set("Content-Type", "text/plain")
	res, _ := testServ.Client().Do(request)
	fmt.Println(res.Status)
	res.Body.Close()

	//Example request for "GET: ...//{ident}".
	request, _ = http.NewRequest(http.MethodGet, testServ.URL+"/123456", nil)
	request.Header.Set("Content-Type", "text/plain")
	request.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
	res, _ = testServ.Client().Do(request)
	fmt.Println(res.Status)
	res.Body.Close()

	// Example request for "POST: .../api/shorten".
	request, _ = http.NewRequest(http.MethodPost, testServ.URL+"/api/shorten", bytes.NewBufferString(`{"url": "https://practicum.test.ru/"}`))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
	res, _ = testServ.Client().Do(request)
	fmt.Println(res.Status)
	res.Body.Close()

	// Example request for "POST: .../api/shorten/batch".
	request, _ = http.NewRequest(http.MethodPost, testServ.URL+"/api/shorten/batch",
		bytes.NewBufferString(`[{"correlation_id": "string_ident1","original_url":"https://practicum.test1.ru/"},`+
			`{"correlation_id":"string_ident2","original_url":"https://practicum.test2.ru/"}]`))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
	res, _ = testServ.Client().Do(request)
	fmt.Println(res.Status)
	res.Body.Close()

	// Example request for "GET: .../api/user/urls".
	request, _ = http.NewRequest(http.MethodGet, testServ.URL+"/api/user/urls", nil)
	token, _ = handler.services.AuthService.BuildJWTString(int32(userID))
	request.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
	res, _ = testServ.Client().Do(request)
	fmt.Println(res.Status)
	res.Body.Close()

	// Example request for "DELETE: ...///api/user/urls".
	request, _ = http.NewRequest(http.MethodDelete, testServ.URL+"/api/user/urls", bytes.NewBufferString(`["123456"]`))
	token, _ = handler.services.AuthService.BuildJWTString(int32(userID))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
	res, _ = testServ.Client().Do(request)
	fmt.Println(res.Status)
	res.Body.Close()

	// Output:
	// 201 Created
	// 307 Temporary Redirect
	// 201 Created
	// 201 Created
	// 200 OK
	// 202 Accepted
}

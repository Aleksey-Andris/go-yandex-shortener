package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	ContentType          = "Content-Type"
	ContentTypeTextPlain = "text/plain"
	ServerURL            = "http://localhost:8080/"
)

type LinkService interface {
	GetFulLink(ident string) (string, error)
	GetIdent(fulLink string) (string, error)
	GenerateIdent(fulLink string) string
}

type linkHandler struct {
	service LinkService
}

func NewLinkHandler(service LinkService) *linkHandler {
	return &linkHandler{service: service}
}

func (h *linkHandler) InitServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.GetLink)

	return mux
}

func (h *linkHandler) GetLink(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		h.GetShortLink(res, req)
	} else if req.Method == http.MethodGet {
		h.GetFulLink(res, req)
	} else {
		http.Error(res, "invalid method", http.StatusBadRequest)
	}

}

func (h *linkHandler) GetShortLink(res http.ResponseWriter, req *http.Request) {
	if strings.Split(req.Header.Get(ContentType), ";")[0] != ContentTypeTextPlain {
		http.Error(res, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "failed reading body", http.StatusBadRequest)
		return
	}

	ident, err := h.service.GetIdent(string(body))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	shortLink := ServerURL + ident

	res.Header().Set(ContentType, ContentTypeTextPlain)
	res.Header().Set("Content-Length", strconv.Itoa(len(shortLink)))
	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(shortLink))
}

func (h *linkHandler) GetFulLink(res http.ResponseWriter, req *http.Request) {
	fulLink, err := h.service.GetFulLink(strings.Replace(req.URL.Path, "/", "", 1))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", fulLink)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

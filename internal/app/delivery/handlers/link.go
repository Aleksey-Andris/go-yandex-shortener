package handlers

import (
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	ContentType          = "Content-Type"
	ContentTypeTextPlain = "text/plain"
)

type LinkService interface {
	GetFulLink(ident string) (string, error)
	GetIdent(fulLink string) (string, error)
	GenerateIdent(fulLink string) string
}

type linkHandler struct {
	service      LinkService
	baseShortURL string
}

func NewLinkHandler(service LinkService, baseShortURL string) *linkHandler {
	return &linkHandler{
		service:      service,
		baseShortURL: baseShortURL}
}

func (h *linkHandler) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", h.GetShortLink)
	router.Get("/{ident}", h.GetFulLink)
	return router
}

func (h *linkHandler) GetShortLink(res http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get(ContentType)
	if strings.Split(contentType, ";")[0] != ContentTypeTextPlain {
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

	shortLink := h.baseShortURL + "/" + ident

	res.Header().Set(ContentType, ContentTypeTextPlain)
	res.Header().Set("Content-Length", strconv.Itoa(len(shortLink)))
	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(shortLink))
}

func (h *linkHandler) GetFulLink(res http.ResponseWriter, req *http.Request) {
	fulLink, err := h.service.GetFulLink(chi.URLParam(req, "ident"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", fulLink)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

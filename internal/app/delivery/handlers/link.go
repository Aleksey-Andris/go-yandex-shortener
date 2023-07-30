package handlers

import (
	"encoding/json"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/gzipmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/logmiddleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io"
	"net/http"
	"strings"
)

const (
	сontentType          = "Content-Type"
	сontentTypeTextPlain = "text/plain"
	сontentTypeAppJSON   = "application/json"
	сontentTypeAppXGZIP  = "application/x-gzip"
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
	router.Use(logmiddleware.WithLogging)
	router.Use(gzipmiddleware.Decompress)
	router.Use(middleware.Compress(5, "application/json", "text/html"))
	router.Post("/", h.GetShortLink)
	router.Post("/api/shorten", h.GetShortLinkByJSON)
	router.Get("/{ident}", h.GetFulLink)
	return router
}

func (h *linkHandler) GetShortLink(res http.ResponseWriter, req *http.Request) {
	ct := strings.Split(req.Header.Get(сontentType), ";")[0]
	if !(ct == сontentTypeTextPlain || ct == сontentTypeAppXGZIP) {
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

	res.Header().Set(сontentType, сontentTypeTextPlain)
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

func (h *linkHandler) GetShortLinkByJSON(res http.ResponseWriter, req *http.Request) {
	ct := req.Header.Get(сontentType)
	if !(ct == сontentTypeAppJSON || ct == сontentTypeAppXGZIP) {
		http.Error(res, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	var request dto.LinkReq

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	ident, err := h.service.GetIdent(request.URL)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	shortLink := h.baseShortURL + "/" + ident

	res.Header().Set(сontentType, сontentTypeAppJSON)
	res.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(res).Encode(&dto.LinkRes{Result: shortLink}); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

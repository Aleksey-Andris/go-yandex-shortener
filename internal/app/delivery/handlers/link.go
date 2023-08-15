package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/gzipmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/logmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/postgresstorage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	GetIdents(linkReq []dto.LinkListReq) ([]dto.LinkListRes, error)
	GenerateIdent(fulLink string) string
}

type linkHandler struct {
	service      LinkService
	baseShortURL string
}

func NewLinkHandler(service LinkService, baseShortURL string) *linkHandler {
	return &linkHandler{
		service:      service,
		baseShortURL: baseShortURL,
	}
}

func (h *linkHandler) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(logmiddleware.WithLogging)
	router.Use(gzipmiddleware.Decompress)
	router.Use(middleware.Compress(5, "application/json", "text/html"))
	router.Post("/", h.GetShortLink)
	router.Post("/api/shorten", h.GetShortLinkByJSON)
	router.Post("/api/shorten/batch", h.GetShortLinkByListJSON)
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

	var status int
	ident, err := h.service.GetIdent(string(body))
	if err != nil {
		if !errors.Is(err, postgresstorage.ErrConflict) {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		status = http.StatusConflict
	} else {
		status = http.StatusCreated
	}

	shortLink := h.baseShortURL + "/" + ident
	res.Header().Set(сontentType, сontentTypeTextPlain)
	res.WriteHeader(status)
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

	var status int
	ident, err := h.service.GetIdent(request.URL)
	if err != nil {
		if !errors.Is(err, postgresstorage.ErrConflict) {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		status = http.StatusConflict
	} else {
		status = http.StatusCreated
	}
	shortLink := h.baseShortURL + "/" + ident

	res.Header().Set(сontentType, сontentTypeAppJSON)
	res.WriteHeader(status)
	if err := json.NewEncoder(res).Encode(&dto.LinkRes{Result: shortLink}); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *linkHandler) GetShortLinkByListJSON(res http.ResponseWriter, req *http.Request) {
	ct := req.Header.Get(сontentType)
	if !(ct == сontentTypeAppJSON || ct == сontentTypeAppXGZIP) {
		http.Error(res, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	var linkReq []dto.LinkListReq
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, "invalid body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &linkReq); err != nil {
		http.Error(res, "invalid format body", http.StatusBadRequest)
		return
	}

	limkResp, err := h.service.GetIdents(linkReq)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, v := range limkResp {
		limkResp[i].ShortURL = h.baseShortURL + "/" + v.ShortURL
	}

	response, err := json.Marshal(&limkResp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set(сontentType, сontentTypeAppJSON)
	res.WriteHeader(http.StatusCreated)
	res.Write(response)
}

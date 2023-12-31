package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/logger"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/postgresstorage"
	"github.com/go-chi/chi"
)

const (
	сontentType          = "Content-Type"
	сontentTypeTextPlain = "text/plain"
	сontentTypeAppJSON   = "application/json"
	сontentTypeAppXGZIP  = "application/x-gzip"
)

func (h *Handler) GetShortLink(res http.ResponseWriter, req *http.Request) {
	userID, err := getUserID(req.Context())
	if err != nil {
		http.Error(res, "failded getting userID", http.StatusBadRequest)
		return
	}

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
	ident, err := h.services.GetIdent(req.Context(), string(body), userID)
	if err != nil {
		if !errors.Is(err, postgresstorage.ErrConflict) {
			http.Error(res, err.Error(), http.StatusInternalServerError)
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

func (h *Handler) GetFulLink(res http.ResponseWriter, req *http.Request) {
	link, err := h.services.GetFulLink(req.Context(), chi.URLParam(req, "ident"))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if link.DeletedFlag {
		http.Error(res, "resurs deleted", http.StatusGone)
		return
	}
	res.Header().Set("Location", link.FulLink)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) GetShortLinkByJSON(res http.ResponseWriter, req *http.Request) {
	userID, err := getUserID(req.Context())
	if err != nil {
		http.Error(res, "failded getting userID", http.StatusBadRequest)
		return
	}

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
	ident, err := h.services.GetIdent(req.Context(), request.URL, userID)
	if err != nil {
		if !errors.Is(err, postgresstorage.ErrConflict) {
			http.Error(res, err.Error(), http.StatusInternalServerError)
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

func (h *Handler) GetShortLinkByListJSON(res http.ResponseWriter, req *http.Request) {
	userID, err := getUserID(req.Context())
	if err != nil {
		http.Error(res, "failded getting userID", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get(сontentType)
	if !(ct == сontentTypeAppJSON || ct == сontentTypeAppXGZIP) {
		http.Error(res, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	var linkReq []dto.LinkListReq
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, "invalid body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &linkReq); err != nil {
		http.Error(res, "invalid format body", http.StatusBadRequest)
		return
	}

	limkResp, err := h.services.GetIdents(req.Context(), linkReq, userID)
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

func (h *Handler) GetLinksByUser(res http.ResponseWriter, req *http.Request) {
	userID, err := getUserID(req.Context())
	if err != nil {
		http.Error(res, "failded getting userID", http.StatusBadRequest)
		return
	}

	linksResp, err := h.services.GetLinksByUserID(req.Context(), userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(linksResp) == 0 {
		res.WriteHeader(http.StatusNoContent)
	}
	for i, v := range linksResp {
		linksResp[i].ShortURL = h.baseShortURL + "/" + v.ShortURL
	}
	response, err := json.Marshal(&linksResp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set(сontentType, сontentTypeAppJSON)
	res.WriteHeader(http.StatusOK)
	res.Write(response)
}

func (h *Handler) DeleteLinksByIdents(res http.ResponseWriter, req *http.Request) {
	userID, err := getUserID(req.Context())
	if err != nil {
		http.Error(res, "failded getting userID", http.StatusBadRequest)
		return
	}

	ct := req.Header.Get(сontentType)
	if !(ct == сontentTypeAppJSON || ct == сontentTypeAppXGZIP) {
		http.Error(res, "invalid Content-Type", http.StatusBadRequest)
		return
	}

	var request []string

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	can, err := h.services.CanDelete(req.Context(), userID, request...)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if !can {
		http.Error(res, "forbidden", http.StatusForbidden)
		return
	}

	h.delChan <- delMesage{
		idents: request,
	}

	res.WriteHeader(http.StatusAccepted)
}

func (h *Handler) flushMessagesDelete(stop <-chan bool) {
	ticker := time.NewTicker(5 * time.Second)
	identsBuf := make([]string, 0)
	for {
		select {
		case msg := <-h.delChan:
			identsBuf = append(identsBuf, msg.idents...)
		case <-ticker.C:
			if len(identsBuf) == 0 {
				continue
			}
			err := h.services.DeleteLinksByIdent(context.Background(), identsBuf...)
			if err != nil {
				logger.Log().Debug("cannot delete links")
				continue
			}
			identsBuf = identsBuf[:0]
		case <-stop:
			close(h.delChan)
			for msg := range h.delChan {
				identsBuf = append(identsBuf, msg.idents...)
			}
			if len(identsBuf) == 0 {
				return
			}
			err := h.services.DeleteLinksByIdent(context.Background(), identsBuf...)
			if err != nil {
				logger.Log().Debug("cannot delete links when stop")
				return
			}
		}
	}
}

func (h *Handler) FlushMessagesDeleteNow() {
	h.stopChan <- true
	close(h.stopChan) 
}

package handlers

import (
	"context"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/gzipmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/logmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type delMesage struct {
	idents []string
}
type Handler struct {
	services     *Service
	baseShortURL string
	delChan      chan delMesage
}

func NewHandler(services *Service, baseShortURL string) *Handler {
	h := &Handler{
		services:     services,
		baseShortURL: baseShortURL,
		delChan:      make(chan delMesage, 1024),
	}
	go h.flushMessagesDelete()
	return h

}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(logmiddleware.WithLogging)
	router.Use(gzipmiddleware.Decompress)
	router.Use(h.userIdentity)
	router.Use(h.setTokenID)
	router.Use(middleware.Compress(5, "application/json", "text/html"))
	router.Post("/", h.GetShortLink)
	router.Post("/api/shorten", h.GetShortLinkByJSON)
	router.Post("/api/shorten/batch", h.GetShortLinkByListJSON)
	router.Get("/{ident}", h.GetFulLink)
	router.Get("/api/user/urls", h.GetLinksByUser)
	router.Delete("/api/user/urls", h.DeleteLinksByIdents)
	return router
}

type Service struct {
	AuthService
	LinkService
}

func NewServices(linkStorage service.LinkStorage, userStorage service.UserStorage) *Service {
	return &Service{
		AuthService: service.NewAauthService(userStorage),
		LinkService: service.NewLinkService(linkStorage),
	}
}

type AuthService interface {
	ParseToken(accessToken string) (int32, bool, error)
	BuildJWTString(userID int32) (string, error)
	CreateUser(ctx context.Context) (int32, error)
}

type LinkService interface {
	GetFulLink(ctx context.Context, ident string) (domain.Link, error)
	GetIdent(ctx context.Context, fulLink string, userID int32) (string, error)
	GetIdents(ctx context.Context, linkReq []dto.LinkListReq, userID int32) ([]dto.LinkListRes, error)
	GenerateIdent(url string) string
	GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error)
	CanDelete(ctx context.Context, userID int32, idents ...string) (bool, error)
	DeleteLinksByIdent(ctx context.Context, idents ...string) error 
}

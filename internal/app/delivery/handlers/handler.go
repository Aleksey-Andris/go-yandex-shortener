package handlers

import (
	"context"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/gzipmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/logmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	services     *Service
	baseShortURL string
}

func NewHandler(services *Service, baseShortURL string) *Handler {
	return &Handler{
		services:     services,
		baseShortURL: baseShortURL}

}

func (h *Handler) InitRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(logmiddleware.WithLogging)
	router.Use(gzipmiddleware.Decompress)
	router.Use(h.userIdentity)
	router.Use(middleware.Compress(5, "application/json", "text/html"))
	router.Group(func(r chi.Router) {
		r.Use(h.setTokenID)
		r.Post("/", h.GetShortLink)
		r.Post("/api/shorten", h.GetShortLinkByJSON)
		r.Post("/api/shorten/batch", h.GetShortLinkByListJSON)
		r.Get("/{ident}", h.GetFulLink)
	})
	router.Get("/api/user/urls", h.GetLinksByUser)
	return router
}

type Service struct {
	AuthService
	LinkService
}

func NewServices(linkStorage service.LinkStorage, userStorage service.UserStorage, baseShortURL string) *Service {
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
	GetFulLink(ctx context.Context, ident string) (string, error)
	GetIdent(ctx context.Context, fulLink string, userID int32) (string, error)
	GetIdents(ctx context.Context, linkReq []dto.LinkListReq, userID int32) ([]dto.LinkListRes, error)
	GenerateIdent(url string) string
	GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error)
}

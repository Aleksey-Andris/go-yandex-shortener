// The handlers package contains settings for HTTP endpoints.
package handlers

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/gzipmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/middlware/logmiddleware"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
)

type delMesage struct {
	idents []string
}

// Handler - structure providing access to services(uscases)
type Handler struct {
	services     *Service
	baseShortURL string
	delChan      chan delMesage
	stopChan     chan bool
}

// NewHandler - constructor for Handler.
func NewHandler(services *Service, baseShortURL string) *Handler {
	h := &Handler{
		services:     services,
		baseShortURL: baseShortURL,
		delChan:      make(chan delMesage, 1),
		stopChan:     make(chan bool),
	}
	go h.flushMessagesDelete(h.stopChan)
	return h

}

// InitRouter - сreating endpoints and connecting middleware.
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

// Service - сomposite structure containing other services.
type Service struct {
	AuthService
	LinkService
}

// NewServices - constructor for Service.
func NewServices(linkStorage service.LinkStorage, userStorage service.UserStorage) *Service {
	return &Service{
		AuthService: service.NewAauthService(userStorage),
		LinkService: service.NewLinkService(linkStorage),
	}
}

// AuthService - interface describing the contract of authorization methods.
type AuthService interface {
	// ParseToken - GWT token parsing method, returns user ID, token validity and error.
	ParseToken(accessToken string) (int32, bool, error)
	// ParseToken - GWT token creating method.
	BuildJWTString(userID int32) (string, error)
	// ParseToken - user creating method.
	CreateUser(ctx context.Context) (int32, error)
}

// LinkService - interface describing the contract for working with link entities.
type LinkService interface {
	// GetFulLink - returns the full link by the shortened URL ident.
	GetFulLink(ctx context.Context, ident string) (domain.Link, error)
	// GetIdent - returns the short link by the full URL.
	GetIdent(ctx context.Context, fulLink string, userID int32) (string, error)
	// GetIdents - returns the short links by the full URLs.
	GetIdents(ctx context.Context, linkReq []dto.LinkListReq, userID int32) ([]dto.LinkListRes, error)
	// GenerateIdent - generates a unique shorts ident.
	GenerateIdent(url string) string
	// GetLinksByUserID - returns all user's links.
	GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error)
	// CanDelete - checks whether the user has the right to delete these links.
	CanDelete(ctx context.Context, userID int32, idents ...string) (bool, error)
	// DeleteLinksByIdent - removes links by their idents.
	DeleteLinksByIdent(ctx context.Context, idents ...string) error
}

package handlers

import (
	"context"
	"errors"
	"net/http"
)

const (
	authorizationHeader           = "Authorization"
	userCTX             YSContext = "YSUserID"
)

type YSContext string

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		cookieToken, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				next.ServeHTTP(res, req)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		tokenString := cookieToken.Value
		userID, valid, err := h.services.ParseToken(tokenString)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if userID == 0 {
			http.Error(res, "not authorization", http.StatusUnauthorized)
		}

		if !valid {
			next.ServeHTTP(res, req)
			return
		}

		request := req.WithContext(context.WithValue(req.Context(), userCTX, userID))

		next.ServeHTTP(res, request)
	})
}

func (h *Handler) setTokenID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		userID, err := getUserID(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if userID > 0 {
			next.ServeHTTP(res, req)
			return
		}
		userID, err = h.getNewUserID(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		tokenVal, err := h.services.BuildJWTString(userID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(res, &http.Cookie{
			Name:     "token",
			Value:    tokenVal,
			HttpOnly: true,
		})
		request := req.WithContext(context.WithValue(req.Context(), userCTX, userID))
		next.ServeHTTP(res, request)
	})
}

func (h *Handler) getNewUserID(ctx context.Context) (int32, error) {
	return h.services.CreateUser(ctx)

}

func getUserID(ctx context.Context) (int32, error) {
	id := ctx.Value(userCTX)
	if id == nil {
		return -1, nil
	}

	idInt, ok := id.(int32)
	if !ok {
		return -1, errors.New("user id of invalid type")
	}
	return idInt, nil
}

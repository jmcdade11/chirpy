package main

import (
	"net/http"

	"github.com/jmcdade11/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		unauthorized(w, "Couldn't find JWT")
		return
	}

	refreshToken, err := cfg.DB.GetRefreshToken(bearerToken)

	if err != nil {
		unauthorized(w, "Couldn't find refresh token")
		return
	}
	cfg.DB.DeleteRefreshToken(refreshToken.ID)

	w.WriteHeader(204)
}

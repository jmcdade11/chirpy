package main

import (
	"net/http"
	"time"

	"github.com/jmcdade11/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		unauthorized(w, "Couldn't find bearer token")
		return
	}

	refresh, err := cfg.DB.GetRefreshToken(bearerToken)
	if err != nil {
		unauthorized(w, "Couldn't find refresh token")
		return
	}
	now := time.Now().UTC()

	if now.Unix() > refresh.Expiration.Unix() {
		unauthorized(w, "Refresh token is expired")
		return
	}

	accessToken, err := auth.CreateJwt(refresh.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Token: accessToken,
	})
}

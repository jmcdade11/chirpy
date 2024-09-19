package main

import (
	"encoding/json"
	"net/http"

	"github.com/jmcdade11/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLoginCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type Response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		unauthorized(w, "")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		unauthorized(w, "")
		return
	}

	token, err := auth.CreateJwt(user.ID, cfg.jwtSecret, params.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}

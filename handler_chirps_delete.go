package main

import (
	"net/http"
	"strconv"

	"github.com/jmcdade11/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		unauthorized(w, "Couldn't find JWT")
		return
	}
	subject, err := auth.ValidateJwt(token, cfg.jwtSecret)
	if err != nil {
		unauthorized(w, "Couldn't validate JWT")
		return
	}
	userId, err := strconv.Atoi(subject)
	if err != nil {
		unauthorized(w, "Couldn't validate JWT")
		return
	}

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse chirp ID")
		return
	}
	chirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp")
		return
	}

	if chirp.AuthorID != userId {
		w.WriteHeader(403)
		return
	}
	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(204)
}

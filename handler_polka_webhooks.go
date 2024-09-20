package main

import (
	"encoding/json"
	"net/http"

	"github.com/jmcdade11/chirpy/internal/webhooks"
)

type Data struct {
	UserId int `json:"user_id"`
}

type Webhook struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	polkaKey, err := webhooks.GetPolkaKey(r.Header)
	if err != nil {
		unauthorized(w, "Invalid api key")
	}
	if polkaKey != cfg.polkaKey {
		unauthorized(w, "Invalid api key")
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	if _, err = cfg.DB.GetUser(params.Data.UserId); err != nil {
		w.WriteHeader(404)
		return
	}

	err = cfg.DB.EnableChirpyRed(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not enable ChirpyRed")
		return
	}

	w.WriteHeader(204)
}

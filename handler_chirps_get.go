package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	authorId := 0
	authorIdString := r.URL.Query().Get("author_id")
	if authorIdString != "" {
		authorId, err = strconv.Atoi(authorIdString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Malformed author id")
			return
		}
	}

	sortParam := r.URL.Query().Get("sort")
	if sortParam != "" && sortParam != "asc" && sortParam != "desc" {
		respondWithError(w, http.StatusBadRequest, "Malformed sort order")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		if authorId != 0 {
			if authorId == dbChirp.AuthorID {
				chirps = append(chirps, Chirp{
					ID:       dbChirp.ID,
					Body:     dbChirp.Body,
					AuthorID: dbChirp.AuthorID,
				})
			}
		} else {
			chirps = append(chirps, Chirp{
				ID:       dbChirp.ID,
				Body:     dbChirp.Body,
				AuthorID: dbChirp.AuthorID,
			})
		}
	}

	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpRetrieve(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse chirp ID")
		return
	}
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	if id < 1 || id > len(dbChirps) {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	chirp := dbChirps[id-1]

	respondWithJSON(w, http.StatusOK, chirp)
}

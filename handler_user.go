package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hendrixthecoder/url_shortener/internal/database"
)

func (appConfig *AppConfig) handlerShortenUrl(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		URL string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	user_id, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong.")
		return
	}

	if !isValidURL(params.URL) {
		respondWithError(w, http.StatusBadRequest, "Invalid URL passed.")
		return
	}

	short_url, err := generateUniqueShortCode(appConfig)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error generating short url, try again!")
		return
	}

	_, err = appConfig.DB.CreateNewShortURL(r.Context(), database.CreateNewShortURLParams{
		ID:        uuid.New(),
		UserID:    user_id,
		ShortUrl:  short_url,
		PlainUrl:  params.URL,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().AddDate(1, 0, 0),
	})

	if err != nil {
		log.Println("Error generating short url:", err)
		respondWithError(w, http.StatusBadRequest, "Error generating short url, try again!")
		return
	}

	respondWithJSON(w, http.StatusOK, parameters{URL: appConfig.AppURL + "/" + short_url})
}

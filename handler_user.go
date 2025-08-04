package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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
		respondWithError(w, 400, "Invalid request data")
		return
	}

	user_id_val := r.Context().Value(contextKey("user_id"))
	user_id_str, ok := user_id_val.(string)
	if !ok {
		log.Println("user_id missing from context or not a string")
		respondWithError(w, 401, "Unauthorized")
		return
	}

	user_id, err := uuid.Parse(user_id_str)
	if err != nil {
		log.Println("Invalid user_id")
		respondWithError(w, 400, "Invalid user.")
		return
	}

	if !isValidURL(params.URL) {
		respondWithError(w, 400, "Invalid URL passed.")
		return
	}

	short_url, err := generateUniqueShortCode(appConfig)
	if err != nil {
		respondWithError(w, 400, "Error generating short url, try again!")
		return
	}

	_, err = appConfig.DB.CreateNewShortURL(r.Context(), database.CreateNewShortURLParams{
		ID:        uuid.New(),
		UserID:    user_id,
		ShortUrl:  short_url,
		PlainUrl:  params.URL,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		log.Println("Error generating short url:", err)
		respondWithError(w, 400, "Error generating short url, try again!")
		return
	}

	respondWithJSON(w, 200, parameters{URL: appConfig.AppURL + "/" + short_url})
}

func (appConfig *AppConfig) handlerGetPlainUrl(w http.ResponseWriter, r *http.Request) {
	short_url := chi.URLParam(r, "short_url")

	if len(short_url) < 6 {
		respondWithError(w, 400, "Invalid short code passed!")
		return
	}

	record, err := appConfig.DB.GetURLEntryByShortURL(r.Context(), short_url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "Short URL not found")
			return
		}

		log.Println("Error fetching URL:", err)
		respondWithError(w, 500, "Error fetching URL")
		return
	}

	http.Redirect(w, r, record.PlainUrl, http.StatusFound)
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hendrixthecoder/url_shortener/internal/database"
)

type URL struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	ShortUrl  string    `json:"short_url"`
	PlainUrl  string    `json:"plain_url"`
}

func databaseURLToUrl(url database.Url) URL {
	return URL{
		ID:        url.ID,
		CreatedAt: url.CreatedAt,
		UpdatedAt: url.UpdatedAt,
		UserID:    url.UserID,
		ShortUrl:  url.ShortUrl,
		PlainUrl:  url.PlainUrl,
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

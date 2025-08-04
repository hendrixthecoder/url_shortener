package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hendrixthecoder/url_shortener/internal/database"
)

func (appConfig *AppConfig) handlerGetUserURLs(w http.ResponseWriter, r *http.Request) {
	userUUID, err := getUserIDFromContext(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong.")
		return
	}

	urls, err := appConfig.DB.GetURLEntriesByUserID(r.Context(), userUUID)
	if err != nil {
		log.Println("Error fetching user url entries", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong!")
		return
	}

	respondWithJSON(w, http.StatusOK, dtoSliceSerializer(urls, databaseURLToUrl))
}

func (appConfig *AppConfig) handlerGetPlainUrl(w http.ResponseWriter, r *http.Request) {
	short_url := chi.URLParam(r, "short_url")

	if len(short_url) < 6 {
		respondWithError(w, http.StatusBadRequest, "Invalid short code passed!")
		return
	}

	record, err := appConfig.DB.GetURLEntryByShortURL(r.Context(), short_url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Short URL not found")
			return
		}

		log.Println("Error fetching URL:", err)
		respondWithError(w, http.StatusInternalServerError, "Error fetching URL")
		return
	}

	ctx := context.Background()

	go func() {
		_, err := appConfig.DB.CreateURLAnalyticRecord(ctx, database.CreateURLAnalyticRecordParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			ShortUrl:  short_url,
			Ip:        getIP(r),
			Referer:   r.Referer(),
			UserAgent: r.Header.Get("User-Agent"),
		})

		if err != nil {
			log.Println("Error creating url analytic record:", err)
		}
	}()

	http.Redirect(w, r, record.PlainUrl, http.StatusFound)
}

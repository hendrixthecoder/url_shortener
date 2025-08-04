package main

import (
	"log"
	"net/http"
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

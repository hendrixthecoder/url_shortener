package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (appConfig *AppConfig) handlerShortenUrl(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		URL string `json:"url"`
	}

	email := r.Context().Value(userContextKey).(string)
	fmt.Println(email)

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, 400, "Invalid request data")
		return
	}

	respondWithJSON(w, 200, parameters{URL: params.URL})
}

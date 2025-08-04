package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hendrixthecoder/url_shortener/internal/database"
)

func (appConfig *AppConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error creating account, try again.")
		return
	}

	_, err = appConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		log.Println("Error creating user:", err)
		respondWithError(w, http.StatusBadRequest, "Could not create user.")
		return
	}

	respondWithJSON(w, http.StatusCreated, "Account created succesfully!")
}

func (appConfig *AppConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := appConfig.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil || !checkPasswordHash(params.Password, user.Password) {
		log.Println("Error decoding request body:", err)
		respondWithError(w, http.StatusBadRequest, "Invalid credentials.")
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error logging in, try again!: %v", err))
		return
	}

	session.Values["user_email"] = user.Email
	session.Values["user_id"] = user.ID.String()
	session.Options.HttpOnly = appConfig.Env == "production"
	session.Options.Secure = appConfig.Env == "production"
	err = session.Save(r, w)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error logging in, try again!: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, "Logged in succesfully!")
}

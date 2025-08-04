package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortCodeLength = 6

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes), nil
}

func checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func isValidURL(u string) bool {
	parsed, err := url.ParseRequestURI(u)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}
	return true
}

func generateShortCode() string {
	b := make([]byte, shortCodeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateUniqueShortCode(appConfig *AppConfig) (string, error) {
	for i := 0; i < 10; i++ {
		code := generateShortCode()

		_, err := appConfig.DB.GetURLEntryByShortURL(context.Background(), code)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return code, nil
			}

			log.Println("Could not fetch entry by short url code:", err)
			continue
		}

	}
	return "", errors.New("could not generate unique short code after 10 attempts")
}

func getUserIDFromContext(r *http.Request) (uuid.UUID, error) {
	user_id_val := r.Context().Value(contextKey("user_id"))
	user_id_str, ok := user_id_val.(string)
	if !ok {
		log.Println("user_id missing from context or not a string")
		return uuid.Nil, errors.New("UserID missing from string.")
	}

	user_id, err := uuid.Parse(user_id_str)
	if err != nil {
		log.Println("Invalid user_id")
		return uuid.Nil, errors.New("Invalid user ID")
	}

	return user_id, nil
}

// Serializer func to serialize slice of data from Goose types to JSON-normalized type.
func dtoSliceSerializer[T any, DTO any](data []T, serializer func(T) DTO) []DTO {
	serialized := make([]DTO, len(data))

	for idx, feed := range data {
		serialized[idx] = serializer(feed)
	}

	return serialized
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

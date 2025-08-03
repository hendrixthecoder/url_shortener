package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"net/url"
	"time"

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

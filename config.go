package main

import (
	"log"
	"os"
)

type Config struct {
	dbString   string
	portString string
	env        string
	appUrl     string
	csrfKey    []byte
}

func LoadConfig() Config {
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not provided in .env")
	}

	dbString := os.Getenv("DB_URL")
	if dbString == "" {
		log.Fatal("DB_URL not provided in .env")
	}

	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("DB_URL not provided in .env")
	}

	appUrl := os.Getenv("APP_URL")
	if appUrl == "" {
		log.Fatal("APP_URL not provided in .env")
	}

	csrfKey := os.Getenv("CSRF_KEY")
	if appUrl == "" {
		log.Fatal("CSRF_KEY not provided in .env")
	}

	return Config{
		dbString:   dbString,
		portString: portString,
		env:        env,
		appUrl:     appUrl,
		csrfKey:    []byte(csrfKey),
	}
}

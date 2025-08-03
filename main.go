package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/hendrixthecoder/url_shortener/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type AppConfig struct {
	DB     *database.Queries
	Env    string
	AppURL string
}

func main() {
	godotenv.Load()

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

	db, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatal("Failed to open Postgres driver:", err)
	}

	conn := database.New(db)

	InitRedisStore()

	appConfig := &AppConfig{
		DB:     conn,
		Env:    env,
		AppURL: appUrl,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "http://"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/users", appConfig.handlerCreateUser)
	v1Router.Post("/login", appConfig.handlerLoginUser)

	v1Router.Post("/shorten", appConfig.authMiddleware(appConfig.handlerShortenUrl))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Println("Server starting on port:", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Could not spin up server:", err)
	}

	log.Println("Server running on port:", portString)
}

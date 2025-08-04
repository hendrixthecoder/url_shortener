package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
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

	config := LoadConfig()

	db, err := sql.Open("postgres", config.dbString)
	if err != nil {
		log.Fatal("Failed to open Postgres driver:", err)
	}

	conn := database.New(db)

	appConfig := &AppConfig{
		DB:     conn,
		Env:    config.env,
		AppURL: config.appUrl,
	}

	ctx := context.Background()
	go StartURLCleaner(appConfig, ctx)
	InitRedisStore()

	router := chi.NewRouter()

	csrfMiddleware := csrf.Protect(
		config.csrfKey,
		csrf.Secure(config.env == "production"),
	)

	router.Use(middleware.Logger)
	router.Use(csrfMiddleware)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "http://"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/register", appConfig.handlerCreateUser)
	v1Router.Post("/login", appConfig.handlerLoginUser)

	v1Router.Post("/shorten", appConfig.authMiddleware(appConfig.handlerShortenUrl))

	v1Router.Get("/urls", appConfig.authMiddleware(appConfig.handlerGetUserURLs))

	v1Router.Get("/{short_url}", appConfig.handlerGetPlainUrl)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + config.portString,
	}

	log.Println("Server starting on port:", config.portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Could not spin up server:", err)
	}
}

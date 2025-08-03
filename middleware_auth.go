package main

import (
	"context"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user_email")

func (appConfig *AppConfig) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-session")
		email, ok := session.Values["user_email"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

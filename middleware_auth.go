package main

import (
	"context"
	"net/http"
)

type contextKey string

func (appConfig *AppConfig) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-session")
		user_id, ok1 := session.Values["user_id"].(string)
		email, ok2 := session.Values["user_email"].(string)

		if !ok1 || user_id == "" || !ok2 || email == "" {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), contextKey("user_id"), user_id)
		ctx = context.WithValue(ctx, contextKey("email"), email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

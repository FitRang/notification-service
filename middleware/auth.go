package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "user-id"
const EmailKey contextKey = "user-email"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("x-user-id")
		emailID := r.Header.Get("x-user-email")
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, EmailKey, emailID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

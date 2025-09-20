package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go-version/internal/api/auth"
	"go-version/internal/contextkeys"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx context.Context) (func(http.Handler) http.Handler, error) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := auth.ParseToken(tokenString)
			if err != nil || !token.Valid {
				fmt.Println("Token parse error:", err)
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract `sub` claim (userId)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				http.Error(w, "missing sub claim", http.StatusUnauthorized)
				return
			}

			// Attach userID to request context
			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

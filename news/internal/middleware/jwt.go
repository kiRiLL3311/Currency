package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kiRiLL3311/Currency/news/internal/config"
)

type contextKey string

const UserIDKey contextKey = "userID"

const (
	TokenIssuer   = "currency-auth-service"
	TokenAudience = "currency-api"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		const bearer = "Bearer "
		if len(authHeader) < len(bearer) || authHeader[:len(bearer)] != bearer {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len(bearer):]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(config.Get("JWT_SECRET")), nil
			},
			jwt.WithIssuer(TokenIssuer),
			jwt.WithAudience(TokenAudience),
		)

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

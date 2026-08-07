
package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kiRiLL3311/Currency/auth/internal/config"
)

const (
	TokenIssuer   = "currency-auth-service"
	TokenAudience = "currency-api"
)

// Issuer → identifies who issued the token.

// Subject → the user the token belongs to.

// Audience → which service should accept the token.

// IssuedAt → when it was created.

// NotBefore → prevents use before its issue time.

// ExpiresAt → expiration.
func GenerateJWT(userID int, email string) (string, error) {

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TokenIssuer,
			Subject:   email,
			Audience:  []string{TokenAudience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Get("JWT_SECRET")))
}

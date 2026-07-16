package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/kiRiLL3311/Currency/auth/internal/config"
)

func GenerateRefreshToken() (string, error) {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func RefreshExpiry() time.Time {

	duration, err := time.ParseDuration(
		config.Get("REFRESH_EXPIRES_IN"),
	)

	if err != nil {
		duration = 720 * time.Hour
	}

	return time.Now().Add(duration)
}

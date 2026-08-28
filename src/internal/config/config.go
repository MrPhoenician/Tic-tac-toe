package config

import (
	"os"
	"time"
)

func GetDataBaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	return "postgres://postgres:1@localhost:5432/tic_tac_toe?sslmode=disable"
}

func GetJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}

	return "tic-tac-dev-secret"
}

func GetAccessTokenTTL() time.Duration {
	if value := os.Getenv("JWT_ACCESS_TTL"); value != "" {
		if ttl, err := time.ParseDuration(value); err == nil {
			return ttl
		}
	}

	return 15 * time.Minute
}

func GetRefreshTokenTTL() time.Duration {
	if value := os.Getenv("JWT_REFRESH_TTL"); value != "" {
		if ttl, err := time.ParseDuration(value); err == nil {
			return ttl
		}
	}

	return 7 * 24 * time.Hour
}

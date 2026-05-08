package app

import (
	"context"
	"errors"
	"strings"
)

type AuthTokens struct {
	Type         string
	AccessToken  string
	RefreshToken string
}

type AuthService interface {
	SignUp(ctx context.Context, login, password string) (*User, error)
	SignIn(ctx context.Context, login, password string) (*AuthTokens, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*AuthTokens, error)
	RefreshRefreshToken(ctx context.Context, refreshToken string) (*AuthTokens, error)
	GetUserByAccessToken(ctx context.Context, accessToken string) (*User, error)
}

func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return "", errors.New("access token is required")
	}

	return token, nil
}

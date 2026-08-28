package middleware

import (
	"context"
	"net/http"
	"tic-tac/internal/app"
	"tic-tac/internal/contextkeys"
	"tic-tac/internal/logger"
)

type UserAuthenticator struct {
	authService app.AuthService
	log         logger.Logger
	next        http.Handler
}

func NewUserAuthenticator(
	authService app.AuthService,
	log logger.Logger,
	next http.Handler,
) *UserAuthenticator {
	return &UserAuthenticator{
		authService: authService,
		log:         log,
		next:        next,
	}
}

func (m *UserAuthenticator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	accessToken, err := app.ExtractBearerToken(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := m.authService.GetUserByAccessToken(r.Context(), accessToken)
	if err != nil {
		m.log.Warn("Unauthorized access attempt", "error", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), contextkeys.UserID, user.ID)
	ctx = context.WithValue(ctx, contextkeys.UserLogin, user.Login)
	r = r.WithContext(ctx)

	m.next.ServeHTTP(w, r)
}

package api

import (
	"net/http"
	"tic-tac/api/middleware"
	"tic-tac/internal/app"
	"tic-tac/internal/logger"
)

type Router struct {
	authHandler *AuthHandler
	gameHandler *GameHandler
	authService app.AuthService
	log         logger.Logger
}

func NewRouter(
	authHandler *AuthHandler,
	gameHandler *GameHandler,
	authService app.AuthService,
	log logger.Logger,
) *Router {
	return &Router{
		authHandler: authHandler,
		gameHandler: gameHandler,
		authService: authService,
		log:         log,
	}
}

func (r *Router) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", r.authHandler.SignUp)
	mux.HandleFunc("POST /signin", r.authHandler.SignIn)
	mux.HandleFunc("POST /auth/access", r.authHandler.RefreshAccessToken)
	mux.HandleFunc("POST /auth/refresh", r.authHandler.RefreshRefreshToken)

	secureMe := middleware.NewUserAuthenticator(r.authService, r.log, http.HandlerFunc(r.authHandler.Me))
	mux.Handle("GET /me", secureMe)

	secureCreate := middleware.NewUserAuthenticator(r.authService, r.log, http.HandlerFunc(r.gameHandler.CreateGame))
	mux.Handle("POST /games", secureCreate)

	secureMove := middleware.NewUserAuthenticator(r.authService, r.log, http.HandlerFunc(r.gameHandler.HandleMove))
	mux.Handle("POST /game/{id}", secureMove)

	secureGet := middleware.NewUserAuthenticator(r.authService, r.log, http.HandlerFunc(r.gameHandler.GetGame))
	mux.Handle("GET /game/{id}", secureGet)

	secureList := middleware.NewUserAuthenticator(r.authService, r.log, http.HandlerFunc(r.gameHandler.ListGames))
	mux.Handle("GET /games", secureList)

	secureCompleted := middleware.NewUserAuthenticator(r.authService, r.log, http.HandlerFunc(r.gameHandler.ListCompletedGames))
	mux.Handle("GET /games/completed", secureCompleted)

	secureLeaderboard := middleware.NewUserAuthenticator(r.authService, r.log, http.HandlerFunc(r.gameHandler.GetLeaderboard))
	mux.Handle("GET /leaderboard", secureLeaderboard)

	joinHandler := middleware.NewUserAuthenticator(r.authService, r.log, http.HandlerFunc(r.gameHandler.JoinGame))
	mux.Handle("POST /game/{id}/join", joinHandler)

	return mux
}

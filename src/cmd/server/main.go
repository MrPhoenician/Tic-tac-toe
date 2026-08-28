package main

import (
	"net/http"
	"tic-tac/api"
	"tic-tac/internal/app"
	"tic-tac/internal/di"
	"tic-tac/internal/logger"

	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(
		di.Module,
		fx.Invoke(func(
			authHandler *api.AuthHandler,
			gameHandler *api.GameHandler,
			authService app.AuthService,
			log logger.Logger,
		) {
			router := api.NewRouter(authHandler, gameHandler, authService, log).Routes()

			log.Info("Server starting on http://localhost:8080")
			if err := http.ListenAndServe(":8080", router); err != nil {
				log.Error("Server failed", "error", err)
			}
		}),
	)
	fxApp.Run()
}

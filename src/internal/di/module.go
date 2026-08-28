package di

import (
	"context"
	"tic-tac/api"
	"tic-tac/internal/config"
	"tic-tac/internal/infrastructure"
	"tic-tac/internal/infrastructure/db"
	"tic-tac/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewPostgresPool(lc fx.Lifecycle, log logger.Logger) (*pgxpool.Pool, error) {
	ctx := context.Background()
	pg, err := db.NewPostgresDB(ctx, config.GetDataBaseURL())
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("Closing PostgreSQL connection")
			pg.Close()
			return nil
		},
	})

	return pg.Pool(), err
}

var Module = fx.Options(
	fx.Provide(logger.NewSlogLogger),
	fx.Provide(NewPostgresPool),
	fx.Provide(infrastructure.NewGameRepository),
	fx.Provide(infrastructure.NewUserRepository),
	fx.Provide(infrastructure.NewJwtProvider),
	fx.Provide(infrastructure.NewGameServiceImpl),
	fx.Provide(infrastructure.NewAuthService),
	fx.Provide(api.NewGameHandler),
	fx.Provide(api.NewAuthHandler),
)

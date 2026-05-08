package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, connString string) (*PostgresDB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("✅ PostgresSQL connected")
	return &PostgresDB{pool: pool}, nil
}

func (db *PostgresDB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *PostgresDB) Close() {
	db.pool.Close()
}

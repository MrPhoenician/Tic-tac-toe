package infrastructure

import (
	"context"
	"fmt"
	"tic-tac/internal/infrastructure/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameRepository struct {
	db *pgxpool.Pool
}

func NewGameRepository(db *pgxpool.Pool) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Save(ctx context.Context, game *model.Game) error {
	query := `INSERT INTO games (id, board, status, player_x, player_o, is_computer, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          ON CONFLICT (id) DO UPDATE
	          SET board = $2, status = $3, player_x = $4, player_o = $5, is_computer = $6`

	_, err := r.db.Exec(ctx, query,
		game.ID,
		game.Board,
		game.Status,
		game.PlayerX,
		game.PlayerO,
		game.IsComputer,
		game.CreatedAt,
	)
	return err
}

func (r *GameRepository) FindByID(ctx context.Context, id string) (*model.Game, error) {
	query := `SELECT id, board, status, player_x, player_o, is_computer, created_at FROM games WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var game model.Game
	var playerX, playerO *string

	err := row.Scan(
		&game.ID,
		&game.Board,
		&game.Status,
		&playerX,
		&playerO,
		&game.IsComputer,
		&game.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("game not found")
		}
		return nil, err
	}

	game.PlayerX = playerX
	game.PlayerO = playerO
	return &game, nil
}

func (r *GameRepository) FindByStatus(ctx context.Context, status string) ([]*model.Game, error) {
	query := `SELECT id, board, status, player_x, player_o, is_computer, created_at FROM games WHERE status = $1`
	rows, err := r.db.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*model.Game
	for rows.Next() {
		var game model.Game
		var playerX, playerO *string
		err := rows.Scan(&game.ID, &game.Board, &game.Status, &playerX, &playerO, &game.IsComputer, &game.CreatedAt)
		if err != nil {
			return nil, err
		}
		game.PlayerX = playerX
		game.PlayerO = playerO
		games = append(games, &game)
	}

	return games, rows.Err()
}

func (r *GameRepository) FindCompletedByUserID(ctx context.Context, userID string) ([]*model.Game, error) {
	query := `
		SELECT id, board, status, player_x, player_o, is_computer, created_at
		FROM games
		WHERE (
			(status = $2 AND player_x = $1::uuid)
			OR (status = $3 AND player_o = $1::uuid)
			OR (status = $4 AND (player_x = $1::uuid OR player_o = $1::uuid))
		)
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID, model.StatusXWon, model.StatusOWon, model.StatusDraw)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*model.Game
	for rows.Next() {
		var game model.Game
		var playerX, playerO *string
		if err := rows.Scan(&game.ID, &game.Board, &game.Status, &playerX, &playerO, &game.IsComputer, &game.CreatedAt); err != nil {
			return nil, err
		}
		game.PlayerX = playerX
		game.PlayerO = playerO
		games = append(games, &game)
	}

	return games, rows.Err()
}

func (r *GameRepository) GetTopPlayers(ctx context.Context, limit int) ([]*model.LeaderboardEntry, error) {
	query := `
		WITH user_game_stats AS (
			SELECT
				u.id AS user_id,
				u.login AS login,
				SUM(CASE
					WHEN (g.status = $1 AND g.player_x = u.id) OR (g.status = $2 AND g.player_o = u.id) THEN 1
					ELSE 0
				END) AS wins,
				SUM(CASE
					WHEN (g.status = $1 AND g.player_o = u.id) OR (g.status = $2 AND g.player_x = u.id) THEN 1
					WHEN g.status = $3 AND (g.player_x = u.id OR g.player_o = u.id) THEN 1
					ELSE 0
				END) AS non_wins
			FROM users u
			LEFT JOIN games g ON g.player_x = u.id OR g.player_o = u.id
			GROUP BY u.id, u.login
		)
		SELECT
			user_id::text,
			login,
			CASE
				WHEN non_wins = 0 THEN wins::float8
				ELSE wins::float8 / non_wins
			END AS win_ratio
		FROM user_game_stats
		WHERE wins > 0 OR non_wins > 0
		ORDER BY win_ratio DESC, login ASC
		LIMIT $4`

	rows, err := r.db.Query(ctx, query, model.StatusXWon, model.StatusOWon, model.StatusDraw, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*model.LeaderboardEntry
	for rows.Next() {
		entry := &model.LeaderboardEntry{}
		if err := rows.Scan(&entry.UserID, &entry.Login, &entry.WinRatio); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

package app

import (
	"context"
	"errors"
)

var (
	ErrGameAlreadyOver = errors.New("game is already over")
	ErrBoardTempered   = errors.New("board has been tampered with")
)

type GameService interface {
	CreateGame(ctx context.Context, userID string, vsComputer bool) (*Game, error)
	JoinGame(ctx context.Context, gameID, userID string) (*Game, error)
	GetGame(ctx context.Context, gameID string) (*Game, error)
	ListAvailableGames(ctx context.Context) ([]*Game, error)
	GetCompletedGamesByUser(ctx context.Context, userID string) ([]*Game, error)
	GetTopPlayers(ctx context.Context, limit int) ([]*LeaderboardEntry, error)
	ProcessPlayerMove(ctx context.Context, gameID string, newBoard [3][3]int) (*Game, error)
	IsGameOver(game *Game) (bool, GameStatus)
	ValidateBoard(original, current *Game) error
}

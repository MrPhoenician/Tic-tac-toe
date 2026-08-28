package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"tic-tac/internal/app"
	"tic-tac/internal/contextkeys"
	"tic-tac/internal/infrastructure/model"
)

type GameServiceImpl struct {
	gameRepo *GameRepository
}

var _ app.GameService = (*GameServiceImpl)(nil)

func NewGameServiceImpl(gameRepo *GameRepository) app.GameService {
	return &GameServiceImpl{gameRepo: gameRepo}
}

func (s *GameServiceImpl) ProcessPlayerMove(ctx context.Context, gameID string, newBoard [3][3]int) (*app.Game, error) {
	dbOrig, err := s.gameRepo.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	original := ToDomain(dbOrig)
	current := &app.Game{
		ID:         gameID,
		Board:      newBoard,
		Status:     original.Status,
		PlayerX:    original.PlayerX,
		PlayerO:    original.PlayerO,
		IsComputer: original.IsComputer,
		CreatedAt:  original.CreatedAt,
	}

	userIDVal := ctx.Value(contextkeys.UserID)
	if userIDVal == nil {
		return nil, errors.New("user not authenticated")
	}
	userID, ok := userIDVal.(string)
	if !ok {
		return nil, errors.New("invalid user ID format")
	}

	if over, _ := s.IsGameOver(original); over {
		return nil, errors.New("game is already over")
	}
	if !s.canUserMove(original, userID) {
		return nil, errors.New("it's not your turn or you are not in this game")
	}
	if err := s.ValidateBoard(original, current); err != nil {
		return nil, err
	}
	if err := validatePlayerMove(original, current); err != nil {
		return nil, err
	}

	if over, status := s.IsGameOver(current); over {
		current.Status = status
		dbUpdated := ToInfrastructure(current)

		if err := s.gameRepo.Save(ctx, dbUpdated); err != nil {
			return nil, err
		}
		return current, nil
	}

	if current.IsComputer {
		computerMove := app.Minimax(current.Board, 0, false)
		current.Board[computerMove.Row][computerMove.Col] = -1

		if over, status := s.IsGameOver(current); over {
			current.Status = status
		} else {
			current.Status = app.StatusPlayerX
		}

		dbUpdated := ToInfrastructure(current)
		if err := s.gameRepo.Save(ctx, dbUpdated); err != nil {
			return nil, err
		}

		return current, nil
	}

	var nextStatus app.GameStatus
	if original.Status == app.StatusPlayerX {
		nextStatus = app.StatusPlayerO
	} else {
		nextStatus = app.StatusPlayerX
	}

	current.Status = nextStatus
	dbUpdated := ToInfrastructure(current)

	if err := s.gameRepo.Save(ctx, dbUpdated); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *GameServiceImpl) canUserMove(game *app.Game, userID string) bool {
	isPlayerX := game.PlayerX != nil && *game.PlayerX == userID
	isPlayerO := game.PlayerO != nil && *game.PlayerO == userID

	if !isPlayerX && !isPlayerO {
		return false
	}
	if game.IsComputer {
		return isPlayerX // только X играет против компьютера
	}
	if game.Status == app.StatusPlayerX {
		return isPlayerX
	}
	if game.Status == app.StatusPlayerO {
		return isPlayerO
	}

	return false
}

func (s *GameServiceImpl) ValidateBoard(original, current *app.Game) error {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if original.Board[i][j] != 0 && original.Board[i][j] != current.Board[i][j] {
				return app.ErrBoardTempered
			}
		}
	}
	return nil
}

func (s *GameServiceImpl) CreateGame(ctx context.Context, userID string, vsComputer bool) (*app.Game, error) {
	dbGame := model.NewGame()
	dbGame.PlayerX = &userID
	dbGame.IsComputer = vsComputer

	if vsComputer {
		dbGame.Status = model.StatusPlayerX
	} else {
		dbGame.Status = model.StatusWaiting
	}

	if err := s.gameRepo.Save(ctx, dbGame); err != nil {
		return nil, err
	}

	return ToDomain(dbGame), nil
}

func validatePlayerMove(original, current *app.Game) error {
	var expectedValue int
	if original.Status == app.StatusPlayerX {
		expectedValue = 1
	} else if original.Status == app.StatusPlayerO {
		expectedValue = -1
	} else {
		return errors.New("invalid game status for move")
	}

	diffCount := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if original.Board[i][j] != current.Board[i][j] {
				if original.Board[i][j] != 0 {
					return errors.New("cannot modify existing moves")
				}
				if current.Board[i][j] != expectedValue {
					return fmt.Errorf("expected move value %d for current turn", expectedValue)
				}
				diffCount++
			}
		}
	}

	if diffCount == 0 {
		return errors.New("no move made")
	}
	if diffCount > 1 {
		return errors.New("only one move allowed per turn")
	}

	return nil
}

func (s *GameServiceImpl) IsGameOver(game *app.Game) (bool, app.GameStatus) {
	return app.IsGameOver(game.Board)
}

func (s *GameServiceImpl) JoinGame(ctx context.Context, gameID, userID string) (*app.Game, error) {
	dbGame, err := s.gameRepo.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if dbGame.Status != model.StatusWaiting {
		return nil, errors.New("game is not waiting for players")
	}
	if dbGame.PlayerX != nil && *dbGame.PlayerX == userID {
		return nil, errors.New("you are already in this game")
	}

	dbGame.PlayerO = &userID
	dbGame.Status = model.StatusPlayerX
	if err := s.gameRepo.Save(ctx, dbGame); err != nil {
		return nil, err
	}

	return ToDomain(dbGame), nil
}

func (s *GameServiceImpl) GetGame(ctx context.Context, gameID string) (*app.Game, error) {
	dbGame, err := s.gameRepo.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	return ToDomain(dbGame), nil
}

func (s *GameServiceImpl) ListAvailableGames(ctx context.Context) ([]*app.Game, error) {
	dbGames, err := s.gameRepo.FindByStatus(ctx, string(model.StatusWaiting))
	if err != nil {
		return nil, err
	}

	var games []*app.Game
	for _, dbGame := range dbGames {
		games = append(games, ToDomain(dbGame))
	}
	return games, nil
}

func (s *GameServiceImpl) GetCompletedGamesByUser(ctx context.Context, userID string) ([]*app.Game, error) {
	dbGames, err := s.gameRepo.FindCompletedByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	games := make([]*app.Game, 0, len(dbGames))
	for _, dbGame := range dbGames {
		games = append(games, ToDomain(dbGame))
	}

	return games, nil
}

func (s *GameServiceImpl) GetTopPlayers(ctx context.Context, limit int) ([]*app.LeaderboardEntry, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}

	entries, err := s.gameRepo.GetTopPlayers(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*app.LeaderboardEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, &app.LeaderboardEntry{
			UserID:   entry.UserID,
			Login:    entry.Login,
			WinRatio: entry.WinRatio,
		})
	}

	return result, nil
}

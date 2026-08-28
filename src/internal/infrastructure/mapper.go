package infrastructure

import (
	"tic-tac/internal/app"
	"tic-tac/internal/infrastructure/model"
)

func ToDomain(dbGame *model.Game) *app.Game {
	if dbGame == nil {
		return nil
	}

	var playerX, playerO *string
	if dbGame.PlayerX != nil {
		tmp := *dbGame.PlayerX
		playerX = &tmp
	}
	if dbGame.PlayerO != nil {
		tmp := *dbGame.PlayerO
		playerO = &tmp
	}

	var board app.Board
	for i := 0; i < 3 && i < len(dbGame.Board); i++ {
		for j := 0; j < 3 && j < len(dbGame.Board[i]); j++ {
			board[i][j] = dbGame.Board[i][j]
		}
	}

	return &app.Game{
		ID:         dbGame.ID,
		Board:      board,
		Status:     app.GameStatus(dbGame.Status),
		PlayerX:    playerX,
		PlayerO:    playerO,
		IsComputer: dbGame.IsComputer,
		CreatedAt:  dbGame.CreatedAt,
	}
}

func ToInfrastructure(domainGame *app.Game) *model.Game {
	if domainGame == nil {
		return nil
	}

	var playerX, playerO *string
	if domainGame.PlayerX != nil {
		tmp := *domainGame.PlayerX
		playerX = &tmp
	}
	if domainGame.PlayerO != nil {
		tmp := *domainGame.PlayerO
		playerO = &tmp
	}

	board := make([][]int, 3)
	for i := 0; i < 3; i++ {
		board[i] = make([]int, 3)
		for j := 0; j < 3; j++ {
			board[i][j] = domainGame.Board[i][j]
		}
	}

	return &model.Game{
		ID:         domainGame.ID,
		Board:      board,
		Status:     model.GameStatus(domainGame.Status),
		PlayerX:    playerX,
		PlayerO:    playerO,
		IsComputer: domainGame.IsComputer,
		CreatedAt:  domainGame.CreatedAt,
	}
}

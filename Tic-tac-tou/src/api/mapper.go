package api

import (
	"tic-tac/api/model"
	"tic-tac/internal/app"
)

func ToDomain(req *model.GameRequest) *app.Game {
	var board app.Board
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = req.Board[i][j]
		}
	}

	return &app.Game{
		ID:    req.ID,
		Board: board,
	}
}

func ToWeb(game *app.Game) model.GameResponse {
	if game == nil {
		return model.GameResponse{}
	}

	var playerX, playerO *string
	if game.PlayerX != nil {
		tmp := *game.PlayerX
		playerX = &tmp
	}
	if game.PlayerO != nil {
		tmp := *game.PlayerO
		playerO = &tmp
	}

	var board [3][3]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			board[i][j] = game.Board[i][j]
		}
	}

	return model.GameResponse{
		ID:         game.ID,
		Board:      board,
		Status:     string(game.Status),
		PlayerX:    playerX,
		PlayerO:    playerO,
		IsComputer: game.IsComputer,
		CreatedAt:  game.CreatedAt,
	}
}

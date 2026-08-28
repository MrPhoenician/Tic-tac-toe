package app

import "time"

type GameStatus string

const (
	StatusWaiting GameStatus = "waiting"
	StatusPlayerX GameStatus = "player_x"
	StatusPlayerO GameStatus = "player_o"
	StatusXWon    GameStatus = "x_won"
	StatusOWon    GameStatus = "o_won"
	StatusDraw    GameStatus = "draw"
)

type Game struct {
	ID         string
	Board      Board
	Status     GameStatus
	PlayerX    *string
	PlayerO    *string
	IsComputer bool
	CreatedAt  time.Time
}

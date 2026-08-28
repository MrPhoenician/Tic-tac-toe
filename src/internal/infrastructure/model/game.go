package model

import (
	"time"

	"github.com/google/uuid"
)

type GameStatus string

const (
	StatusWaiting GameStatus = "waiting"  // ожидание второго игрока
	StatusPlayerX GameStatus = "player_x" // ходит игрок X
	StatusPlayerO GameStatus = "player_o" // ходит игрок O
	StatusXWon    GameStatus = "x_won"    // победил X
	StatusOWon    GameStatus = "o_won"    // победил O
	StatusDraw    GameStatus = "draw"     // ничья
)

type Game struct {
	ID         string     `db:"id"`
	Board      [][]int    `db:"board"`
	Status     GameStatus `db:"status"`
	PlayerX    *string    `db:"player_x"`
	PlayerO    *string    `db:"player_o"`
	IsComputer bool       `db:"is_computer"`
	CreatedAt  time.Time  `db:"created_at"`
}

func NewGame() *Game {
	id := uuid.New().String()
	board := [][]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	return &Game{
		ID:         id,
		Board:      board,
		Status:     StatusWaiting,
		IsComputer: false,
		CreatedAt:  time.Now(),
	}
}

package model

import "time"

type GameRequest struct {
	ID    string
	Board [3][3]int
}

type GameResponse struct {
	ID         string    `json:"id"`
	Board      [3][3]int `json:"board"`
	Status     string    `json:"status"`
	PlayerX    *string   `json:"player_x,omitempty"`
	PlayerO    *string   `json:"player_o,omitempty"`
	IsComputer bool      `json:"is_computer"`
	CreatedAt  time.Time `json:"created_at"`
}

type LeaderboardResponse struct {
	UserID   string  `json:"user_id"`
	Login    string  `json:"login"`
	WinRatio float64 `json:"win_ratio"`
}

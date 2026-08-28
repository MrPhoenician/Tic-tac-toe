package app

type Move struct {
	Row, Col int
	Score    int
}

func Minimax(board Board, depth int, isMaximizing bool) Move {
	gameEnded, result := IsGameOver(board)

	if gameEnded {
		if result == StatusXWon {
			return Move{Score: 10 - depth}
		} else if result == StatusOWon {
			return Move{Score: depth - 10}
		} else {
			return Move{Score: 0}
		}
	}

	if isMaximizing {
		best := Move{Score: -1000}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == 0 {
					board[i][j] = 1
					move := Minimax(board, depth+1, false)
					board[i][j] = 0
					if move.Score > best.Score {
						best.Score = move.Score
						best.Row = i
						best.Col = j
					}
				}
			}
		}
		return best
	} else {
		best := Move{Score: 1000}
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == 0 {
					board[i][j] = -1
					move := Minimax(board, depth+1, true)
					board[i][j] = 0
					if move.Score < best.Score {
						best.Score = move.Score
						best.Row = i
						best.Col = j
					}
				}
			}
		}
		return best
	}
}

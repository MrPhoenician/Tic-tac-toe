package app

func IsGameOver(board Board) (bool, GameStatus) {
	winner := checkWinner(board)
	if winner != "" {
		if winner == "X" {
			return true, StatusXWon
		}
		return true, StatusOWon
	}

	if isBoardFull(board) {
		return true, StatusDraw
	}

	xCount, oCount := countMoves(board)
	if xCount == oCount {
		return false, StatusPlayerX // X ходит первым
	}
	return false, StatusPlayerO
}

func checkWinner(board Board) string {
	for i := 0; i < 3; i++ {
		if board[i][0] != 0 && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return playerFromInt(board[i][0])
		}
	}

	for j := 0; j < 3; j++ {
		if board[0][j] != 0 && board[0][j] == board[1][j] && board[1][j] == board[2][j] {
			return playerFromInt(board[0][j])
		}
	}

	if board[0][0] != 0 && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return playerFromInt(board[0][0])
	}
	if board[0][2] != 0 && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return playerFromInt(board[0][2])
	}

	return ""
}

func playerFromInt(val int) string {
	if val == 1 {
		return "X"
	}
	if val == -1 {
		return "O"
	}
	return ""
}

func isBoardFull(board Board) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func countMoves(board Board) (xCount, oCount int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == 1 {
				xCount++
			} else if board[i][j] == -1 {
				oCount++
			}
		}
	}
	return xCount, oCount
}

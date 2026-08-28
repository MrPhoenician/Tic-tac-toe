package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tic-tac/api/model"
	"tic-tac/internal/app"
	"tic-tac/internal/contextkeys"
	"tic-tac/internal/logger"
)

type GameHandler struct {
	service app.GameService
	log     logger.Logger
}

func NewGameHandler(service app.GameService, log logger.Logger) *GameHandler {
	return &GameHandler{service: service, log: log}
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextkeys.UserID).(string)

	var req struct {
		VsComputer bool `json:"vs_computer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid JSON", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	game, err := h.service.CreateGame(r.Context(), userID, req.VsComputer)
	if err != nil {
		h.log.Error("Create game failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest) // ← теперь err.Error() работает
		return
	}

	resp := ToWeb(game)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Encode failed", "error", err)
	}
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextkeys.UserID).(string)
	gameID := r.PathValue("id")

	game, err := h.service.JoinGame(r.Context(), gameID, userID)
	if err != nil {
		h.log.Error("Join game failed", "error", err, "game_id", gameID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ToWeb(game)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Encode failed", "error", err)
	}
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")

	game, err := h.service.GetGame(r.Context(), gameID)
	if err != nil {
		h.log.Error("Get game failed", "error", err, "game_id", gameID)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := ToWeb(game)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Encode failed", "error", err)
	}
}

func (h *GameHandler) ListGames(w http.ResponseWriter, r *http.Request) {
	games, err := h.service.ListAvailableGames(r.Context())
	if err != nil {
		h.log.Error("List games failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp []model.GameResponse
	for _, game := range games {
		resp = append(resp, ToWeb(game)) // ← теперь ToWeb принимает *app.Game
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Encode failed", "error", err)
	}
}

func (h *GameHandler) HandleMove(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextkeys.UserID).(string)
	gameID := r.PathValue("id")

	var req model.GameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid JSON", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.ID != gameID {
		http.Error(w, "ID mismatch", http.StatusBadRequest)
		return
	}

	currentGame, err := h.service.GetGame(r.Context(), gameID)
	if err != nil {
		h.log.Error("Game not found", "error", err, "game_id", gameID)
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	isParticipant := false
	if currentGame.PlayerX != nil && *currentGame.PlayerX == userID {
		isParticipant = true
	}
	if currentGame.PlayerO != nil && *currentGame.PlayerO == userID {
		isParticipant = true
	}

	if !isParticipant {
		h.log.Warn("User not in game", "user_id", userID, "game_id", gameID)
		http.Error(w, "you are not a participant in this game", http.StatusForbidden)
		return
	}

	canMove := false
	if currentGame.IsComputer {
		// В игре с компьютером — игрок всегда X
		if currentGame.PlayerX != nil && *currentGame.PlayerX == userID {
			canMove = (currentGame.Status == app.StatusPlayerX)
		}
	} else {
		if currentGame.PlayerX != nil && *currentGame.PlayerX == userID {
			canMove = (currentGame.Status == app.StatusPlayerX)
		}
		if currentGame.PlayerO != nil && *currentGame.PlayerO == userID {
			canMove = (currentGame.Status == app.StatusPlayerO)
		}
	}

	if !canMove {
		h.log.Warn("Not your turn", "user_id", userID, "game_id", gameID, "status", currentGame.Status)
		http.Error(w, "it's not your turn", http.StatusForbidden)
		return
	}

	game, err := h.service.ProcessPlayerMove(r.Context(), gameID, req.Board)
	if err != nil {
		h.log.Error("Move failed", "error", err, "game_id", gameID)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ToWeb(game)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Encode failed", "error", err)
	}
}

func (h *GameHandler) ListCompletedGames(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextkeys.UserID).(string)

	games, err := h.service.GetCompletedGamesByUser(r.Context(), userID)
	if err != nil {
		h.log.Error("List completed games failed", "error", err, "user_id", userID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]model.GameResponse, 0, len(games))
	for _, game := range games {
		resp = append(resp, ToWeb(game))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Encode failed", "error", err)
	}
}

func (h *GameHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("n")
	if limitParam == "" {
		limitParam = "10"
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		http.Error(w, "n must be a positive integer", http.StatusBadRequest)
		return
	}

	entries, err := h.service.GetTopPlayers(r.Context(), limit)
	if err != nil {
		h.log.Error("Get leaderboard failed", "error", err, "limit", limit)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]model.LeaderboardResponse, 0, len(entries))
	for _, entry := range entries {
		resp = append(resp, model.LeaderboardResponse{
			UserID:   entry.UserID,
			Login:    entry.Login,
			WinRatio: entry.WinRatio,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Encode failed", "error", err)
	}
}

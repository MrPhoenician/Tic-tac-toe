package api

import (
	"encoding/json"
	"net/http"
	"tic-tac/api/model"
	"tic-tac/internal/app"
	"tic-tac/internal/contextkeys"
	"tic-tac/internal/logger"
)

type AuthHandler struct {
	authService app.AuthService
	log         logger.Logger
}

func NewAuthHandler(authService app.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req model.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid JSON in signup", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.authService.SignUp(r.Context(), req.Login, req.Password)
	if err != nil {
		h.log.Warn("Signup failed", "error", err, "login", req.Login)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := model.SignUpResponse{
		ID:    user.ID,
		Login: user.Login,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Failed to encode signup response", "error", err)
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req model.JwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid JSON in signin", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.SignIn(r.Context(), req.Login, req.Password)
	if err != nil {
		h.log.Warn("Signin failed", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toJWTResponse(tokens)); err != nil {
		h.log.Error("Failed to encode signin response", "error", err)
	}
}

func (h *AuthHandler) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid JSON in refresh access", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.Warn("Refresh access token failed", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toJWTResponse(tokens)); err != nil {
		h.log.Error("Failed to encode refresh access response", "error", err)
	}
}

func (h *AuthHandler) RefreshRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Invalid JSON in refresh refresh", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.RefreshRefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.Warn("Refresh refresh token failed", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toJWTResponse(tokens)); err != nil {
		h.log.Error("Failed to encode refresh refresh response", "error", err)
	}
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(contextkeys.UserID).(string)
	login, _ := r.Context().Value(contextkeys.UserLogin).(string)

	resp := model.UserResponse{
		ID:    userID,
		Login: login,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("Failed to encode me response", "error", err)
	}
}

func toJWTResponse(tokens *app.AuthTokens) model.JwtResponse {
	if tokens == nil {
		return model.JwtResponse{}
	}

	return model.JwtResponse{
		Type:         tokens.Type,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

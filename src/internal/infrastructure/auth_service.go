package infrastructure

import (
	"context"
	"errors"
	"tic-tac/internal/app"
	"tic-tac/internal/infrastructure/model"
)

type AuthService struct {
	userRepo    *UserRepository
	jwtProvider *JwtProvider
}

var _ app.AuthService = (*AuthService)(nil)

func NewAuthService(userRepo *UserRepository, jwtProvider *JwtProvider) app.AuthService {
	return &AuthService{userRepo: userRepo, jwtProvider: jwtProvider}
}

func (s *AuthService) SignUp(ctx context.Context, login, password string) (*app.User, error) {
	if len(login) < 3 {
		return nil, errors.New("login must be at least 3 characters")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	user, err := s.userRepo.Create(ctx, login, password)
	if err != nil {
		return nil, err
	}

	return toAppUser(user), nil
}

func (s *AuthService) SignIn(ctx context.Context, login, password string) (*app.AuthTokens, error) {
	if login == "" || password == "" {
		return nil, errors.New("login and password are required")
	}

	user, err := s.userRepo.FindByLogin(ctx, login)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !s.userRepo.ValidatePassword(user, password) {
		return nil, errors.New("invalid credentials")
	}

	return s.buildAuthTokens(user)
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*app.AuthTokens, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	if err := s.jwtProvider.ValidateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	user, err := s.getModelUserByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &app.AuthTokens{
		Type:        "Bearer",
		AccessToken: accessToken,
	}, nil
}

func (s *AuthService) RefreshRefreshToken(ctx context.Context, refreshToken string) (*app.AuthTokens, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	if err := s.jwtProvider.ValidateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	user, err := s.getModelUserByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return s.buildAuthTokens(user)
}

func (s *AuthService) GetUserByAccessToken(ctx context.Context, accessToken string) (*app.User, error) {
	if accessToken == "" {
		return nil, errors.New("access token is required")
	}

	if err := s.jwtProvider.ValidateAccessToken(accessToken); err != nil {
		return nil, err
	}

	user, err := s.getModelUserByToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return toAppUser(user), nil
}

func (s *AuthService) buildAuthTokens(user *model.User) (*app.AuthTokens, error) {
	accessToken, err := s.jwtProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtProvider.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &app.AuthTokens{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) getModelUserByToken(ctx context.Context, token string) (*model.User, error) {
	userID, err := s.jwtProvider.GetUUIDFromToken(token)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(ctx, userID)
}

func toAppUser(user *model.User) *app.User {
	if user == nil {
		return nil
	}

	return &app.User{
		ID:    user.ID,
		Login: user.Login,
	}
}

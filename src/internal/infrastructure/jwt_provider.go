package infrastructure

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"tic-tac/internal/config"
	"tic-tac/internal/infrastructure/model"
	"time"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type JwtProvider struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type jwtClaims struct {
	Subject string `json:"sub"`
	Type    string `json:"type"`
	Expires int64  `json:"exp"`
	Issued  int64  `json:"iat"`
}

func NewJwtProvider() *JwtProvider {
	return &JwtProvider{
		secret:          []byte(config.GetJWTSecret()),
		accessTokenTTL:  config.GetAccessTokenTTL(),
		refreshTokenTTL: config.GetRefreshTokenTTL(),
	}
}

func (p *JwtProvider) GenerateAccessToken(user *model.User) (string, error) {
	return p.generateToken(user, tokenTypeAccess, p.accessTokenTTL)
}

func (p *JwtProvider) GenerateRefreshToken(user *model.User) (string, error) {
	return p.generateToken(user, tokenTypeRefresh, p.refreshTokenTTL)
}

func (p *JwtProvider) ValidateAccessToken(token string) error {
	_, err := p.validateToken(token, tokenTypeAccess)
	return err
}

func (p *JwtProvider) ValidateRefreshToken(token string) error {
	_, err := p.validateToken(token, tokenTypeRefresh)
	return err
}

func (p *JwtProvider) GetUUIDFromToken(token string) (string, error) {
	claims, err := p.parseToken(token)
	if err != nil {
		return "", err
	}

	if claims.Subject == "" {
		return "", errors.New("token subject is empty")
	}

	return claims.Subject, nil
}

func (p *JwtProvider) generateToken(user *model.User, tokenType string, ttl time.Duration) (string, error) {
	if user == nil || user.ID == "" {
		return "", errors.New("user is required")
	}

	header, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	claims, err := json.Marshal(jwtClaims{
		Subject: user.ID,
		Type:    tokenType,
		Expires: time.Now().Add(ttl).Unix(),
		Issued:  now,
	})
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(header)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claims)
	payload := encodedHeader + "." + encodedClaims
	signature := p.sign(payload)

	return payload + "." + signature, nil
}

func (p *JwtProvider) validateToken(token, expectedType string) (*jwtClaims, error) {
	claims, err := p.parseToken(token)
	if err != nil {
		return nil, err
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s", expectedType)
	}

	if time.Now().Unix() >= claims.Expires {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func (p *JwtProvider) parseToken(token string) (*jwtClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	payload := parts[0] + "." + parts[1]
	expectedSignature := p.sign(payload)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, errors.New("invalid token signature")
	}

	decodedClaims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims jwtClaims
	if err := json.Unmarshal(decodedClaims, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	return &claims, nil
}

func (p *JwtProvider) sign(payload string) string {
	h := hmac.New(sha256.New, p.secret)
	_, _ = h.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

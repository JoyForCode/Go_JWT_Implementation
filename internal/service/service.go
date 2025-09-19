package service

import (
	"errors"
	"jwt_clean/internal/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	jwtSecret []byte
	users     map[string]string //temporary in memory database
}

func NewAuthService(jwtSecret []byte) *Service {
	return &Service{jwtSecret: jwtSecret,
		users: map[string]string{
			"admin": "password123",
			"user1": "password@123",
		},
	}
}

func (s *Service) newClaims(username, issuer string, duration time.Duration, tokenType auth.TokenType) auth.Claims {
	return auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		TokenType: tokenType,
		UserID:    username,
	}
}

func (s *Service) GenerateToken(username string) (string, error) { //Test function for checks
	claims := s.newClaims(username, "myApp", time.Hour, auth.AccessTokenType)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) GenerateAccessToken(username string) (string, error) {
	claims := s.newClaims(username, "myApp", 15*time.Minute, auth.AccessTokenType)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) GenerateRefreshToken(username string) (string, error) {
	claims := s.newClaims(username, "myApp", 7*24*time.Hour, auth.RefreshTokenType)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) GenerateTokenPair(username string) (*auth.TokenPair, error) {
	accessToken, err := s.GenerateAccessToken(username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.GenerateRefreshToken(username)
	if err != nil {
		return nil, err
	}

	return &auth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		TokenType:    "Bearer",
	}, nil
}

func (s *Service) Login(username, password string) (*auth.TokenPair, error) {
	storedPassword, exists := s.users[username]
	if !exists || storedPassword != password {
		return nil, errors.New("invalid credentials")
	}

	return s.GenerateTokenPair(username)
}

func (s *Service) ValidateAccessToken(tokenString string) (*auth.Claims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != auth.AccessTokenType {
		return nil, errors.New("invalid token type for access")
	}

	return claims, nil
}

func (s *Service) ValidateRefreshToken(tokenString string) (*auth.Claims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != auth.RefreshTokenType {
		return nil, errors.New("invalid token type for refresh")
	}

	return claims, nil
}

func (s *Service) RefreshAccessToken(refreshToken string) (*auth.TokenPair, error) {
	claims, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	accessToken, err := s.GenerateAccessToken(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &auth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
		TokenType:    "Bearer",
	}, nil
}

func (s *Service) ParseToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*auth.Claims); ok {
		return claims, nil
	}

	return nil, errors.New("could not parse claims")
}

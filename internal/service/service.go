package service

import (
	"errors"
	"jwt_clean/internal/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	jwtSecret []byte
}

func NewAuthService(jwtSecret []byte) *Service {
	return &Service{jwtSecret: jwtSecret}
}

func (s *Service) newClaims(username, issuer string, duration time.Duration) auth.Claims {
	return auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
}

func (s *Service) GenerateToken(username string) (string, error) {
	claims := s.newClaims(username, "myApp", time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
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

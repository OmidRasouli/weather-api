package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) secret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}
	return []byte(secret), nil
}

func (s *AuthService) GenerateToken(username string, ttl time.Duration) (string, time.Time, error) {
	sec, err := s.secret()
	if err != nil {
		return "", time.Time{}, err
	}

	exp := time.Now().Add(ttl)
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(sec)
	return signed, exp, err
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	sec, err := s.secret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return sec, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) ValidateCredentials(username, password string) bool {
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	return username == adminUsername && password == adminPassword
}

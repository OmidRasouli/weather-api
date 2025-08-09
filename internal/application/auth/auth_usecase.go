package auth

import (
	"errors"
	"time"

	"github.com/OmidRasouli/weather-api/internal/domain/services"
)

type UseCase struct {
	authService *services.AuthService
}

func NewUseCase(authService *services.AuthService) *UseCase {
	return &UseCase{
		authService: authService,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	ExpiresAt string `json:"expiresAt"`
}

func (uc *UseCase) Login(req LoginRequest) (*LoginResponse, error) {
	if !uc.authService.ValidateCredentials(req.Username, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, exp, err := uc.authService.GenerateToken(req.Username, 24*time.Hour)
	if err != nil {
		return nil, errors.New("could not issue token")
	}

	return &LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresAt: exp.UTC().Format(time.RFC3339),
	}, nil
}

func (uc *UseCase) ValidateToken(tokenString string) (string, error) {
	claims, err := uc.authService.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}

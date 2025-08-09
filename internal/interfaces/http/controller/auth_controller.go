package controller

import (
	"net/http"

	authUseCase "github.com/OmidRasouli/weather-api/internal/application/auth"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUseCase *authUseCase.UseCase
}

func NewAuthController(authUseCase *authUseCase.UseCase) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	var req authUseCase.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	response, err := ac.authUseCase.Login(req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
		return
	}

	c.JSON(http.StatusOK, response)
}

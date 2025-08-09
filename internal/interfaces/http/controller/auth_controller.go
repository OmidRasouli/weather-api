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

// DTO used only for Swagger documentation to avoid parser issues with external package types.
type LoginRequestDTO struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"secret"`
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return a JWT access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequestDTO true "Login credentials"
// @Success      200  {object}  map[string]string "token response"
// @Failure      400  {object}  map[string]string "validation error"
// @Failure      401  {object}  map[string]string "invalid credentials"
// @Failure      500  {object}  map[string]string "server error"
// @Router       /login [post]
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

package auth

// Deprecated: This package is deprecated. Use the following instead:
// - internal/domain/services/auth_service.go for JWT operations
// - internal/application/auth/auth_usecase.go for business logic
// - internal/interfaces/http/middleware/auth_middleware.go for HTTP middleware
//
// This file remains for backward compatibility only.

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func secret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, errors.New("JWT_SECRET not set")
	}
	return []byte(s), nil
}

func GenerateToken(username string, ttl time.Duration) (string, time.Time, error) {
	sec, err := secret()
	if err != nil {
		return "", time.Time{}, err
	}
	exp := time.Now().Add(ttl)
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(sec)
	return signed, exp, err
}

func parseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	sec, err := secret()
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

// Deprecated: Use domain/services/AuthService.ValidateToken instead
func ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	return parseToken(tokenString)
}

// Deprecated: Use middleware/JWTAuth instead
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			return
		}
		claims, err := parseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		// Attach subject (username) to context
		if claims.Subject != "" {
			c.Set("user", claims.Subject)
		}
		c.Next()
	}
}

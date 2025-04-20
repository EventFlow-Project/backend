package services

import (
	"errors"
	"time"

	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/models"

	"github.com/golang-jwt/jwt"
)

type JWTService struct {
	config *config.Config
}

func NewJWTService(config *config.Config) *JWTService {
	return &JWTService{
		config: config,
	}
}

func (s *JWTService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})
}

func (s *JWTService) GetUserIDFromToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	return userID, nil
}

package services

import (
	"errors"

	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       ports.AuthRepository
	jwtService *JWTService
}

func NewAuthService(repo ports.AuthRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(credentials models.RegistrationCredentials) (*models.User, error) {
	if _, err := s.repo.GetUserByEmail(credentials.Email); err == nil {
		return nil, errors.New("user already exists")
	}

	user, err := s.repo.CreateUserWithPassword(credentials.Email, credentials.Password, credentials.Name, credentials.Role, credentials.Description, credentials.ActivityArea)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	return s.jwtService.GenerateToken(user)
}

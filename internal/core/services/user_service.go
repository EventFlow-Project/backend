package services

import (
	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
)

type UserService struct {
	repo   ports.UserRepository
	config *config.Config
}

func NewUserService(repo ports.UserRepository, config *config.Config) *UserService {
	return &UserService{
		repo:   repo,
		config: config,
	}
}

func (s *UserService) GetUserInfo(accessToken string) (*models.SafeUser, error) {
	user, err := s.repo.GetUserInfo(accessToken)
	if err != nil {
		return nil, err
	}
	return user.ToSafeUser(), nil
}

func (s *UserService) EditUserInfo(accessToken string, info *models.EditUserInfo) (*models.SafeUser, error) {
	user, err := s.repo.EditUserInfo(accessToken, info)
	if err != nil {
		return nil, err
	}
	return user.ToSafeUser(), nil
}

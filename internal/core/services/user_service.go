package services

import (
	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
)

type UserService struct {
	repo       ports.UserRepository
	config     *config.Config
	jwtService *JWTService
}

func NewUserService(repo ports.UserRepository, config *config.Config, jwtService *JWTService) *UserService {
	return &UserService{
		repo:       repo,
		config:     config,
		jwtService: jwtService,
	}
}

func (s *UserService) GetUserInfo(accessToken string) (*models.SafeUser, error) {
	userID, err := s.jwtService.GetUserIDFromToken(accessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user.ToSafeUser(), nil
}

func (s *UserService) EditUserInfo(accessToken string, info *models.EditUserInfo) (*models.SafeUser, error) {
	userID, err := s.jwtService.GetUserIDFromToken(accessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.EditUserInfo(userID, info)
	if err != nil {
		return nil, err
	}

	return user.ToSafeUser(), nil
}

func (s *UserService) SearchUsersByName(name string) ([]*models.SearchUserResponse, error) {
	users, err := s.repo.SearchUsersByName(name)
	if err != nil {
		return nil, err
	}

	searchResults := make([]*models.SearchUserResponse, len(users))
	for i, user := range users {
		searchResults[i] = user.ToSearchResponse()
	}

	return searchResults, nil
}

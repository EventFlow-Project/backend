package services

import (
	"errors"

	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
)

type FriendService struct {
	repo ports.FriendRepository
}

func NewFriendService(repo ports.FriendRepository) *FriendService {
	return &FriendService{
		repo: repo,
	}
}

func (s *FriendService) SendFriendRequest(fromID, toID string) (*models.FriendRequestResponse, error) {
	exists, err := s.repo.CheckExistingRequest(fromID, toID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("friend request already exists")
	}

	return s.repo.CreateFriendRequest(fromID, toID)
}

func (s *FriendService) RespondToFriendRequest(friendID string, accept bool) error {
	request, err := s.repo.GetFriendRequestByFromID(friendID)
	if err != nil {
		return err
	}

	status := "rejected"
	if accept {
		status = "accepted"
	}

	return s.repo.UpdateFriendRequestStatus(request.ID, status)
}

func (s *FriendService) GetFriendsList(userID string) (*models.FriendListResponse, error) {
	friends, err := s.repo.GetFriendsList(userID)
	if err != nil {
		return nil, err
	}

	return &models.FriendListResponse{
		Friends: friends,
	}, nil
}

func (s *FriendService) RemoveFriend(userID, friendID string) error {
	return s.repo.RemoveFriend(userID, friendID)
}

func (s *FriendService) GetIncomingFriendRequests(userID string) ([]models.FriendRequestResponse, error) {
	return s.repo.GetIncomingFriendRequests(userID)
}

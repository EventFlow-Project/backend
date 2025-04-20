package ports

import "github.com/EventFlow-Project/backend/internal/core/models"

type FriendRepository interface {
	CreateFriendRequest(fromID, toID string) (*models.FriendRequestResponse, error)
	UpdateFriendRequestStatus(requestID string, status string) error
	GetFriendsList(userID string) ([]models.SafeUser, error)
	RemoveFriend(userID, friendID string) error
	GetFriendRequest(requestID string) (*models.FriendRequest, error)
	GetFriendRequestByFromID(fromID string) (*models.FriendRequest, error)
	GetIncomingFriendRequests(userID string) ([]models.FriendRequestResponse, error)
	CheckExistingRequest(fromID, toID string) (bool, error)
	CreateFriendship(userID, friendID string) error
	DeleteFriendship(userID, friendID string) error
}

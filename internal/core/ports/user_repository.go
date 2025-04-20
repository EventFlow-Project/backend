package ports

import "github.com/EventFlow-Project/backend/internal/core/models"

type UserRepository interface {
	GetUserByID(userID string) (*models.User, error)
	EditUserInfo(userID string, info *models.EditUserInfo) (*models.User, error)
	SearchUsersByName(name string) ([]*models.User, error)
}

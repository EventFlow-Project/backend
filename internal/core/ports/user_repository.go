package ports

import "github.com/EventFlow-Project/backend/internal/core/models"

type UserRepository interface {
	GetUserInfo(accessToken string) (*models.User, error)
	EditUserInfo(accessToken string, info *models.EditUserInfo) (*models.User, error)
}

package ports

import "github.com/EventFlow-Project/backend/internal/core/models"

type AuthRepository interface {
	CreateUserWithPassword(email, password, name, role, description, activityArea string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
	UpdateUser(user *models.User) error
}

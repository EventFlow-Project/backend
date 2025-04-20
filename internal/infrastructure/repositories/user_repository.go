package repositories

import (
	"errors"
	"time"

	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/EventFlow-Project/backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) ports.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) GetUserByID(userID string) (*models.User, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}
	var user models.User

	result := r.db.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepositoryImpl) EditUserInfo(userID string, info *models.EditUserInfo) (*models.User, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}
	var user models.User

	result := r.db.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, result.Error
	}

	if info.Email != user.Email {
		var existingUser models.User
		result := r.db.DB.Where("email = ? AND id != ?", info.Email, userID).First(&existingUser)

		if result.Error == nil {
			return nil, errors.New("email already in use")
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
	}
	user.Email = info.Email
	user.Name = info.Name
	user.Avatar = info.Avatar
	user.UpdatedAt = time.Now()

	result = r.db.DB.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepositoryImpl) SearchUsersByName(name string) ([]*models.User, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}
	var users []*models.User

	result := r.db.DB.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

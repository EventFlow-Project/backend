package repositories

import (
	"errors"
	"time"

	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/EventFlow-Project/backend/internal/core/services"
	"github.com/EventFlow-Project/backend/internal/infrastructure/database"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db         *database.Database
	jwtService *services.JWTService
}

func NewUserRepository(db *database.Database, jwtService *services.JWTService) ports.UserRepository {
	return &UserRepositoryImpl{
		db:         db,
		jwtService: jwtService,
	}
}

func (r *UserRepositoryImpl) GetUserInfo(accessToken string) (*models.User, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	token, err := r.jwtService.ValidateToken(accessToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user ID not found in token")
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

func (r *UserRepositoryImpl) EditUserInfo(accessToken string, info *models.EditUserInfo) (*models.User, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	token, err := r.jwtService.ValidateToken(accessToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("user ID not found in token")
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

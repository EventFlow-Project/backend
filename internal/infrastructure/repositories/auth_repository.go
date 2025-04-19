package repositories

import (
	"errors"
	"time"

	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/EventFlow-Project/backend/internal/infrastructure/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
	db *database.Database
}

func NewAuthRepository(db *database.Database) ports.AuthRepository {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (r *AuthRepositoryImpl) CreateUserWithPassword(email, password, name, role, description, activityArea string) (*models.User, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         name,
		PasswordHash: string(hashedPassword),
		Avatar:       "",
		Role:         role,
		Description:  description,
		ActivityArea: activityArea,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result := r.db.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (r *AuthRepositoryImpl) UpdateUser(user *models.User) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}

	user.UpdatedAt = time.Now()
	result := r.db.DB.Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *AuthRepositoryImpl) GetUserByName(name string) (*models.User, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var user models.User
	result := r.db.DB.Where("name = ?", name).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var user models.User
	result := r.db.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

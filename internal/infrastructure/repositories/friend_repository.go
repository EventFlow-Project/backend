package repositories

import (
	"errors"
	"time"

	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/EventFlow-Project/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendRepositoryImpl struct {
	db *database.Database
}

func NewFriendRepository(db *database.Database) ports.FriendRepository {
	return &FriendRepositoryImpl{
		db: db,
	}
}

func (r *FriendRepositoryImpl) CreateFriendRequest(fromID, toID string) (*models.FriendRequestResponse, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	request := &models.FriendRequest{
		ID:        uuid.New().String(),
		FromID:    fromID,
		ToID:      toID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.DB.Create(request).Error; err != nil {
		return nil, err
	}

	return &models.FriendRequestResponse{
		ID:        request.ID,
		FromID:    request.FromID,
		ToID:      request.ToID,
		Status:    request.Status,
		CreatedAt: request.CreatedAt,
	}, nil
}

func (r *FriendRepositoryImpl) UpdateFriendRequestStatus(requestID string, status string) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}
	var request models.FriendRequest

	if err := r.db.DB.First(&request, "id = ?", requestID).Error; err != nil {
		return err
	}

	tx := r.db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Model(&request).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if status == "accepted" {
		r.CreateFriendship(request.FromID, request.ToID)
	}

	return tx.Commit().Error
}

func (r *FriendRepositoryImpl) GetFriendsList(userID string) ([]models.SafeUser, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}
	var users []*models.User

	err := r.db.DB.Model(&models.User{}).
		Distinct("users.*").
		Joins("JOIN friendships ON ((friendships.user_id = ? AND friendships.friend_id = users.id) OR (friendships.friend_id = ? AND friendships.user_id = users.id))", userID, userID).
		Where("friendships.deleted_at IS NULL").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	safeUsers := make([]models.SafeUser, len(users))
	for i, user := range users {
		safeUsers[i] = models.SafeUser{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		}
	}

	return safeUsers, nil
}

func (r *FriendRepositoryImpl) RemoveFriend(userID, friendID string) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}

	tx := r.db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	now := time.Now()
	if err := tx.Model(&models.Friendship{}).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			userID, friendID, friendID, userID).
		Update("deleted_at", now).Error; err != nil {
		tx.Rollback()

		return err
	}

	if err := tx.Model(&models.FriendRequest{}).
		Where("(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)",
			userID, friendID, friendID, userID).
		Update("deleted_at", now).Error; err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}

func (r *FriendRepositoryImpl) GetFriendRequest(requestID string) (*models.FriendRequest, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}
	var request models.FriendRequest

	if err := r.db.DB.First(&request, "id = ?", requestID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("friend request not found")
		}
		return nil, err
	}

	return &request, nil
}

func (r *FriendRepositoryImpl) CheckExistingRequest(fromID, toID string) (bool, error) {
	if r.db == nil || r.db.DB == nil {
		return false, errors.New("database connection is not initialized")
	}
	var count int64

	err := r.db.DB.Model(&models.FriendRequest{}).
		Where("from_id = ? AND to_id = ? AND status = ?",
			fromID, toID, "pending").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *FriendRepositoryImpl) CreateFriendship(userID, friendID string) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}

	friendship1 := &models.Friendship{
		ID:        uuid.New().String(),
		UserID:    userID,
		FriendID:  friendID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	friendship2 := &models.Friendship{
		ID:        uuid.New().String(),
		UserID:    friendID,
		FriendID:  userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return r.db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(friendship1).Error; err != nil {
			return err
		}

		if err := tx.Create(friendship2).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *FriendRepositoryImpl) DeleteFriendship(userID, friendID string) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}

	return r.db.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friendID, friendID, userID).
		Delete(&models.Friendship{}).Error
}

func (r *FriendRepositoryImpl) GetFriendRequestByFromID(fromID string) (*models.FriendRequest, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var request models.FriendRequest
	if err := r.db.DB.First(&request, "from_id = ? AND status = ?", fromID, "pending").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("friend request not found")
		}
		return nil, err
	}

	return &request, nil
}

func (r *FriendRepositoryImpl) GetIncomingFriendRequests(userID string) ([]models.FriendRequestResponse, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var requests []models.FriendRequest
	if err := r.db.DB.Where("to_id = ? AND status = ?", userID, "pending").Find(&requests).Error; err != nil {
		return nil, err
	}

	responses := make([]models.FriendRequestResponse, len(requests))
	for i, request := range requests {
		var user models.User
		if err := r.db.DB.First(&user, "id = ?", request.FromID).Error; err != nil {
			return nil, err
		}

		responses[i] = models.FriendRequestResponse{
			ID:         request.ID,
			FromID:     request.FromID,
			ToID:       request.ToID,
			Status:     request.Status,
			CreatedAt:  request.CreatedAt,
			FromName:   user.Name,
			FromAvatar: user.Avatar,
		}
	}

	return responses, nil
}

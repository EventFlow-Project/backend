package services

import (
	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/google/uuid"
)

type MinioService struct {
	repo   ports.MinioRepository
	config *config.Config
}

func NewMinioService(repo ports.MinioRepository, config *config.Config) *MinioService {
	return &MinioService{
		repo:   repo,
		config: config,
	}
}

func (s *MinioService) UploadImage(base64Data string) (string, error) {
	fileName := uuid.New().String() + ".jpg"

	return s.repo.UploadBase64Image(base64Data, fileName)
}

func (s *MinioService) DeleteImage(fileName string) error {
	return s.repo.DeleteImage(fileName)
}

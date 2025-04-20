package repositories

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioRepositoryImpl struct {
	client *minio.Client
	config *config.Config
}

func NewMinioRepository(cfg *config.Config) (ports.MinioRepository, error) {
	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	exists, err := minioClient.BucketExists(context.Background(), cfg.Minio.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), cfg.Minio.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MinioRepositoryImpl{
		client: minioClient,
		config: cfg,
	}, nil
}

func (r *MinioRepositoryImpl) UploadBase64Image(base64Data string, fileName string) (string, error) {
	base64Data = strings.TrimPrefix(base64Data, "data:image/")
	if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
		base64Data = base64Data[idx+8:]
	}

	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	_, err = r.client.PutObject(
		context.Background(),
		r.config.Minio.BucketName,
		fileName,
		bytes.NewReader(imageData),
		int64(len(imageData)),
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to minio: %w", err)
	}

	protocol := "http"
	if r.config.Minio.UseSSL {
		protocol = "https"
	}
	fileURL := fmt.Sprintf("%s://%s/%s/%s", protocol, r.config.Minio.PublicEndpoint, r.config.Minio.BucketName, fileName)

	return fileURL, nil
}

func (r *MinioRepositoryImpl) DeleteImage(fileName string) error {
	err := r.client.RemoveObject(
		context.Background(),
		r.config.Minio.BucketName,
		fileName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to delete file from minio: %w", err)

	}

	return nil
}

package services

import (
	"backend/internal/configs"
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"sync"
)

type MinioService interface {
	UploadToMinio(ctx context.Context, bucketName, objectName string, file io.Reader, size int64, contentType string) error
	DeleteFromMinio(ctx context.Context, bucketName, objectName string) error
	BatchDeleteFromMinio(ctx context.Context, bucketName string, objectNames []string) error
}

type minioService struct {
	minioClient *minio.Client
}

var (
	globalMinioService MinioService
	minioOnce          sync.Once
)

// Minio桶名
const (
	ImagesBucket = "images"
	FilesBucket  = "files"
)

func initMinioService(minioClient *minio.Client) {
	globalMinioService = &minioService{minioClient: minioClient}
}

func GetMinioService() MinioService {
	if globalMinioService == nil {
		minioOnce.Do(func() {
			initMinioService(configs.MinioClient)
		})
	}
	return globalMinioService
}

// UploadToMinio 上传文件到Minio
func (s *minioService) UploadToMinio(ctx context.Context, bucketName, objectName string, file io.Reader, size int64, contentType string) error {
	_, err := s.minioClient.PutObject(
		ctx,
		bucketName,
		objectName,
		file,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func (s *minioService) DeleteFromMinio(ctx context.Context, bucketName, objectName string) error {
	return s.minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (s *minioService) BatchDeleteFromMinio(ctx context.Context, bucketName string, objectNames []string) error {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, objectName := range objectNames {
			objectsCh <- minio.ObjectInfo{
				Key: objectName,
			}
		}
	}()

	var errs []error
	for err := range s.minioClient.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			errs = append(errs, err.Err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

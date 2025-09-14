package services

import (
	"backend/internal/configs"
	"backend/internal/models"
	"backend/internal/models/dto"
	"backend/internal/repositories"
	"backend/pkg/fileUtil"
	"backend/pkg/logger"
	"backend/pkg/minioContent"
	"backend/pkg/redisContent"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type UploadService interface {
	UploadImg(email string, fileHeader *multipart.FileHeader) (string, error)
	UploadToMinio(bucketName, objectName string, file io.Reader, size int64, contentType string) error
	UploadFile(id int, files []*multipart.FileHeader) error
	DeleteFile(ids []int) error
	DeleteFromMinio(bucketName, objectName string) error
	BatchDeleteFromMinio(bucketName string, objectNames []string) error
}

type uploadService struct {
	authRepository  repository.AuthRepository
	imageRepository repository.ImageRepository
	fileRepository  repository.FileRepository
}

func NewUploadService(authRepository repository.AuthRepository,
	imageRepository repository.ImageRepository,
	fileRepository repository.FileRepository) UploadService {
	return &uploadService{
		authRepository:  authRepository,
		imageRepository: imageRepository,
		fileRepository:  fileRepository,
	}
}

// UploadImg 上传图片服务
func (s *uploadService) UploadImg(email string, fileHeader *multipart.FileHeader) (string, error) {
	// 读取文件
	file, err := fileHeader.Open()
	if err != nil {
		logger.Error("读取文件失败", logger.Err(err))
		return "", err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			logger.Error("关闭文件失败", logger.Err(err))
		}
	}(file)

	// 计算文件的md5值
	md5Hash, err := fileUtil.GetFileMD5(file)
	if err != nil {
		logger.Error("计算文件MD5值失败", logger.Err(err))
		return "", err
	}

	// 根据MD5值查询数据库中是否已经有该图片
	image, err := s.imageRepository.GetImageByMD5(md5Hash)
	if err != nil {
		logger.Error("根据md5值查询图片失败", logger.String("md5", md5Hash), logger.Err(err))
		return "", err
	}

	if image != (models.Image{}) {
		// 数据库中已有该图片，更新用户URL，然后直接返回图片的URL
		user, err := s.authRepository.SelectUserByEmail(email)
		if err != nil {
			logger.Error("获取用户信息失败", logger.String("email", email), logger.Err(err))
			return "", err
		}
		user.Pic = image.Url
		err = s.authRepository.UpdateUser(user, nil)
		if err != nil {
			logger.Error("更新用户信息失败", logger.String("email", email), logger.Err(err))
			return "", err
		}

		return image.Url, nil
	}

	// 数据库中没有该文件，上传到minio即数据库中
	_, err = file.Seek(0, 0) // 重置文件指针位置
	if err != nil {
		return "", err
	}

	// 上传到minio
	extensionName := fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):]
	// 按照年/月/日的方式获取当前日期
	date := time.Now().Format("2006/01/02")
	objectName := date + "/" + md5Hash + extensionName

	contentType := fileHeader.Header.Get("Content-Type")
	err = s.UploadToMinio(minioContent.ImagesBucket, objectName, file, fileHeader.Size, contentType)
	if err != nil {
		return "", err
	}

	// 写入数据库，写入图片表，并更新用户表
	// 开启事务
	tx := configs.MysqlDb.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}

	image = models.Image{
		Md5: md5Hash,
		Url: "/" + minioContent.ImagesBucket + "/" + objectName,
	}
	err = s.imageRepository.InsertImage(image, tx)
	if err != nil {
		tx.Rollback() // 回滚事务
		return "", err
	}

	// 更新用户表
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		tx.Rollback() // 回滚事务
		return "", err
	}

	// 更新前删除缓存
	redisClient := configs.RedisClient
	emailKey := redisContent.UserInfoKey + email
	idKey := fmt.Sprintf(redisContent.UserInfoKey+"%d", user.Id)
	err = redisClient.Del(context.Background(), emailKey, idKey).Err()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	user.Pic = image.Url
	err = s.authRepository.UpdateUser(user, tx)
	if err != nil {
		tx.Rollback() // 回滚事务
		return "", err
	}

	// 更新缓存
	userDto := dto.UserDto{
		Id:          user.Id,
		Email:       user.Email,
		Name:        user.Name,
		Pic:         user.Pic,
		Description: user.Description,
	}
	userJson, err := json.Marshal(userDto)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	redisClient.Set(context.Background(), emailKey, userJson, 24*time.Hour)
	redisClient.Set(context.Background(), idKey, userJson, 24*time.Hour)

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		// 提交事务异常
		tx.Rollback()
		logger.Error("提交事务异常", logger.Err(err))
		return "", err
	}

	return image.Url, nil
}

// UploadToMinio 上传到minio
func (s *uploadService) UploadToMinio(bucketName, objectName string, file io.Reader, size int64, contentType string) error {
	// 获取Minio客户端
	minioClient := configs.MinioClient

	ctx := context.Background()
	// 上传到minio
	_, err := minioClient.PutObject(
		ctx,
		bucketName,
		objectName,
		file,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		logger.Error("上传Minio失败",
			logger.String("bucketName", bucketName),
			logger.String("objectName", objectName),
			logger.Err(err))
	}

	return err
}

// UploadFile 上传文件
func (s *uploadService) UploadFile(id int, files []*multipart.FileHeader) error {
	for _, fileHeader := range files {
		// 读文件
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}

		// 上传到minio
		objectName := strconv.Itoa(id) + "/" + fileHeader.Filename
		err = s.UploadToMinio(minioContent.FilesBucket, objectName, file, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return err
		}

		// 插入数据库
		fileDb := models.File{
			TaskId: id,
			Name:   fileHeader.Filename,
			Url:    "/" + minioContent.FilesBucket + "/" + objectName,
		}
		err = s.fileRepository.Insert(fileDb, nil)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *uploadService) DeleteFile(ids []int) error {
	// 删除minio中的附件
	files, err := s.fileRepository.GetFileByIds(ids)
	if err != nil {
		return err
	}

	var objectNames []string
	for _, file := range files {
		objectNames = append(objectNames, strconv.Itoa(file.TaskId)+"/"+file.Name)
	}
	err = s.BatchDeleteFromMinio(minioContent.FilesBucket, objectNames)
	if err != nil {
		return err
	}

	// 删除数据库中的附件
	err = s.fileRepository.DeleteByIds(ids, nil)

	return err
}

// DeleteFromMinio 删除Minio中的附件
func (s *uploadService) DeleteFromMinio(bucketName, objectName string) error {
	minioClient := configs.MinioClient
	ctx := context.Background()
	err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		logger.Error("从minio中删除附件失败",
			logger.String("bucketName", bucketName),
			logger.String("objectName", objectName),
			logger.Err(err))
	}
	return err
}

// BatchDeleteFromMinio 批量删除Minio中的附件
func (s *uploadService) BatchDeleteFromMinio(bucketName string, objectNames []string) error {
	minioClient := configs.MinioClient
	ctx := context.Background()

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
	for err := range minioClient.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		errs = append(errs, err.Err)
	}
	if len(errs) > 0 {
		logger.Warn("从minio中批量删除附件异常",
			logger.String("bucketName", bucketName),
			logger.Err(errs[0]))
		return errs[0]
	}
	return nil
}

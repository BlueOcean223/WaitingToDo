package service

import (
	"back/configs"
	"back/models"
	"back/repository"
	"back/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

type UploadService struct {
	authRepository  *repository.AuthRepository
	imageRepository *repository.ImageRepository
}

func NewUploadService(authRepository *repository.AuthRepository,
	imageRepository *repository.ImageRepository) *UploadService {
	return &UploadService{
		authRepository:  authRepository,
		imageRepository: imageRepository,
	}
}

// UploadImg 上传图片
func (s *UploadService) UploadImg(email string, fileHeader *multipart.FileHeader) (string, error) {
	// 读取文件
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 计算文件的md5值
	md5Hash, err := utils.GetFileMD5(file)
	if err != nil {
		return "", err
	}

	// 根据MD5值查询数据库中是否已经有该图片
	image, err := s.imageRepository.GetImageByMD5(md5Hash)
	if err != nil {
		return "", err
	}

	if image != (models.Image{}) {
		// 数据库中已有该图片，更新用户URL，然后直接返回图片的URL
		user, err := s.authRepository.SelectUserByEmail(email)
		if err != nil {
			return "", err
		}
		user.Pic = image.Url
		err = s.authRepository.UpdateUser(user)
		if err != nil {
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
	err = s.UploadImgToMinio(utils.ImagesBucket, objectName, file, fileHeader.Size, contentType)
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
		Url: "/" + utils.ImagesBucket + "/" + objectName,
	}
	err = s.imageRepository.InsertImage(image)
	if err != nil {
		tx.Rollback() // 回滚事务
		return "", err
	}

	// 更新用户表
	redisClient := configs.RedisClient
	key := utils.UserInfoKey + email
	var user models.User

	if redisClient.Exists(context.Background(), key).Val() == 1 {
		val, err := redisClient.Get(context.Background(), key).Bytes()
		if err != nil {
			tx.Rollback()
			return "", err
		}

		err = json.Unmarshal(val, &user)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	} else {
		user, err = s.authRepository.SelectUserByEmail(email)
		if err != nil {
			tx.Rollback() // 回滚事务
			return "", err
		}

		// 将用户信息写入缓存
		val, err := json.Marshal(user)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		redisClient.Set(context.Background(), key, val, 24*time.Hour)
	}

	// 更新前删除缓存
	emailKey := key
	idKey := fmt.Sprintf(utils.UserInfoKey+"%d", user.Id)
	err = redisClient.Del(context.Background(), emailKey, idKey).Err()
	if err != nil {
		tx.Rollback()
		return "", err
	}

	user.Pic = image.Url
	err = s.authRepository.UpdateUser(user)
	if err != nil {
		tx.Rollback() // 回滚事务
		return "", err
	}

	// 更新缓存
	userJson, err := json.Marshal(user)
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
		return "", err
	}

	return image.Url, nil
}

func (s *UploadService) UploadImgToMinio(bucketName, objectName string, file io.Reader, size int64, contentType string) error {
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

	return err
}

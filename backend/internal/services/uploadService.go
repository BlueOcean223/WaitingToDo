package services

import (
	"backend/internal/configs"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/pkg/fileUtil"
	"backend/pkg/logger"
	"backend/pkg/redisContent"
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type UploadService interface {
	UploadImg(email string, fileHeader *multipart.FileHeader) (string, error)
	UploadFile(id int, files []*multipart.FileHeader) error
	DeleteFile(ids []int) error
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
	err = GetMinioService().UploadToMinio(context.Background(), ImagesBucket, objectName, file, fileHeader.Size, contentType)
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
		Url: "/" + ImagesBucket + "/" + objectName,
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

	user.Pic = image.Url
	err = s.authRepository.UpdateUser(user, tx)
	if err != nil {
		tx.Rollback() // 回滚事务
		return "", err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		// 提交事务异常
		tx.Rollback()
		logger.Error("提交事务异常", logger.Err(err))
		return "", err
	}

	// 删除redis缓存
	emailKey := redisContent.UserInfoKey + email
	idKey := fmt.Sprintf(redisContent.UserInfoKey+"%d", user.Id)
	configs.RedisClient.Del(context.Background(), emailKey, idKey)

	return image.Url, nil
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
		err = GetMinioService().UploadToMinio(context.Background(), FilesBucket, objectName, file, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
		if err != nil {
			return err
		}

		// 插入数据库
		fileDb := models.File{
			TaskId: id,
			Name:   fileHeader.Filename,
			Url:    "/" + FilesBucket + "/" + objectName,
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
	err = GetMinioService().BatchDeleteFromMinio(context.Background(), FilesBucket, objectNames)
	if err != nil {
		return err
	}

	// 删除数据库中的附件
	err = s.fileRepository.DeleteByIds(ids, nil)

	return err
}

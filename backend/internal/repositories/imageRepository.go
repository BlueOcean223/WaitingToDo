package repository

import (
	"backend/internal/models"
	"errors"
	"gorm.io/gorm"
)

type ImageRepository interface {
	GetImageByMD5(md5 string) (models.Image, error)
	InsertImage(image models.Image, tx *gorm.DB) error
}

type imageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{db: db}
}

// GetImageByMD5 根据md5值查询图片
func (s *imageRepository) GetImageByMD5(md5 string) (models.Image, error) {
	var image models.Image
	err := s.db.Where("md5 = ?", md5).First(&image).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 处理未查询到的情况
		return models.Image{}, nil
	}

	return image, err
}

// InsertImage 插入图片信息
func (s *imageRepository) InsertImage(image models.Image, tx *gorm.DB) error {
	db := s.db
	if tx != nil {
		db = tx
	}
	return db.Create(&image).Error
}

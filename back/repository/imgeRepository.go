package repository

import (
	"back/configs"
	"back/models"
	"errors"
	"gorm.io/gorm"
)

// GetImageByMD5 根据md5值查询图片
func GetImageByMD5(md5 string) (models.Image, error) {
	db := configs.MysqlDb
	var image models.Image
	err := db.Where("md5 = ?", md5).First(&image).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 处理未查询到的情况
		return models.Image{}, nil
	}

	return image, err
}

// InsertImage 插入图片信息
func InsertImage(image models.Image) error {
	db := configs.MysqlDb
	err := db.Create(&image).Error
	return err
}

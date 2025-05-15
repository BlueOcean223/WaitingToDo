package repository

import (
	"back/configs"
	"back/models"
	"errors"
	"gorm.io/gorm"
)

// SelectUserByEmail 根据邮箱查询用户
func SelectUserByEmail(email string) (models.User, error) {
	db := configs.MysqlDb
	// 根据邮箱查询用户
	var user models.User
	result := db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 处理未查询到的情况
		return models.User{}, nil
	}

	return user, result.Error
}

// InsertUser 插入用户
func InsertUser(user models.User) error {
	db := configs.MysqlDb
	result := db.Create(&user)
	return result.Error
}

// UpdateUser 更新用户
func UpdateUser(user models.User) error {
	db := configs.MysqlDb
	result := db.Save(&user)
	return result.Error
}

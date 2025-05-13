package repository

import (
	"back/configs"
	"back/models"
)

// SelectUserByEmail 根据邮箱查询用户
func SelectUserByEmail(email string) (models.User, error) {
	db := configs.MysqlDb
	// 根据邮箱查询用户
	var user models.User
	result := db.Where("email = ?", email).First(&user)

	return user, result.Error
}

package repository

import (
	"back/models"
	"errors"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// SelectUserByEmail 根据邮箱查询用户
func (s *AuthRepository) SelectUserByEmail(email string) (models.User, error) {
	// 根据邮箱查询用户
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 处理未查询到的情况
		return models.User{}, nil
	}

	return user, result.Error
}

// SelectUserById 根据Id查询用户
func (s *AuthRepository) SelectUserById(id int) (models.User, error) {
	var user models.User
	result := s.db.Where("id = ?", id).First(&user)

	return user, result.Error
}

// SelectUsersByIds 根据Ids批量查询用户
func (s *AuthRepository) SelectUsersByIds(ids []int) ([]models.User, error) {
	var users []models.User
	result := s.db.Where("id in (?)", ids).Find(&users)

	return users, result.Error
}

// InsertUser 插入用户
func (s *AuthRepository) InsertUser(user models.User, tx *gorm.DB) error {
	db := s.db
	if tx != nil {
		db = tx
	}
	result := db.Create(&user)
	return result.Error
}

// UpdateUser 更新用户
func (s *AuthRepository) UpdateUser(user models.User, tx *gorm.DB) error {
	db := s.db
	if tx != nil {
		db = tx
	}
	return db.Save(&user).Error
}

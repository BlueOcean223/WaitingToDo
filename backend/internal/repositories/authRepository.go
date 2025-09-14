package repository

import (
	"backend/internal/models"
	"errors"
	"gorm.io/gorm"
)

type AuthRepository interface {
	SelectUserByEmail(email string) (models.User, error)
	SelectUserById(id int) (models.User, error)
	SelectUsersByIds(ids []int) ([]models.User, error)
	InsertUser(user models.User, tx *gorm.DB) error
	UpdateUser(user models.User, tx *gorm.DB) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// SelectUserByEmail 根据邮箱查询用户
func (s *authRepository) SelectUserByEmail(email string) (models.User, error) {
	// 根据邮箱查询用户
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 处理未查询到的情况
		return models.User{}, nil
	}

	return user, err
}

// SelectUserById 根据Id查询用户
func (s *authRepository) SelectUserById(id int) (models.User, error) {
	var user models.User
	err := s.db.Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 处理未查询到的情况
		return models.User{}, nil
	}

	return user, err
}

// SelectUsersByIds 根据Ids批量查询用户
func (s *authRepository) SelectUsersByIds(ids []int) ([]models.User, error) {
	var users []models.User
	err := s.db.Where("id in (?)", ids).Find(&users).Error

	return users, err
}

// InsertUser 插入用户
func (s *authRepository) InsertUser(user models.User, tx *gorm.DB) error {
	db := s.db
	if tx != nil {
		db = tx
	}
	return db.Create(&user).Error
}

// UpdateUser 更新用户
func (s *authRepository) UpdateUser(user models.User, tx *gorm.DB) error {
	db := s.db
	if tx != nil {
		db = tx
	}
	return db.Model(&user).Updates(user).Error
}

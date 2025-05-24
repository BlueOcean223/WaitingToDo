package repository

import (
	"back/models"
	"errors"
	"gorm.io/gorm"
)

type FriendRepository struct {
	db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

// GetFriendList 根据用户id查询好友列表
func (s *FriendRepository) GetFriendList(userId int) ([]models.User, error) {
	var friends []models.User
	subQuery := s.db.Table("friends").Select("friend_id").Where("user_id = ? and status = ?", userId, 1)
	err := s.db.Where("id in (?)", subQuery).Find(&friends)

	return friends, err.Error
}

// GetFriendRelation 根据用户id和好友id查询好友关系
func (s *FriendRepository) GetFriendRelation(userId, friendId int) (models.Friend, error) {
	var friend models.Friend
	err := s.db.Where("user_id = ? and friend_id = ?", userId, friendId).First(&friend)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return models.Friend{}, nil
	}

	return friend, err.Error
}

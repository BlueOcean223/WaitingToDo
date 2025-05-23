package repository

import (
	"back/models"
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
	subQuery := s.db.Table("friends").Select("friend_id").Where("user_id = ?", userId)
	err := s.db.Where("id in (?)", subQuery).Find(&friends)

	return friends, err.Error
}

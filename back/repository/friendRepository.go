package repository

import (
	"back/models"
	"errors"
	"gorm.io/gorm"
)

type FriendRepository struct {
	Db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) *FriendRepository {
	return &FriendRepository{Db: db}
}

// GetFriendList 根据用户id查询好友列表
func (s *FriendRepository) GetFriendList(userId int) ([]models.User, error) {
	var friends []models.User
	subQuery := s.Db.Table("friends").Select("friend_id").Where("user_id = ? and status = ?", userId, 1)
	err := s.Db.Where("id in (?)", subQuery).Find(&friends)

	return friends, err.Error
}

// GetFriendRelation 根据用户id和好友id查询好友关系
func (s *FriendRepository) GetFriendRelation(userId, friendId int) (models.Friend, error) {
	var friend models.Friend
	err := s.Db.Where("user_id = ? and friend_id = ?", userId, friendId).First(&friend)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		return models.Friend{}, nil
	}

	return friend, err.Error
}

// AddFriendRequest 添加好友请求
func (s *FriendRepository) AddFriendRequest(friend models.Friend) error {
	return s.Db.Create(&friend).Error
}

// UpdateFriend 更新好友关系
func (s *FriendRepository) UpdateFriend(friend models.Friend) error {
	return s.Db.Save(&friend).Error
}

// DeleteFriend 删除好友关系
func (s *FriendRepository) DeleteFriend(id int) error {
	return s.Db.Where("id = ?", id).Delete(&models.Friend{}).Error
}

// DeleteByUserIdAndFriendId 根据用户id和好友id删除好友关系
func (s *FriendRepository) DeleteByUserIdAndFriendId(userId, friendId int) error {
	return s.Db.Where("user_id = ? AND friend_id = ?", userId, friendId).Delete(&models.Friend{}).Error
}

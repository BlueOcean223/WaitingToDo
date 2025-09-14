package repository

import (
	"backend/internal/models"
	"errors"
	"gorm.io/gorm"
)

type FriendRepository interface {
	GetFriendList(userId int) ([]models.User, error)
	GetFriendRelation(userId, friendId int) (models.Friend, error)
	AddFriendRequest(friend *models.Friend, tx *gorm.DB) error
	UpdateFriend(friend models.Friend, tx *gorm.DB) error
	DeleteFriend(id int, tx *gorm.DB) error
	DeleteByUserIdAndFriendId(userId, friendId int, tx *gorm.DB) error
	GetIsFriend(userId, friendId int) ([]models.Friend, error)
	DeleteIsFriend(userId, friendId int, tx *gorm.DB) error
}

type friendRepository struct {
	Db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) FriendRepository {
	return &friendRepository{Db: db}
}

// GetFriendList 根据用户id查询好友列表
func (s *friendRepository) GetFriendList(userId int) ([]models.User, error) {
	var friends []models.User
	subQuery := s.Db.Table("friends").Select("friend_id").Where("user_id = ? and status = ?", userId, 1)
	err := s.Db.Where("id in (?)", subQuery).Find(&friends).Error

	return friends, err
}

// GetFriendRelation 根据用户id和好友id查询好友关系
func (s *friendRepository) GetFriendRelation(userId, friendId int) (models.Friend, error) {
	var friend models.Friend
	err := s.Db.Where("user_id = ? and friend_id = ?", userId, friendId).First(&friend).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Friend{}, nil
	}

	return friend, err
}

// AddFriendRequest 添加好友请求
func (s *friendRepository) AddFriendRequest(friend *models.Friend, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Create(&friend).Error
}

// UpdateFriend 更新好友关系
func (s *friendRepository) UpdateFriend(friend models.Friend, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Model(&friend).Updates(&friend).Error
}

// DeleteFriend 删除好友关系
func (s *friendRepository) DeleteFriend(id int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("id = ?", id).Delete(&models.Friend{}).Error
}

// DeleteByUserIdAndFriendId 根据用户id和好友id删除好友关系
func (s *friendRepository) DeleteByUserIdAndFriendId(userId, friendId int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("user_id = ? AND friend_id = ?", userId, friendId).Delete(&models.Friend{}).Error
}

// GetIsFriend 查询是否已经是好友关系
func (s *friendRepository) GetIsFriend(userId, friendId int) ([]models.Friend, error) {
	var friends []models.Friend
	err := s.Db.Where("user_id = ? AND friend_id = ?", userId, friendId).Find(&friends).Error
	return friends, err
}

// DeleteIsFriend 删除已经是好友的多余好友关系
func (s *friendRepository) DeleteIsFriend(userId, friendId int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("user_id = ? AND friend_id = ? AND status = ?", userId, friendId, 0).Delete(&models.Friend{}).Error
}

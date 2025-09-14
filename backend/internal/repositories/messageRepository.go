package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	InsertMessage(message models.Message, tx *gorm.DB) error
	GetUnreadMessageCount(userId int) (int64, error)
	GetMessageList(page, pageSize, userId int) ([]models.Message, error)
	Update(message models.Message, tx *gorm.DB) error
	Delete(id int, tx *gorm.DB) error
	ReadAllMessage(userId int, tx *gorm.DB) error
	SelectUserInfoByIds(ids []int) ([]models.User, error)
}

type messageRepository struct {
	Db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{Db: db}
}

// InsertMessage 插入消息
func (s *messageRepository) InsertMessage(message models.Message, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Create(&message).Error
}

// GetUnreadMessageCount GetUnreadMessage 查询用户未读消息数量
func (s *messageRepository) GetUnreadMessageCount(userId int) (int64, error) {
	var count int64
	err := s.Db.Table("messages").Where("to_id = ? and is_read = ?", userId, 0).Count(&count).Error
	return count, err
}

// GetMessageList 查询用户消息列表
func (s *messageRepository) GetMessageList(page, pageSize, userId int) ([]models.Message, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	var messages []models.Message
	// 分页查询
	err := s.Db.Where("to_id = ?", userId).Order("send_time desc, id").
		Offset(offset).Limit(pageSize).Find(&messages).Error
	return messages, err
}

// Update 更新消息
func (s *messageRepository) Update(message models.Message, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Model(&message).Updates(&message).Error
}

// Delete 删除消息
func (s *messageRepository) Delete(id int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Delete(&models.Message{}, id).Error
}

// ReadAllMessage 将全部消息设置为已读
func (s *messageRepository) ReadAllMessage(userId int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Model(&models.Message{}).Where("to_id = ?", userId).Update("is_read", 1).Error
}

// SelectUserInfoByIds 根据用户id获取用户信息
func (s *messageRepository) SelectUserInfoByIds(ids []int) ([]models.User, error) {
	var users []models.User
	err := s.Db.Where("id in ?", ids).Find(&users).Error

	return users, err
}

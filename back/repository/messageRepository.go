package repository

import (
	"back/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	Db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{Db: db}
}

// InsertMessage 插入消息
func (s *MessageRepository) InsertMessage(message models.Message, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Create(&message).Error
}

// GetUnreadMessageCount GetUnreadMessage 查询用户未读消息数量
func (s *MessageRepository) GetUnreadMessageCount(userId int) (int64, error) {
	var count int64
	err := s.Db.Table("messages").Where("to_id = ? and is_read = ?", userId, 0).Count(&count).Error
	return count, err
}

// GetMessageList 查询用户消息列表
func (s *MessageRepository) GetMessageList(page, pageSize, userId int) ([]models.Message, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize

	var messages []models.Message
	// 分页查询
	err := s.Db.Where("to_id = ?", userId).Order("send_time desc").
		Offset(offset).Limit(pageSize).Find(&messages).Error
	return messages, err
}

// Update 更新消息
func (s *MessageRepository) Update(message models.Message, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Save(&message).Error
}

// Delete 删除消息
func (s *MessageRepository) Delete(id int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Delete(&models.Message{}, id).Error
}

// ReadAllMessage 将全部消息设置为已读
func (s *MessageRepository) ReadAllMessage(userId int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Model(&models.Message{}).Where("to_id = ?", userId).Update("is_read", 1).Error
}

// SelectUserInfoByIds 根据用户id获取用户信息
func (s *MessageRepository) SelectUserInfoByIds(ids []int) ([]models.User, error) {
	var users []models.User
	err := s.Db.Where("id in ?", ids).Find(&users).Error

	return users, err
}

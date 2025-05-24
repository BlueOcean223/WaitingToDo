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
func (s *MessageRepository) InsertMessage(message models.Message) error {
	return s.Db.Create(&message).Error
}

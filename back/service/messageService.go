package service

import (
	"back/models"
	"back/models/dto"
	"back/repository"
)

type MessageService struct {
	MessageRepository *repository.MessageRepository
}

func NewMessageService(messageRepository *repository.MessageRepository) *MessageService {
	return &MessageService{
		MessageRepository: messageRepository,
	}
}

// GetUnreadMessageCount 查询用户未读消息数量
func (s *MessageService) GetUnreadMessageCount(userId int) (int64, error) {
	return s.MessageRepository.GetUnreadMessageCount(userId)
}

// GetMessageList 获取用户消息列表
func (s *MessageService) GetMessageList(page, pageSize, userId int) ([]dto.MessageDto, error) {
	messages, err := s.MessageRepository.GetMessageList(page, pageSize, userId)
	if err != nil {
		return nil, err
	}

	// 封装成Dto列表
	var messageDtoList []dto.MessageDto
	for _, message := range messages {
		messageDtoList = append(messageDtoList, dto.MessageDto{
			Id:          message.Id,
			Title:       message.Title,
			Description: message.Description,
			SendTime:    message.SendTime,
			FromId:      message.FromId,
			ToId:        message.ToId,
			IsRead:      message.IsRead,
			Type:        message.Type,
			OutId:       message.OutId,
		})
	}
	return messageDtoList, nil
}

// UpdateMessage 更新信息
func (s *MessageService) UpdateMessage(message models.Message) error {
	return s.MessageRepository.Update(message)
}

// DeleteMessage 删除信息
func (s *MessageService) DeleteMessage(messageId int) error {
	return s.MessageRepository.Delete(messageId)
}

// ReadAllMessage 全部已读
func (s *MessageService) ReadAllMessage(userId int) error {
	return s.MessageRepository.ReadAllMessage(userId)
}

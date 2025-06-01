package service

import (
	"back/configs"
	"back/models"
	"back/models/dto"
	"back/models/vo"
	"back/repository"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"time"
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

// HandleRequest 处理请求
func (s *MessageService) HandleRequest(messageVo vo.MessageVo) error {
	// 更新消息状态为已读
	message := models.Message{
		Id:          messageVo.Id,
		Title:       messageVo.Title,
		Description: messageVo.Description,
		FromId:      messageVo.FromId,
		ToId:        messageVo.ToId,
		Type:        messageVo.Type,
		SendTime:    messageVo.SendTime,
		OutId:       messageVo.OutId,
		IsRead:      1,
	}
	err := s.UpdateMessage(message)
	if err != nil {
		return err
	}
	// 如果是团队消息，且为拒绝。则无需发消息
	if messageVo.Type == 2 && messageVo.RequestAction == 0 {
		return nil
	}

	// 使用消息队列发送处理消息
	return s.SendMQMessage(messageVo)
}
func (s *MessageService) SendMQMessage(messageVo vo.MessageVo) error {
	// 准备MQ消息
	mqMsg := models.MQMessage{
		MessageType: messageVo.Type,
		ActionType:  models.FriendRequestReject,
		RelationID:  messageVo.OutId,
		RequesterID: messageVo.FromId,
		ReceiverID:  messageVo.ToId,
	}

	if messageVo.RequestAction == 1 {
		if messageVo.Type == 1 {
			mqMsg.ActionType = models.FriendRequestAccept
		} else {
			mqMsg.ActionType = models.TeamRequestAccept
		}
	}

	// 准备MQ连接
	MQConn := configs.RabbitMQConn
	channel, err := MQConn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// 声明交换机
	err = channel.ExchangeDeclare(
		"social",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	body, err := json.Marshal(mqMsg)
	if err != nil {
		return err
	}
	var routingKey string
	if messageVo.Type == 1 {
		// 添加好友
		routingKey = configs.AppConfigs.RabbitMQConfig.Queues["friend_request"].RoutingKey
	} else {
		// 小组邀请
		routingKey = configs.AppConfigs.RabbitMQConfig.Queues["team_request"].RoutingKey
	}
	// 发送消息
	err = channel.Publish(
		"social",
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	return err
}

// AddMessage 添加消息
func (s *MessageService) AddMessage(message models.Message) error {
	// 填充消息
	message.SendTime = time.Now().Format("2006-01-02 15:04:05")
	// 插入数据库
	return s.MessageRepository.InsertMessage(message)
}

package services

import (
	"backend/internal/configs"
	"backend/internal/models"
	"backend/internal/models/dto"
	"backend/internal/models/vo"
	"backend/internal/repositories"
	"backend/pkg/logger"
	"backend/pkg/myError"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

type MessageService interface {
	GetUnreadMessageCount(userId int) (int64, error)
	GetMessageList(page, pageSize, userId int) ([]dto.MessageDto, error)
	UpdateMessage(message models.Message) error
	DeleteMessage(messageId int) error
	ReadAllMessage(userId int) error
	HandleRequest(messageVo vo.MessageVo) error
	AddMessage(message models.Message) error
	SendMQMessage(messageVo vo.MessageVo) error
}

type messageService struct {
	MessageRepository repository.MessageRepository
}

func NewMessageService(messageRepository repository.MessageRepository) MessageService {
	return &messageService{
		MessageRepository: messageRepository,
	}
}

// GetUnreadMessageCount 查询用户未读消息数量
func (s *messageService) GetUnreadMessageCount(userId int) (int64, error) {
	count, err := s.MessageRepository.GetUnreadMessageCount(userId)
	if err != nil {
		logger.Error("查询用户未读消息数量失败",
			logger.Int("userId", userId),
			logger.Err(err))
		return 0, err
	}
	return count, nil
}

// GetMessageList 获取用户消息列表
func (s *messageService) GetMessageList(page, pageSize, userId int) ([]dto.MessageDto, error) {
	messages, err := s.MessageRepository.GetMessageList(page, pageSize, userId)
	if err != nil {
		logger.Error("获取用户消息列表失败",
			logger.Int("userId", userId),
			logger.Err(err))
		return nil, err
	}

	// 收集用户id
	var userIds []int
	for _, message := range messages {
		userIds = append(userIds, message.FromId)
	}

	// 收集用户信息
	userInfoMap := make(map[int]models.User)
	userInfoList, err := s.MessageRepository.SelectUserInfoByIds(userIds)
	if err != nil {
		logger.Error("查询用户信息列表失败",
			logger.Err(err))
		return nil, err
	}
	for _, userInfo := range userInfoList {
		userInfoMap[userInfo.Id] = userInfo
	}

	// 封装成Dto列表
	var messageDtoList []dto.MessageDto
	for _, message := range messages {
		// 如果为普通消息，则根据用户昵称初始化消息标题
		title := message.Title
		if message.Type == 0 {
			title = fmt.Sprintf(title, userInfoMap[message.FromId].Name)
		}
		// 如果是添加好友或小组邀请请求，则根据用户昵称初始化消息内容
		description := message.Description
		if message.Type == 1 || message.Type == 2 {
			description = fmt.Sprintf(description, userInfoMap[message.FromId].Name)
		}
		messageDtoList = append(messageDtoList, dto.MessageDto{
			Id:          message.Id,
			Title:       title,
			Description: description,
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
func (s *messageService) UpdateMessage(message models.Message) error {
	err := s.MessageRepository.Update(message, nil)
	if err != nil {
		logger.Error("更新消息失败",
			logger.Int("id", message.Id),
			logger.Err(err))
	}
	return err
}

// DeleteMessage 删除信息
func (s *messageService) DeleteMessage(messageId int) error {
	err := s.MessageRepository.Delete(messageId, nil)
	if err != nil {
		logger.Error("删除消息失败",
			logger.Int("id", messageId),
			logger.Err(err))
	}
	return err
}

// ReadAllMessage 全部已读
func (s *messageService) ReadAllMessage(userId int) error {
	err := s.MessageRepository.ReadAllMessage(userId, nil)
	if err != nil {
		logger.Error("消息一键已读失败",
			logger.Int("id", userId),
			logger.Err(err))
	}
	return err
}

// HandleRequest 处理请求
func (s *messageService) HandleRequest(messageVo vo.MessageVo) error {
	// 更新消息状态为已读
	message := models.Message{
		Id:     messageVo.Id,
		IsRead: 1,
	}
	err := s.UpdateMessage(message)
	if err != nil {
		return err
	}
	// 特殊处理团队消息
	if messageVo.Type == 2 {
		// 如果是拒绝，则直接返回
		if messageVo.RequestAction == 0 {
			return nil
		}
		// 如果是同意，则先检查小组任务是否还存在
		taskRepo := repository.NewTaskRepository(configs.MysqlDb)
		teamTask, err := taskRepo.GetTaskListByIds([]int{messageVo.OutId})
		if err != nil {
			logger.Error("接受小组邀请失败",
				logger.Int("taskId", messageVo.OutId),
				logger.Int("messageId", messageVo.Id),
				logger.Err(err))
			return err
		}
		if len(teamTask) == 0 {
			// 小组任务以及不存在，无法接受邀请
			return myError.NewMyError("该小组任务已不存在，无法接受邀请！")
		}
	}

	// 使用消息队列发送处理消息
	return s.SendMQMessage(messageVo)
}
func (s *messageService) SendMQMessage(messageVo vo.MessageVo) error {
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
		logger.Error("连接MQ管道失败",
			logger.Int("id", messageVo.Id),
			logger.Err(err))
		return err
	}
	defer func(channel *amqp091.Channel) {
		err := channel.Close()
		if err != nil {
			logger.Error("关闭MQ管道失败",
				logger.Int("id", messageVo.Id),
				logger.Err(err))
		}
	}(channel)

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
		logger.Error("json格式化失败",
			logger.Int("id", messageVo.Id),
			logger.Err(err))
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
	if err != nil {
		logger.Error("生产者发送消息至MQ失败",
			logger.Int("id", messageVo.Id),
			logger.Err(err))
	}
	return err
}

// AddMessage 添加消息
func (s *messageService) AddMessage(message models.Message) error {
	// 填充消息
	message.SendTime = time.Now().Format("2006-01-02 15:04:05")
	// 插入数据库
	err := s.MessageRepository.InsertMessage(message, nil)
	if err != nil {
		logger.Error("添加消息失败",
			logger.Int("from_id", message.FromId),
			logger.Int("to_id", message.ToId),
			logger.Err(err))
	}
	return err
}

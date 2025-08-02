package consumer

import (
	"back/configs"
	"back/models"
	"back/repository"
	"back/service"
	"encoding/json"
	"log"
)

/*
	用于处理好友请求的消费者
*/

type FriendConsumer struct {
	*BaseConsumer
}

func NewFriendConsumer() *FriendConsumer {
	queueConfig := configs.AppConfigs.RabbitMQConfig.Queues[configs.FriendRequestQueue]
	handler := func(body []byte) error {
		return handleFriendMessage(body)
	}

	baseConsumer := NewBaseConsumer(
		"FriendConsumer",
		queueConfig.Name,
		queueConfig.RoutingKey,
		"social",
		handler,
	)

	return &FriendConsumer{BaseConsumer: baseConsumer}
}

// handleFriendMessage 处理好友请求消息
func handleFriendMessage(body []byte) error {
	// 解析消息
	var mqMessage models.MQMessage
	if err := json.Unmarshal(body, &mqMessage); err != nil {
		return err
	}

	friendRepository := repository.NewFriendRepository(configs.MysqlDb)
	friendService := service.NewFriendService(nil, friendRepository, nil)

	switch mqMessage.ActionType {
	case models.FriendRequestAccept:
		return friendService.AcceptFriendRequest(mqMessage.RequesterID, mqMessage.ReceiverID)
	case models.FriendRequestReject:
		return friendService.RejectFriendRequest(mqMessage.RelationID)
	default:
		log.Printf("未知的好友请求操作类型: %v", mqMessage.ActionType)
		return nil
	}
}

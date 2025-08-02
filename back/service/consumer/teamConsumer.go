package consumer

import (
	"back/configs"
	"back/models"
	"encoding/json"
	"log"
)

/*
	用于处理小组邀请的消费者
*/

type TeamConsumer struct {
	*BaseConsumer
}

func NewTeamConsumer() *TeamConsumer {
	queueConfig := configs.AppConfigs.RabbitMQConfig.Queues[configs.TeamRequestQueue]
	handler := func(body []byte) error {
		return handleTeamMessage(body)
	}

	baseConsumer := NewBaseConsumer(
		"TeamConsumer",
		queueConfig.Name,
		queueConfig.RoutingKey,
		"social",
		handler,
	)

	return &TeamConsumer{BaseConsumer: baseConsumer}
}

// 处理小队邀请消息
func handleTeamMessage(body []byte) error {
	// 解析消息
	var mqMessage models.MQMessage
	if err := json.Unmarshal(body, &mqMessage); err != nil {
		return err
	}

	// 开启事务
	tx := configs.MysqlDb.Begin()

	// 向任务关系表插入数据
	teamTask := models.TeamTask{
		TaskId: mqMessage.RelationID,
		UserId: mqMessage.ReceiverID,
		Status: 0,
	}
	if err := tx.Create(&teamTask).Error; err != nil {
		tx.Rollback()
		log.Printf("TeamConsumer:向任务关系表插入数据异常：%v", err)
	}

	// 将任务表相应任务状态改为未完成
	if err := tx.Model(&models.Task{}).Where("id = ?", mqMessage.RelationID).
		Update("status", 0).Error; err != nil {
		tx.Rollback()
		log.Printf("TeamConsumer:更新任务状态异常：%v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("TeamConsumer:提交事务异常：%v", err)
	}

	return nil
}

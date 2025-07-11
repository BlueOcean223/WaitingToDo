package service

import (
	"back/configs"
	"back/models"
	"back/models/dto"
	"back/repository"
	"back/utils/redisContent"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type FriendService struct {
	authRepository    *repository.AuthRepository
	friendRepository  *repository.FriendRepository
	messageRepository *repository.MessageRepository
}

func NewFriendService(authRepository *repository.AuthRepository,
	friendRepository *repository.FriendRepository, messageRepository *repository.MessageRepository) *FriendService {
	return &FriendService{
		authRepository:    authRepository,
		friendRepository:  friendRepository,
		messageRepository: messageRepository,
	}
}

// GetFriendInfo 根据好友id获取好友信息
func (s *FriendService) GetFriendInfo(friendId int) (dto.FriendInfoDto, error) {
	user, err := s.authRepository.SelectUserById(friendId)
	if user == (models.User{}) {
		return dto.FriendInfoDto{}, err
	}
	friendInfoDto := dto.FriendInfoDto{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		Description: user.Description,
		Pic:         user.Pic,
		IsFriend:    1,
	}
	return friendInfoDto, nil
}

// GetFriendList 根据用户id获取好友列表
func (s *FriendService) GetFriendList(userId int) ([]dto.FriendInfoDto, error) {
	friendList, err := s.friendRepository.GetFriendList(userId)
	if err != nil {
		return nil, err
	}
	var friends []dto.FriendInfoDto
	for _, friend := range friendList {
		friendInfoDto := dto.FriendInfoDto{
			Id:          friend.Id,
			Name:        friend.Name,
			Email:       friend.Email,
			Description: friend.Description,
			Pic:         friend.Pic,
			IsFriend:    1,
		}
		friends = append(friends, friendInfoDto)
	}

	return friends, nil
}

// SearchUserByEmail 根据邮箱搜索用户，用于添加好友
func (s *FriendService) SearchUserByEmail(userEmail, searchEmail string) (dto.FriendInfoDto, error) {
	searchUser, err := s.authRepository.SelectUserByEmail(searchEmail)
	if searchUser == (models.User{}) || err != nil {
		return dto.FriendInfoDto{}, err
	}

	// 根据用户id和好友id查询两者目前关系
	// 根据用户邮箱查询用户
	redisClient := configs.RedisClient
	key := redisContent.UserInfoKey + userEmail
	var user models.User
	if redisClient.Exists(context.Background(), key).Val() == 1 {
		val, err := redisClient.Get(context.Background(), key).Bytes()
		if err != nil {
			return dto.FriendInfoDto{}, err
		}

		err = json.Unmarshal(val, &user)
		if err != nil {
			return dto.FriendInfoDto{}, err
		}
	} else {
		user, err = s.authRepository.SelectUserByEmail(userEmail)
		if err != nil {
			return dto.FriendInfoDto{}, err
		}

		// 将用户信息写入缓存
		val, err := json.Marshal(user)
		if err != nil {
			return dto.FriendInfoDto{}, err
		}
		redisClient.Set(context.Background(), key, val, 24*time.Hour)
	}

	// 查询关系
	relation, err := s.friendRepository.GetFriendRelation(user.Id, searchUser.Id)
	if err != nil {
		return dto.FriendInfoDto{}, err
	}

	var isFriend int
	// 关系不存在
	if relation == (models.Friend{}) {
		isFriend = 2
	} else {
		isFriend = relation.Status
	}

	// 封装返回数据
	friendInfoDto := dto.FriendInfoDto{
		Id:          searchUser.Id,
		Name:        searchUser.Name,
		Email:       searchUser.Email,
		Description: searchUser.Description,
		Pic:         searchUser.Pic,
		IsFriend:    isFriend,
	}
	return friendInfoDto, nil
}

// AddFriend 添加好友请求
func (s *FriendService) AddFriend(userId, friendId int) error {
	// 开启事务
	tx := s.friendRepository.Db.Begin()
	// 向好友关系表写入数据
	friend := models.Friend{
		UserId:   userId,
		FriendId: friendId,
		Status:   0,
	}
	err := s.friendRepository.AddFriendRequest(&friend, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 获取用户信息
	redisClient := configs.RedisClient
	key := fmt.Sprintf(redisContent.UserInfoKey+"%d", userId)
	var user models.User

	if redisClient.Exists(context.Background(), key).Val() == 1 {
		val, err := redisClient.Get(context.Background(), key).Bytes()
		if err != nil {
			tx.Rollback()
			return err
		}

		err = json.Unmarshal(val, &user)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		user, err = s.authRepository.SelectUserById(userId)
		if err != nil {
			tx.Rollback()
			return err
		}

		// 将用户信息写入缓存
		val, err := json.Marshal(user)
		if err != nil {
			tx.Rollback()
			return err
		}

		redisClient.Set(context.Background(), key, val, 24*time.Hour)
	}

	// 向请求添加对象发送消息
	message := models.Message{
		Title:       "好友请求",
		Description: "你好，我是%s，我想与你成为好友",
		FromId:      userId,
		ToId:        friendId,
		Type:        1,
		SendTime:    time.Now().Format("2006-01-02 15:04:05"),
		OutId:       friend.Id,
		IsRead:      0,
	}
	err = s.messageRepository.InsertMessage(message, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 执行完毕，提交事务
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// AcceptFriendRequest 接受好友请求
func (s *FriendService) AcceptFriendRequest(userId, friendId int) error {
	// 检查两人是否已经是好友
	friendShip, err := s.friendRepository.GetIsFriend(userId, friendId)
	if err != nil {
		return err
	}
	if len(friendShip) > 1 {
		// 两人已经是好友，无需重复添加，同时删除好友关系表中的记录
		err = s.friendRepository.DeleteIsFriend(userId, friendId, nil)
		return err
	}

	// 开启事务
	tx := s.friendRepository.Db.Begin()
	friend, err := s.friendRepository.GetFriendRelation(userId, friendId)
	if err != nil {
		tx.Rollback()
		return err
	}
	friend.Status = 1
	err = s.friendRepository.UpdateFriend(friend, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 反向添加好友
	newFriend := models.Friend{
		UserId:   friend.FriendId,
		FriendId: friend.UserId,
		Status:   1,
	}
	err = s.friendRepository.AddFriendRequest(&newFriend, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// RejectFriendRequest 拒绝好友请求
func (s *FriendService) RejectFriendRequest(id int) error {
	return s.friendRepository.DeleteFriend(id, nil)
}

// StartFriendConsumer 监听消息队列，处理好友请求
func StartFriendConsumer() {
	for {
		// 建立连接
		MQConn, err := amqp091.Dial(configs.AppConfigs.RabbitMQConfig.Dsn)
		if err != nil {
			log.Printf("FriendConsumer:Failed to connect to RabbitMQ: %v", err)
			time.Sleep(time.Minute)
			continue
		}
		// 注册 connection 级别关闭通知
		connCloseCh := make(chan *amqp091.Error)
		MQConn.NotifyClose(connCloseCh)

		// 打开 channel
		channel, err := MQConn.Channel()
		if err != nil {
			log.Printf("FriendConsumer:Failed to open channel: %v, retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// 监听Channel级别的关闭
		chanCloseCh := make(chan *amqp091.Error)
		channel.NotifyClose(chanCloseCh)

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
		if err != nil {
			log.Printf("FriendConsumer:Failed to declare exchange: %v", err)
			channel.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// 声明队列
		queue, err := channel.QueueDeclare(
			"friend_requests",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("FriendConsumer:Failed to declare queue: %v", err)
			channel.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// 绑定队列
		err = channel.QueueBind(
			queue.Name,
			configs.AppConfigs.RabbitMQConfig.Queues["friend_request"].RoutingKey,
			"social",
			false,
			nil,
		)
		if err != nil {
			log.Printf("FriendConsumer:Failed to bind queue: %v", err)
			channel.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// 消费消息
		msgs, err := channel.Consume(
			queue.Name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("FriendConsumer:Failed to register consumer: %v", err)
			channel.Close()
			MQConn.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("FriendConsumer start")

	loop:
		for {
			select {
			case err = <-connCloseCh:
				log.Printf("FriendConsumer:Connection closed: %v", err)
				break loop
			case err = <-chanCloseCh:
				log.Printf("FriendConsumer:Channel closed: %v", err)
				break loop
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("FriendConsumer:Channel closed")
					break loop
				}
				var mqMessage models.MQMessage
				if err := json.Unmarshal(msg.Body, &mqMessage); err != nil {
					log.Printf("FriendConsumer:Failed to unmarshal message: %v", err)
					msg.Nack(false, true) // 重新入队
					continue
				}

				friendRepository := repository.NewFriendRepository(configs.MysqlDb)
				friendService := NewFriendService(nil, friendRepository, nil)

				var processErr error
				switch mqMessage.ActionType {
				case models.FriendRequestAccept:
					processErr = friendService.AcceptFriendRequest(mqMessage.RequesterID, mqMessage.ReceiverID)
				case models.FriendRequestReject:
					processErr = friendService.RejectFriendRequest(mqMessage.RelationID)
				}

				if processErr != nil {
					log.Printf("FriendConsumer:Failed to process message: %v", processErr)
					msg.Nack(false, true) // 处理失败，重新入队
				} else {
					if err := msg.Ack(false); err != nil {
						log.Printf("FriendConsumer:Failed to ack message: %v", err)
					}
				}
			}
		}

		// 断开连接
		err = channel.Close()
		if err != nil {
			log.Printf("FriendConsumer:Failed to close channel: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		err = MQConn.Close()
		if err != nil {
			log.Printf("FriendConsumer:Failed to close connection: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Message channel closed, reconnecting...")
		time.Sleep(5 * time.Second)
	}
}

// DeleteFriend 删除好友
func (s *FriendService) DeleteFriend(userId, friendId int) error {
	// 开启事务
	tx := s.friendRepository.Db.Begin()

	// 双向删除好友关系
	err := s.friendRepository.DeleteByUserIdAndFriendId(userId, friendId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.friendRepository.DeleteByUserIdAndFriendId(friendId, userId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		// 提交事务异常
		tx.Rollback()
		return err
	}
	return nil
}

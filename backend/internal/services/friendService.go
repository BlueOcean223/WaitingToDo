package services

import (
	"backend/internal/configs"
	"backend/internal/models"
	"backend/internal/models/dto"
	"backend/internal/repositories"
	"backend/pkg/logger"
	"backend/pkg/redisContent"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type FriendService interface {
	GetFriendInfo(friendId int) (dto.FriendInfoDto, error)
	GetFriendList(userId int) ([]dto.FriendInfoDto, error)
	SearchUserByEmail(userEmail, searchEmail string) (dto.FriendInfoDto, error)
	AddFriend(userId, friendId int) error
	AcceptFriendRequest(userId, friendId int) error
	RejectFriendRequest(id int) error
	DeleteFriend(userId, friendId int) error
}

type friendService struct {
	authRepository    repository.AuthRepository
	friendRepository  repository.FriendRepository
	messageRepository repository.MessageRepository
}

func NewFriendService(authRepository repository.AuthRepository,
	friendRepository repository.FriendRepository, messageRepository repository.MessageRepository) FriendService {
	return &friendService{
		authRepository:    authRepository,
		friendRepository:  friendRepository,
		messageRepository: messageRepository,
	}
}

// GetFriendInfo 根据好友id获取好友信息
func (s *friendService) GetFriendInfo(friendId int) (dto.FriendInfoDto, error) {
	user, err := s.authRepository.SelectUserById(friendId)
	if err != nil {
		logger.Error("查询好友信息失败",
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
		return dto.FriendInfoDto{}, err
	}
	if user == (models.User{}) {
		logger.Warn("好友不存在",
			logger.String("friend_id", fmt.Sprintf("%d", friendId)))
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
func (s *friendService) GetFriendList(userId int) ([]dto.FriendInfoDto, error) {
	friendList, err := s.friendRepository.GetFriendList(userId)
	if err != nil {
		logger.Error("查询好友列表失败",
			logger.Int("user_id", userId),
			logger.Err(err))
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
func (s *friendService) SearchUserByEmail(userEmail, searchEmail string) (dto.FriendInfoDto, error) {
	searchUser, err := s.authRepository.SelectUserByEmail(searchEmail)
	if searchUser == (models.User{}) || err != nil {
		logger.Warn("搜索用户不存在",
			logger.String("search_email", searchEmail),
			logger.Err(err))
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
			logger.Error("从Redis获取用户信息失败",
				logger.String("user_email", userEmail),
				logger.Err(err))
			return dto.FriendInfoDto{}, err
		}

		err = json.Unmarshal(val, &user)
		if err != nil {
			logger.Error("反序列化用户信息失败",
				logger.String("user_email", userEmail),
				logger.Err(err))
			return dto.FriendInfoDto{}, err
		}
	} else {
		user, err = s.authRepository.SelectUserByEmail(userEmail)
		if err != nil {
			logger.Error("查询用户信息失败",
				logger.String("user_email", userEmail),
				logger.Err(err))
			return dto.FriendInfoDto{}, err
		}

		// 将用户信息写入缓存
		val, err := json.Marshal(user)
		if err != nil {
			logger.Error("序列化用户信息失败",
				logger.String("user_email", userEmail),
				logger.Err(err))
			return dto.FriendInfoDto{}, err
		}
		redisClient.Set(context.Background(), key, val, 24*time.Hour)
	}

	// 查询关系
	relation, err := s.friendRepository.GetFriendRelation(user.Id, searchUser.Id)
	if err != nil {
		logger.Error("查询好友关系失败",
			logger.String("user_id", fmt.Sprintf("%d", user.Id)),
			logger.String("friend_id", fmt.Sprintf("%d", searchUser.Id)),
			logger.Err(err))
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
func (s *friendService) AddFriend(userId, friendId int) error {
	// 开启事务
	tx := configs.MysqlDb.Begin()
	// 向好友关系表写入数据
	friend := models.Friend{
		UserId:   userId,
		FriendId: friendId,
		Status:   0,
	}
	err := s.friendRepository.AddFriendRequest(&friend, tx)
	if err != nil {
		tx.Rollback()
		logger.Error("添加好友关系失败",
			logger.Int("user_id", userId),
			logger.Int("friend_id", friendId),
			logger.Err(err))
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
			logger.Error("从Redis获取用户信息失败",
				logger.Int("user_id", userId),
				logger.Err(err))
			return err
		}

		err = json.Unmarshal(val, &user)
		if err != nil {
			tx.Rollback()
			logger.Error("反序列化用户信息失败",
				logger.Int("user_id", userId),
				logger.Err(err))
			return err
		}
	} else {
		user, err = s.authRepository.SelectUserById(userId)
		if err != nil {
			tx.Rollback()
			logger.Error("查询用户信息失败",
				logger.Int("user_id", userId),
				logger.Err(err))
			return err
		}

		// 将用户信息写入缓存
		val, err := json.Marshal(user)
		if err != nil {
			tx.Rollback()
			logger.Error("序列化用户信息失败",
				logger.Int("user_id", userId),
				logger.Err(err))
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
		logger.Error("发送好友请求消息失败",
			logger.Int("user_id", userId),
			logger.Int("friend_id", friendId),
			logger.Err(err))
		return err
	}

	// 执行完毕，提交事务
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		logger.Error("提交好友请求事务失败",
			logger.Int("user_id", userId),
			logger.Int("friend_id", friendId),
			logger.Err(err))
		return err
	}
	return nil
}

// AcceptFriendRequest 接受好友请求
func (s *friendService) AcceptFriendRequest(userId, friendId int) error {
	// 检查两人是否已经是好友
	friendShip, err := s.friendRepository.GetIsFriend(userId, friendId)
	if err != nil {
		logger.Error("查询好友关系失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
		return err
	}
	if len(friendShip) > 1 {
		// 两人已经是好友，无需重复添加，同时删除好友关系表中的记录
		err = s.friendRepository.DeleteIsFriend(userId, friendId, nil)
		if err != nil {
			logger.Error("删除重复好友关系失败",
				logger.String("user_id", fmt.Sprintf("%d", userId)),
				logger.String("friend_id", fmt.Sprintf("%d", friendId)),
				logger.Err(err))
			return err
		}
		logger.Info("好友关系已存在，删除重复记录",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)))
		return err
	}

	// 开启事务
	tx := configs.MysqlDb.Begin()
	friend, err := s.friendRepository.GetFriendRelation(userId, friendId)
	if err != nil {
		tx.Rollback()
		logger.Error("获取好友关系失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
		return err
	}
	friend.Status = 1
	err = s.friendRepository.UpdateFriend(friend, tx)
	if err != nil {
		tx.Rollback()
		logger.Error("更新好友关系失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
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
		logger.Error("反向添加好友关系失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		logger.Error("提交好友关系事务失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
		return err
	}
	return nil
}

// RejectFriendRequest 拒绝好友请求
func (s *friendService) RejectFriendRequest(id int) error {
	err := s.friendRepository.DeleteFriend(id, nil)
	if err != nil {
		logger.Error("拒绝好友请求失败",
			logger.String("friend_request_id", fmt.Sprintf("%d", id)),
			logger.Err(err))
		return err
	}

	return nil
}

// DeleteFriend 删除好友
func (s *friendService) DeleteFriend(userId, friendId int) error {
	// 开启事务
	tx := configs.MysqlDb.Begin()

	// 双向删除好友关系
	err := s.friendRepository.DeleteByUserIdAndFriendId(userId, friendId, tx)
	if err != nil {
		tx.Rollback()
		logger.Error("删除好友关系失败(1)",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
		return err
	}

	err = s.friendRepository.DeleteByUserIdAndFriendId(friendId, userId, tx)
	if err != nil {
		tx.Rollback()
		logger.Error("删除好友关系失败(2)",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
		return err
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		// 提交事务异常
		tx.Rollback()
		logger.Error("提交删除好友事务失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.String("friend_id", fmt.Sprintf("%d", friendId)),
			logger.Err(err))
		return err
	}
	return nil
}

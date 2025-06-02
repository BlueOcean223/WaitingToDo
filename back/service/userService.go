package service

import (
	"back/configs"
	"back/models"
	"back/models/vo"
	"back/repository"
	"back/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type UserService struct {
	authRepository *repository.AuthRepository
}

func NewUserService(authRepository *repository.AuthRepository) *UserService {
	return &UserService{authRepository: authRepository}
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(email, password string) error {
	// 检查数据库是否有该邮箱的用户
	user, err := s.authRepository.SelectUserByEmail(email)
	// 查询异常
	if err != nil {
		return err
	}
	// 用户不存在
	if user == (models.User{}) {
		return utils.NewMyError("用户不存在")
	}
	// 更新密码
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashPassword
	return s.authRepository.UpdateUser(user)
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(userVo vo.UserVo) error {
	user, err := s.authRepository.SelectUserByEmail(userVo.Email)
	if err != nil {
		return err
	}
	if user == (models.User{}) {
		return utils.NewMyError("用户不存在")
	}

	// 先删除缓存
	redisClient := configs.RedisClient
	emailKey := utils.UserInfoKey + user.Email
	idKey := fmt.Sprintf(utils.UserInfoKey+"%d", user.Id)
	err = redisClient.Del(context.Background(), emailKey, idKey).Err()
	if err != nil {
		return err
	}

	// 更新个人信息
	user.Name = userVo.Name
	user.Description = userVo.Description
	err = s.authRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	// 更新缓存
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = redisClient.Set(context.Background(), emailKey, userJson, 24*time.Hour).Err()
	err = redisClient.Set(context.Background(), idKey, userJson, 24*time.Hour).Err()

	return err
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(id int) (vo.UserVo, error) {
	// 先从缓存中获取
	redisClient := configs.RedisClient
	key := fmt.Sprintf(utils.UserInfoKey+"%d", id)
	var user models.User
	var err error

	if redisClient.Exists(context.Background(), key).Val() == 1 {
		val, err := redisClient.Get(context.Background(), key).Bytes()
		if err != nil {
			return vo.UserVo{}, err
		}
		err = json.Unmarshal(val, &user)
		if err != nil {
			return vo.UserVo{}, err
		}
	} else {
		// 从数据库中获取
		user, err = s.authRepository.SelectUserById(id)
		if err != nil {
			return vo.UserVo{}, err
		}
	}

	return vo.UserVo{
		Name:        user.Name,
		Description: user.Description,
		Email:       user.Email,
		Pic:         user.Pic,
	}, nil
}
